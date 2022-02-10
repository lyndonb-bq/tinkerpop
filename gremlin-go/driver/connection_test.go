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
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

const runIntegration = true
const testHost string = "localhost"
const testPort int = 8182

func TestConnection(t *testing.T) {
	t.Run("Test g.V().count()", func(t *testing.T) {
		if runIntegration {
			remote, err := NewDriverRemoteConnection(testHost, testPort)
			assert.Nil(t, err)
			assert.NotNil(t, remote)
			g := _Traversal().WithRemote(remote)

			// Drop vertices that were added.
			_, err = g.V().Drop().Iterate()
			assert.Nil(t, err)

			// Read count again, should be 0.
			results, err := g.V().Count().ToList()
			assert.Nil(t, err)
			assert.NotNil(t, results)
			assert.Equal(t, 1, len(results))
			var count int32
			count, err = results[0].GetInt32()
			assert.Nil(t, err)
			assert.Equal(t, int32(0), count)
			time.Sleep(50 * time.Millisecond)

			// Read count again, should be 0, same as above.
			results, err = g.V().HasLabel("Test").Count().ToList()
			assert.Nil(t, err)
			assert.NotNil(t, results)
			assert.Equal(t, 1, len(results))
			count, err = results[0].GetInt32()
			assert.Nil(t, err)
			assert.Equal(t, int32(0), count)
			time.Sleep(50 * time.Millisecond)

			// Add 5 vertices.
			_, err = g.
				AddV("person").Property("name", "Lyndon").
				AddV("person").Property("name", "Simon").
				AddV("person").Property("name", "Yang").
				AddV("person").Property("name", "Rithin").
				AddV("person").Property("name", "Alexey").
				Iterate()
			assert.Nil(t, err)
			time.Sleep(50 * time.Millisecond)

			// Read count again, should be 5.
			results, err = g.V().Count().ToList()
			assert.Nil(t, err)
			assert.NotNil(t, results)
			assert.Equal(t, 1, len(results))
			count, err = results[0].GetInt32()
			assert.Nil(t, err)
			assert.Equal(t, int32(5), count)
			time.Sleep(50 * time.Millisecond)

			// Drop vertices that were added.
			_, err = g.V().Drop().Iterate()
			assert.Nil(t, err)
			time.Sleep(50 * time.Millisecond)
		}
	})

	t.Run("Test addV()", func(t *testing.T) {
		if runIntegration {
			remote, err := NewDriverRemoteConnection(testHost, testPort)
			g := _Traversal().WithRemote(remote)

			// Add 5 vertices.
			_, err = g.
				AddV("person").Property("name", "Lyndon").
				AddV("person").Property("name", "Simon").
				AddV("person").Property("name", "Yang").
				AddV("person").Property("name", "Rithin").
				AddV("person").Property("name", "Alexey").
				Iterate()
			assert.Nil(t, err)
		}
	})

	t.Run("Test createConnection", func(t *testing.T) {
		if runIntegration {
			connection, err := createConnection(testHost, testPort, newLogHandler(&defaultLogger{}, Info, language.English))
			assert.Nil(t, err)
			assert.NotNil(t, connection)
			err = connection.close()
			assert.Nil(t, err)
		}
	})

	t.Run("Test connection.write()", func(t *testing.T) {
		if runIntegration {
			connection, err := createConnection(testHost, testPort, newLogHandler(&defaultLogger{}, Info, language.English))
			assert.Nil(t, err)
			assert.NotNil(t, connection)
			request := makeStringRequest("g.V().count()")
			resultSet, err := connection.write(&request)
			assert.Nil(t, err)
			assert.NotNil(t, resultSet)
			result := resultSet.one()
			assert.NotNil(t, result)
			assert.Equal(t, "[0]", result.GetString())
			err = connection.close()
			assert.Nil(t, err)
		}
	})

	t.Run("Test client.submit()", func(t *testing.T) {
		if runIntegration {
			client, err := NewClient(testHost, testPort)
			assert.Nil(t, err)
			assert.NotNil(t, client)
			resultSet, err := client.Submit("g.V().count()")
			assert.Nil(t, err)
			assert.NotNil(t, resultSet)
			result := resultSet.one()
			assert.NotNil(t, result)
			assert.Equal(t, "[0]", result.GetString())
			err = client.Close()
			assert.Nil(t, err)
		}
	})
}
