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

package shared

import (
	"log"
	"sort"

	"github.com/google/uuid"
)

const (
	// BusSettingPriority
	BusSettingPriority = "priority"
)

// Bus is the interface used to communicate on a bus
type Bus interface {
	// BroadcastMessage sends the message to all clients
	BroadcastMessage(message interface{})
	// TODO: Add error channel in case the bus cannot deliver
	// SendMessageAsync sends the message to a specific client asynchronously
	SendMessage(uuid uuid.UUID, message interface{}) <-chan interface{}
}

// BusService represents an implementation of a bus
type BusService interface {
	// GetPriority returns the priority of the bus
	Priority() int
	// HandleBroadcast handles bus broadcasts
	HandleBroadcast(message interface{})
	// HandleMessage handles bus messages
	HandleMessage(uuid uuid.UUID, message interface{}) (<-chan interface{}, error)
}

// Endpoint represents an endpoint connected to a bus (service)
type Endpoint interface {
	// HandleBroadcast handles bus broadcasts
	HandleBroadcast(message interface{})
	// HandleMessage handles bus messages
	HandleMessage(message interface{}) <-chan interface{}
}

// UseBus allows a plugin to have access to logging
type UseBus struct {
	buses []BusService
}

// Make sure UseBus implements required interfaces
var _ PluginListener = (*UseBus)(nil)
var _ Bus = (*UseBus)(nil)

// ByPriority implements sort.Interface for []BusService based on the priorities
type ByPriority []BusService

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority() < a[j].Priority() }

// PluginLoaded allows the plugin to check if a loaded plugin is of any interest
func (b *UseBus) PluginLoaded(plugin Plugin) {
	if bus, ok := plugin.(BusService); ok {
		b.buses = append(b.buses, bus)
		sort.Sort(ByPriority(b.buses))
	}
}

// BroadcastMessage sends the message to all clients
func (b *UseBus) BroadcastMessage(message interface{}) {
	for _, bus := range b.buses {
		bus.HandleBroadcast(message)
	}
}

// SendMessage sends the message to a specific client asynchronously
func (b *UseBus) SendMessage(uuid uuid.UUID, message interface{}) <-chan interface{} {
	for _, bus := range b.buses {
		resChannel, err := bus.HandleMessage(uuid, message)
		if err != nil {
			return nil
		}
		return resChannel
	}

	log.Printf("Cannot send message")
	return nil
}
