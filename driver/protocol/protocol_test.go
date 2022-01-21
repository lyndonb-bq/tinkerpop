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

package protocol

import (
	"github.com/stretchr/testify/assert"
	"gremlin-go/driver/results"
	"testing"
)

// TODO: remove this file when sandbox is no longer needed
func Test(t *testing.T) {

	t.Run("Test DataReceived nil message", func(t *testing.T) {
		protocol := NewGremlinServerWSProtocol()
		statusCode, err := protocol.DataReceived(nil, map[string]results.ResultSet{})
		assert.Equal(t, 0, statusCode)
		assert.Nil(t, err)
	})

	t.Run("Test DataReceived nil message", func(t *testing.T) {
		protocol := NewGremlinServerWSProtocol()
		protocol.DataReceived(nil, nil)
	})

	t.Run("Test Protocol Connection Made", func(t *testing.T) {
		protocol := NewGremlinServerWSProtocol()
		protocol.DataReceived(nil, nil)
	})
}
