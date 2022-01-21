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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSerializer(t *testing.T) {
	t.Run("test serialize and deserialize request message", func(t *testing.T) {
		var u, _ = uuid.Parse("41d2e28a-20a4-4ab0-b379-d810dede3786")
		testRequest := Request{
			RequestID: u,
			Op:        "eval",
			Processor: "traversal",
			Args:      map[interface{}]interface{}{"test_key": "test_val"},
		}
		serializer := GraphBinarySerializer{}
		serialized, _ := serializer.SerializeMessage(&testRequest)
		deserialized, _ := serializer.deserializeRequestMessage(&serialized)
		assert.Equal(t, testRequest, deserialized)
	})

	t.Run("test serialize and deserialize response message", func(t *testing.T) {
		var u, _ = uuid.Parse("41d2e28a-20a4-4ab0-b379-d810dede3786")
		testResponse := Response{
			RequestID: u,
			ResponseStatus: ResponseStatus{
				Code:       200,
				Message:    "",
				Attributes: map[interface{}]interface{}{"attr_key": "attr_val"},
			},
			ResponseResult: ResponseResult{
				Data: map[interface{}]interface{}{"data_key": "data_val"},
				Meta: map[interface{}]interface{}{"meta_key": "meta_val"},
			},
		}
		serializer := GraphBinarySerializer{}
		serialized, _ := serializer.serializeResponseMessage(&testResponse)
		deserialized, _ := serializer.DeserializeMessage(&serialized)
		assert.Equal(t, testResponse, deserialized)
	})
}
