/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package gremlingo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test(t *testing.T) {
	t.Run("Test protocol connectionMade", func(t *testing.T) {
		protocol := newGremlinServerWSProtocol(newLogHandler(&defaultLogger{}, Info, language.English))
		transport := getTransportLayer(Gorilla, "host", 1234)
		assert.NotPanics(t, func() { protocol.connectionMade(transport) })
	})

	t.Run("Test dataReceived actual message", func(t *testing.T) {
		protocol := newGremlinServerWSProtocol(newLogHandler(&defaultLogger{}, Info, language.English))
		err := protocol.write("1+1", nil)
		assert.NotNil(t, err)
	})
}
