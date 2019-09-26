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
	"encoding/gob"

	"github.com/google/uuid"
)

var UnencryptedStorageUUID = uuid.MustParse("A700A163-BDFE-4AE4-A357-E5E28389C3E7")
var EncryptedStorageUUID = uuid.MustParse("C9BFE170-745F-4C7B-954F-95BEA16AA3EC")

// Storage is the interface used to store
type Storage interface {
	// Log is the logging method
	Read(path string) <-chan interface{}
	Write(path string, value interface{})
}

// UseStorage allows a plugin to have access to storage
type UseStorage struct {
	UUID uuid.UUID
	UseBus
}

// Make sure UseStorage implements required interfces
var _ Storage = (*UseStorage)(nil)

const (
	StorageMessageTypeRead = iota
	StorageMessageTypeWrite
)

type StorageMessage struct {
	Type  int
	Path  string
	Value interface{}
}

// Read is the read method
func (s *UseStorage) Read(path string) <-chan interface{} {
	return s.SendMessage(s.UUID, StorageMessage{Type: StorageMessageTypeRead, Path: path})
}

// Write is the write method
func (s *UseStorage) Write(path string, value interface{}) {
	s.SendMessage(s.UUID, StorageMessage{Type: StorageMessageTypeWrite, Path: path, Value: value})
}

func init() {
	gob.Register(StorageMessage{})
}
