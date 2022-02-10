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

			_, err = g.V().Drop().Iterate()
			assert.Nil(t, err)

			// Read count, expect there to be 0 vertices.
			results, err := g.V().Count().ToList()
			assert.Nil(t, err)
			assert.NotNil(t, results)
			assert.Equal(t, 1, len(results))
			var count int32
			count, err = results[0].GetInt32()
			assert.Nil(t, err)
			assert.Equal(t, int32(0), count)
			time.Sleep(50 * time.Millisecond)

			// Read count, expect there to be 0 vertices.
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
			// 157 -> 179?
			//

			// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 1 119 73 39 145 98 71 31 129 112 123 151 184 194 29 229     0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 2  0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 3 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]
			// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 85 52 180 61 184 76 64 221 160 191 245 217 183 45 236 96    0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 13 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 6 76 121 110 100 111 110 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 5 83 105 109 111 110 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 4 89 97 110 103 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 6 82 105 116 104 105 110 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 6 65 108 101 120 101 121 0 0 0 4 110 111 110 101 0 0 0 0 0 0 0 0 3 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]
			// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 183 175 134 165 139 105 79 112 176 216 122 93 254 176 241 4 0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 15 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 6 76 121 110 100 111 110 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 5 83 105 109 111 110 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 4 89 97 110 103 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 6 82 105 116 104 105 110 0 0 0 4 97 100 100 86 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 6 112 101 114 115 111 110 0 0 0 8 112 114 111 112 101 114 116 121 0 0 0 1 0 0 0 1 9 0 0 0 0 2 3 0 0 0 0 4 110 97 109 101 3 0 0 0 0 6 65 108 101 120 101 121 0 0 0 4 110 111 110 101 0 0 0 0 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0]1624 <nil>
			results, err = g.V().Count().ToList()
			assert.Nil(t, err)
			assert.NotNil(t, results)
			assert.Equal(t, 1, len(results))
			count, err = results[0].GetInt32()
			assert.Nil(t, err)
			assert.Equal(t, int32(5), count)

			time.Sleep(50 * time.Millisecond)
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
			connection, err := createConnection(newLogHandler(&defaultLogger{}, Info, language.English), testHost, testPort)
			assert.Nil(t, err)
			assert.NotNil(t, connection)
			err = connection.close()
			assert.Nil(t, err)
		}
	})

	t.Run("Test write", func(t *testing.T) {
		if runIntegration {
			connection, err := createConnection(newLogHandler(&defaultLogger{}, Info, language.English), testHost, testPort)
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

	t.Run("Test client submit", func(t *testing.T) {
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
