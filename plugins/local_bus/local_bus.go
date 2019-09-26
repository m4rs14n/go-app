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

	"github.com/google/uuid"
	"github.com/m4rs14n/go-app/shared"
)

var id = uuid.MustParse("B49A64D6-8F06-4053-9E30-F5A237EE208A")
var settings = shared.SetupSettings(id, "LocalBus", "This is a bus plugin")

// bus is the Bus Plugin
type localBus struct {
	shared.SimplePlugin
}

// Make sure we implement required interfaces
var _ shared.PluginListener = (*localBus)(nil)
var _ shared.BusService = (*localBus)(nil)

var instance = &localBus{
	shared.SimplePlugin{Settings: settings},
}

// NewPlugin returns an instance of the plugin
func NewPlugin() (shared.Plugin, error) {
	return instance, nil
}

// TODO: Synchronize access to this data structure
var endpoints = make(map[uuid.UUID]shared.Endpoint)

// PluginLoaded allows the plugin to check if a loaded plugin is of any interest
func (b *localBus) PluginLoaded(plugin shared.Plugin) {
	if endpoint, ok := plugin.(shared.Endpoint); ok {
		endpoints[plugin.GetSettings().ID()] = endpoint
	}
}

// GetPriority returns the priority of the bus
func (b *localBus) Priority() int {
	return 0
}

// BroadcastMessage sends the message to all clients
func (b *localBus) HandleBroadcast(message interface{}) {
	for _, endpoint := range endpoints {
		go func(endpoint shared.Endpoint, message interface{}) {
			endpoint.HandleBroadcast(message)
		}(endpoint, message)
	}
}

// SendMessage sends the message to a specific client asynchronously
func (b *localBus) HandleMessage(uuid uuid.UUID, message interface{}) (<-chan interface{}, error) {
	if endpoint, ok := endpoints[uuid]; ok {
		return endpoint.HandleMessage(message), nil
	}

	return nil, errors.New("Invalid endpoint")
}
