// Copyright (c) 2019, Hojat Parta
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  * Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//  * Neither the name of  nor the names of its contributors may be used to
//    endorse or promote products derived from this software without specific
//    prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/m4rs14n/go-app/shared"
)

var id = uuid.MustParse("5AB218CD-A9D1-41A6-877D-5454AF9994C2")
var settings = shared.SetupSettings(id, "UnixBus", "This is a bus plugin")

// bus is the Bus Plugin
type unixBus struct {
	shared.SimplePlugin
	listeners []net.Listener
	waitGroup sync.WaitGroup
}

// Make sure we implement required interfaces
var _ shared.PluginListener = (*unixBus)(nil)
var _ shared.BusService = (*unixBus)(nil)

var instance = &unixBus{
	SimplePlugin: shared.SimplePlugin{Settings: settings},
}

// NewPlugin returns an instance of the plugin
func NewPlugin() (shared.Plugin, error) {
	return instance, nil
}

const endpointsDir = "/tmp/endpoints"

const (
	messageTypeBroadcast = iota
	messageTypeSend
	messageTypeResult
)

type message struct {
	Type    int
	UUID    uuid.UUID
	Message interface{}
}

// Stop method
func (b *unixBus) Stop() {
	for _, l := range b.listeners {
		l.Close()
	}
	b.waitGroup.Wait()
	b.SimplePlugin.Stop()
}

// PluginLoaded allows the plugin to check if a loaded plugin is of any interest
func (b *unixBus) PluginLoaded(plugin shared.Plugin) {
	if endpoint, ok := plugin.(shared.Endpoint); ok {
		uuid := plugin.GetSettings().ID()
		sockAddr := fmt.Sprintf("%s/%v.sock", endpointsDir, uuid)

		if err := os.RemoveAll(sockAddr); err != nil {
			log.Fatal(err)
		}

		l, err := net.Listen("unix", sockAddr)
		if err != nil {
			log.Fatal("listen error: ", err)
		}

		b.listeners = append(b.listeners, l)

		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					// log.Printf("accept error: %v", err)
					break
				}

				b.waitGroup.Add(1)
				go func(c net.Conn) {
					defer c.Close()
					defer b.waitGroup.Done()

					dec := gob.NewDecoder(c)
					var req message
					err := dec.Decode(&req)

					if err == nil {
						switch req.Type {
						case messageTypeBroadcast:
							endpoint.HandleBroadcast(req.Message)

						case messageTypeSend:
							ch, err := endpoint.HandleMessage(req.Message)
							if err == nil {
								enc := gob.NewEncoder(c)
								for res := range ch {
									err = enc.Encode(message{Type: messageTypeResult, Message: res})
									if err != nil {
										break
									}
								}
							}

						default:
							log.Fatal("invalid command")
						}
					}
				}(conn)
			}

			if err := os.RemoveAll(sockAddr); err != nil {
				log.Fatal(err)
			}
		}()
	}
}

func broadcastMessage(uuid uuid.UUID, msg interface{}) {
	sockAddr := fmt.Sprintf("%s/%v.sock", endpointsDir, uuid)

	c, err := net.Dial("unix", sockAddr)
	if err != nil {
		return
	}
	defer c.Close()

	enc := gob.NewEncoder(c)
	err = enc.Encode(message{Type: messageTypeBroadcast, Message: msg})
}

// GetPriority returns the priority of the bus
func (b *unixBus) Priority() int {
	return 100 // Local bus will have a higher priority
}

// BroadcastMessage sends the message to all clients
func (b *unixBus) HandleBroadcast(msg interface{}) {
	filepath.Walk(endpointsDir, func(path string, info os.FileInfo, err error) error {
		if info.Mode()&os.ModeDir == 0 {
			filename := filepath.Base(path)
			ext := filepath.Ext(path)
			name := filename[0 : len(filename)-len(ext)]
			if uuid, err := uuid.Parse(name); err == nil {
				go broadcastMessage(uuid, msg)
			}
		}
		return nil
	})
}

// SendMessage sends the message to a specific client asynchronously
func (b *unixBus) HandleMessage(uuid uuid.UUID, msg interface{}) (<-chan interface{}, error) {
	sockAddr := fmt.Sprintf("%s/%v.sock", endpointsDir, uuid)

	c, err := net.Dial("unix", sockAddr)
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(c)
	err = enc.Encode(message{Type: messageTypeSend, Message: msg})
	if err != nil {
		return nil, err
	}

	ch := make(chan interface{})

	go func() {
		defer close(ch)
		defer c.Close()

		dec := gob.NewDecoder(c)
		var r message
		for {
			err := dec.Decode(&r)
			if err != nil {
				break
			}

			switch r.Type {
			case messageTypeResult:
				ch <- r.Message

			default:
				log.Printf("Invalid response so ignoring")
			}
		}
	}()

	return ch, nil
}

func init() {
	if _, err := os.Stat(endpointsDir); os.IsNotExist(err) {
		os.MkdirAll(endpointsDir, 0700)
	}
}
