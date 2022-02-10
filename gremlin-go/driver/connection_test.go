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

const runIntegration = false
const testHost string = "localhost"
const testPort int = 8182

func TestConnection(t *testing.T) {
	t.Run("Test connect", func(t *testing.T) {
		if runIntegration {
			connection := connection{testHost, testPort, Gorilla, newLogHandler(&defaultLogger{}, Info, language.English), nil, nil, nil}
			err := connection.connect()
			assert.Nil(t, err)
		}
	})

	t.Run("Test write", func(t *testing.T) {
		if runIntegration {
			connection := connection{testHost, testPort, Gorilla, newLogHandler(&defaultLogger{}, Info, language.English), nil, nil, nil}
			err := connection.connect()
			assert.Nil(t, err)
			request := makeStringRequest("g.V().count()")
			resultSet, err := connection.write(&request)
			assert.Nil(t, err)
			assert.NotNil(t, resultSet)
			result := resultSet.one()
			assert.NotNil(t, result)
			assert.Equal(t, "[0]", result.GetString())
		}
	})

	t.Run("Test client submit", func(t *testing.T) {
		if runIntegration {
			connection := connection{testHost, testPort, Gorilla, newLogHandler(&defaultLogger{}, Info, language.English), nil, nil, nil}
			err := connection.connect()
			assert.Nil(t, err)
			client := NewClient(testHost, testPort)
			resultSet, err := client.Submit("g.V().count()")
			assert.Nil(t, err)
			assert.NotNil(t, resultSet)
			result := resultSet.one()
			assert.NotNil(t, result)
			assert.Equal(t, "[0]", result.GetString())
		}
	})
}
