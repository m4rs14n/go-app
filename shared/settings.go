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
	"github.com/google/uuid"
)

// Settings is a map of Plugin settings/configurations
type Settings map[string]interface{}

const (
	// SettingID is id constant key
	SettingID = "id"
	// SettingName is name constant key
	SettingName = "name"
	// SettingDescription is description constant key
	SettingDescription = "description"
)

// SetupSettings is a helper function to create plugin settings
func SetupSettings(id uuid.UUID, name string, description string) Settings {
	return Settings{
		SettingID:          id,
		SettingName:        name,
		SettingDescription: description,
	}
}

// ID returns the id field
func (s Settings) ID() uuid.UUID {
	return s[SettingID].(uuid.UUID)
}

// Name returns the name field
func (s Settings) Name() string {
	return s[SettingName].(string)
}

// Description returns the description field
func (s Settings) Description() string {
	return s[SettingDescription].(string)
}
