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
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGraphBinaryV1(t *testing.T) {
	t.Run("test simple types", func(t *testing.T) {
		writer := graphBinaryWriter{}
		reader := graphBinaryReader{}
		buff := bytes.Buffer{}
		t.Run("test int", func(t *testing.T) {
			var x int = 33000
			buff.Write(writer.write(x).([]byte))
			val := reader.read(&buff).(int64)
			assert.Equal(t, x, int(val))
		})
		t.Run("test int32", func(t *testing.T) {
			var x int32 = 33000
			buff.Write(writer.write(x).([]byte))
			assert.Equal(t, x, reader.read(&buff))
		})
		t.Run("test long", func(t *testing.T) {
			var x int64 = 2147483648
			buff.Write(writer.write(x).([]byte))
			assert.Equal(t, x, reader.read(&buff))
		})
		t.Run("test string", func(t *testing.T) {
			var x = "serialize this!"
			buff.Write(writer.write(x).([]byte))
			assert.Equal(t, x, reader.read(&buff))
		})
		t.Run("test string", func(t *testing.T) {
			var x, _ = uuid.Parse("41d2e28a-20a4-4ab0-b379-d810dede3786")
			buff.Write(writer.write(x).([]byte))
			assert.Equal(t, x, reader.read(&buff))
		})
	})

	t.Run("test nested types", func(t *testing.T) {
		writer := graphBinaryWriter{}
		reader := graphBinaryReader{}
		buff := bytes.Buffer{}
		t.Run("test map", func(t *testing.T) {
			var x int32 = 666
			var m = map[interface{}]interface{}{
				"marko": x,
				"noone": "blah",
			}
			buff.Write(writer.write(m).([]byte))
			assert.Equal(t, m, reader.read(&buff))
		})
	})

	t.Run("test long", func(t *testing.T) {
		var x = 100
		var y int32 = 100
		var z int64 = 100
		var s = "serialize this!"
		var u, _ = uuid.Parse("41d2e28a-20a4-4ab0-b379-d810dede3786")
		var a int64 = 666
		var m = map[interface{}]interface{}{
			"marko": a,
			"noone": "blah",
		}
		writer := graphBinaryWriter{}
		reader := graphBinaryReader{}
		buff := bytes.Buffer{}
		buff.Write(writer.write(s).([]byte))
		fmt.Println(reader.read(&buff))
		buff.Write(writer.write(x).([]byte))
		fmt.Println(reader.read(&buff))
		buff.Write(writer.write(y).([]byte))
		fmt.Println(reader.read(&buff))
		buff.Write(writer.write(u).([]byte))
		fmt.Println(reader.read(&buff))
		buff.Write(writer.write(z).([]byte))
		fmt.Println(reader.read(&buff))
		buff.Write(writer.write(m).([]byte))
		out := reader.read(&buff)
		fmt.Println(reflect.DeepEqual(m, out))
		assert.Equal(t, m, out)
	})
}
