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
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test(t *testing.T) {
	t.Run("Test dataReceived nil message", func(t *testing.T) {
		protocol := newGremlinServerWSProtocol(newLogHandler(&defaultLogger{}, Info, language.English))
		statusCode, err := protocol.dataReceived(nil, map[string]ResultSet{})
		assert.Equal(t, uint16(0), statusCode)
		assert.Equal(t, err, errors.New("malformed ws or wss URL"))
	})

	t.Run("Test dataReceived actual message", func(t *testing.T) {
		protocol := newGremlinServerWSProtocol(newLogHandler(&defaultLogger{}, Info, language.English))
		data := map[interface{}]interface{}{"data_key": "data_val"}
		var u, _ = uuid.Parse("41d2e28a-20a4-4ab0-b379-d810dede3786")
		resultSets := map[string]ResultSet{}
		resultSet := newChannelResultSet(u.String())
		resultSets[u.String()] = resultSet
		testResponseStatusOK := response{
			requestID: u,
			responseStatus: responseStatus{
				code:       http.StatusOK,
				message:    "",
				attributes: map[interface{}]interface{}{"attr_key": "attr_val"},
			},
			responseResult: responseResult{
				data: data,
				meta: map[interface{}]interface{}{"meta_key": "meta_val"},
			},
		}
		testResponseStatusNoContent := response{
			requestID: u,
			responseStatus: responseStatus{
				code:       http.StatusNoContent,
				message:    "",
				attributes: map[interface{}]interface{}{"attr_key": "attr_val"},
			},
			responseResult: responseResult{},
		}
		serializer := graphBinarySerializer{}
		message, err := serializer.serializeResponseMessage(&testResponseStatusOK)
		assert.Nil(t, err)
		code, err := protocol.dataReceived(message, resultSets)
		assert.Nil(t, err)
		assert.Equal(t, uint16(http.StatusOK), code)

		message, err = serializer.serializeResponseMessage(&testResponseStatusNoContent)
		assert.Nil(t, err)
		code, err = protocol.dataReceived(message, resultSets)
		assert.Nil(t, err)
		assert.Equal(t, uint16(http.StatusNoContent), code)

		result1 := resultSet.one()
		assert.Equal(t, fmt.Sprintf("%v", data), result1.AsString())
		result2 := resultSet.one()
		assert.Equal(t, fmt.Sprintf("%v", make([]interface{}, 0)), result2.AsString())
	})

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
