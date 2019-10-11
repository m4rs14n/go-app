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
	"errors"
	"log"

	"github.com/m4rs14n/go-app/shared"
)

var settings = shared.SetupSettings(shared.UnencryptedStorageUUID, "UnencryptedStorage", "This is a simple unencrypted storage plugin")

// unencryptedStorage is a unencrypted storage plugin
type unencryptedStorage struct {
	shared.SimplePlugin
}

// Make sure we implement required interfaces
var _ shared.Endpoint = (*unencryptedStorage)(nil)

var instance = &unencryptedStorage{
	shared.SimplePlugin{Settings: settings},
}

// NewPlugin returns an instance of the plugin
func NewPlugin() (shared.Plugin, error) {
	return instance, nil
}

// Simple in memory storage for now
var disk = make(map[string]interface{})

// BroadcastMessage sends the message to all clients
func (s *unencryptedStorage) HandleBroadcast(message interface{}) {
	// Do nothing
}

// SendMessage sends the message to a specific client asynchronously
func (s *unencryptedStorage) HandleMessage(message interface{}) (<-chan interface{}, error) {
	if storageMsg, ok := message.(shared.StorageMessage); ok {
		switch storageMsg.Type {
		case shared.StorageMessageTypeRead:
			ch := make(chan interface{})

			go func() {
				ch <- disk[storageMsg.Path]
				close(ch)
			}()

			return ch, nil

		case shared.StorageMessageTypeWrite:
			disk[storageMsg.Path] = storageMsg.Value
			return nil, nil
		}
	}

	log.Print("Invalid message")
	return nil, errors.New("Invalid message")
}
