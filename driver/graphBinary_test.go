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
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func writeToBuffer(value interface{}, buffer *bytes.Buffer) []byte {
	writer := GraphBinaryWriter{}
	val, err := writer.write(value, buffer)
	if err != nil {
		panic(err)
	}
	return val.([]byte)
}

func readToValue(buff *bytes.Buffer) interface{} {
	reader := GraphBinaryReader{}
	val, err := reader.read(buff)
	if err != nil {
		panic(err)
	}
	return val
}

func TestGraphBinaryV1(t *testing.T) {
	t.Run("test simple types", func(t *testing.T) {
		buff := bytes.Buffer{}
		t.Run("test int", func(t *testing.T) {
			var x int = 33000
			writeToBuffer(x, &buff)
			val := readToValue(&buff).(int64)
			assert.Equal(t, x, int(val))
		})
		t.Run("test int32", func(t *testing.T) {
			var x int32 = 33000
			writeToBuffer(x, &buff)
			assert.Equal(t, x, readToValue(&buff))
		})
		t.Run("test long", func(t *testing.T) {
			var x int64 = 2147483648
			writeToBuffer(x, &buff)
			assert.Equal(t, x, readToValue(&buff))
		})
		t.Run("test string", func(t *testing.T) {
			var x = "serialize this!"
			writeToBuffer(x, &buff)
			assert.Equal(t, x, readToValue(&buff))
		})
		t.Run("test string", func(t *testing.T) {
			var x, _ = uuid.Parse("41d2e28a-20a4-4ab0-b379-d810dede3786")
			writeToBuffer(x, &buff)
			assert.Equal(t, x, readToValue(&buff))
		})
	})

	t.Run("test nested types", func(t *testing.T) {
		buff := bytes.Buffer{}
		t.Run("test map", func(t *testing.T) {
			var x int32 = 666
			var m = map[interface{}]interface{}{
				"marko": x,
				"noone": "blah",
			}
			writeToBuffer(m, &buff)
			assert.Equal(t, m, readToValue(&buff))
		})
	})
}
