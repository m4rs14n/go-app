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
	"errors"
	"log"
	"os"
	"path/filepath"
	plgin "plugin"
)

// Plugin is the general interface to implement
type Plugin interface {
	Start(done <-chan struct{})
	GetSettings() Settings
	Stop()
}

// PluginListener is the interface that, when implemented, recieves the notifications of plugins getting loaded
type PluginListener interface {
	PluginLoaded(plugin Plugin)
}

// SimplePlugin implements basic functionality of a plugin
type SimplePlugin struct {
	Settings Settings
}

// Make sure we implement both interfaces
var _ Plugin = (*SimplePlugin)(nil)

// Start method
func (p *SimplePlugin) Start(done <-chan struct{}) {
	log.Printf("Starting plugin %s", p.Settings.Name())
}

// GetSettings returns the settings/configurations
func (p *SimplePlugin) GetSettings() Settings {
	return p.Settings
}

// Stop method
func (p *SimplePlugin) Stop() {
	log.Printf("Stopping plugin %s", p.Settings.Name())
}

var plugins = make(map[string]Plugin)
var listeners []PluginListener

// LoadPlugin is the helper to load a plugin
func LoadPlugin(path string) (plugin Plugin, err error) {
	var dylib *plgin.Plugin
	if dylib, err = plgin.Open(path); err != nil {
		return
	}

	var factorySymbol plgin.Symbol
	if factorySymbol, err = dylib.Lookup("NewPlugin"); err != nil {
		return
	}

	switch factory := factorySymbol.(type) {
	case func() (Plugin, error):
		plugin, err = factory()
	case *func() (Plugin, error):
		plugin, err = (*factory)()
	default:
		return nil, errors.New("Cannot cast the factory function")
	}

	if err != nil {
		return
	}

	// TODO: verify unique plugin uuid and other error checks

	for _, listener := range listeners {
		// Send to all the listeners already loaded
		listener.PluginLoaded(plugin)
	}

	if listener, ok := plugin.(PluginListener); ok {
		for _, plugin := range plugins {
			// If the loaded plugin is a listener then send all the loaded plugins to it
			listener.PluginLoaded(plugin)
		}

		listeners = append(listeners, listener)
	}

	plugins[plugin.GetSettings().Name()] = plugin

	return
}

// LoadAllPlugins loads all the plugins in a directory
func LoadAllPlugins(libDir string) chan<- struct{} {
	ch := make(chan struct{})
	filepath.Walk(libDir, func(path string, info os.FileInfo, err error) error {
		if info.Mode()&os.ModeDir == 0 {
			if ext := filepath.Ext(path); ext == ".so" {
				plugin, err := LoadPlugin(path)
				if err != nil {
					log.Print(err)
					return err
				}

				plugin.Start(ch)
			}
		}
		return nil
	})
	return ch
}

// StopAllPlugins stops all plugins by closing the done channel
func StopAllPlugins(done chan<- struct{}) {
	close(done)
	for _, plugin := range plugins {
		plugin.Stop()
	}
}

// GetPlugin returns a loaded plugin
func GetPlugin(name string) Plugin {
	return plugins[name]
}
