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
	"bufio"
	"fmt"
	"os"

	"github.com/m4rs14n/go-app/shared"
)

func main() {
	done := make(chan struct{})
	unixBus, err := shared.LoadPlugin("unix_bus.so")
	if err != nil {
		panic(err)
	}
	unixBus.Start(done)

	plugin, err := shared.LoadPlugin("my_plugin.so")
	if err != nil {
		panic(err)
	}
	plugin.Start(done)

	if storage, ok := plugin.(shared.Storage); ok {
		if ch := storage.Read("/path"); ch != nil {
			response := <-ch
			fmt.Printf("Read: %v\n", response)
		}
	}

	// Read a char to make sure we get all the messages (boardcast is async)
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	shared.StopAllPlugins(done)
}
