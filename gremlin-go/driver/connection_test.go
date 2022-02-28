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
	"github.com/stretchr/testify/assert"
	"gitlab.com/avarf/getenvs"
	"golang.org/x/text/language"
	"reflect"
	"sort"
	"testing"
)

const personLabel = "Person"
const testLabel = "Test"
const nameKey = "name"

func dropGraph(t *testing.T, g *GraphTraversalSource) {
	// Drop vertices that were added.
	_, promise, err := g.V().Drop().Iterate()
	assert.Nil(t, err)
	assert.NotNil(t, promise)
	<-promise
}

func getTestNames() []string {
	return []string{"Lyndon", "Yang", "Simon", "Rithin", "Alexey", "Stephen"}
}

func addTestData(t *testing.T, g *GraphTraversalSource) {
	testNames := getTestNames()

	// Add vertices to traversal.
	var traversal *GraphTraversal
	for _, name := range testNames {
		if traversal == nil {
			traversal = g.AddV(personLabel).Property(nameKey, name)
		} else {
			traversal = traversal.AddV(personLabel).Property(nameKey, name)
		}
	}

	// Commit traversal.
	_, promise, err := traversal.Iterate()
	assert.Nil(t, err)
	<-promise
}

func readTestDataVertexProperties(t *testing.T, g *GraphTraversalSource) {
	// Read names from graph
	var names []string
	results, err := g.V().HasLabel(personLabel).Properties(nameKey).ToList()
	for _, result := range results {
		vp, err := result.GetVertexProperty()
		assert.Nil(t, err)
		names = append(names, vp.value.(string))
	}
	assert.Nil(t, err)
	assert.NotNil(t, names)
	assert.True(t, sortAndCompareTwoStringSlices(names, getTestNames()))
}

func readTestDataValues(t *testing.T, g *GraphTraversalSource) {
	// Read names from graph
	var names []string
	results, err := g.V().HasLabel(personLabel).Values(nameKey).ToList()
	for _, result := range results {
		names = append(names, result.GetString())
	}
	assert.Nil(t, err)
	assert.NotNil(t, names)
	assert.True(t, sortAndCompareTwoStringSlices(names, getTestNames()))
}

func readCount(t *testing.T, g *GraphTraversalSource, label string, expected int) {
	// Generate traversal.
	var traversal *GraphTraversal
	if label != "" {
		traversal = g.V().HasLabel(label).Count()
	} else {
		traversal = g.V().Count()
	}

	// Get results from traversal.
	results, err := traversal.ToList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))

	// Read count from results.
	var count int32
	count, err = results[0].GetInt32()
	assert.Nil(t, err)

	// Check count.
	assert.Equal(t, int32(expected), count)
}

func sortAndCompareTwoStringSlices(s1 []string, s2 []string) bool {
	sort.Slice(s1, func(i, j int) bool {
		return s1[i] < s1[j]
	})
	sort.Slice(s2, func(i, j int) bool {
		return s2[i] < s2[j]
	})
	return reflect.DeepEqual(s1, s2)
}

func TestConnection(t *testing.T) {
	testHost := getenvs.GetEnvString("GREMLIN_SERVER_HOSTNAME", "localhost")
	testPort, _ := getenvs.GetEnvInt("GREMLIN_SERVER_PORT", 8182)
	runIntegration, _ := getenvs.GetEnvBool("RUN_INTEGRATION_TESTS", true)

	t.Run("Test DriverRemoteConnection GraphTraversal", func(t *testing.T) {
		if runIntegration {
			remote, err := NewDriverRemoteConnection(testHost, testPort)
			assert.Nil(t, err)
			assert.NotNil(t, remote)
			g := Traversal_().WithRemote(remote)

			// Drop the graph and check that it is empty.
			dropGraph(t, g)
			readCount(t, g, "", 0)
			readCount(t, g, testLabel, 0)
			readCount(t, g, personLabel, 0)

			// Add data and check that the size of the graph is correct.
			addTestData(t, g)
			readCount(t, g, "", len(getTestNames()))
			readCount(t, g, testLabel, 0)
			readCount(t, g, personLabel, len(getTestNames()))

			// Read test data out of the graph and check that it is correct.
			readTestDataVertexProperties(t, g)
			readTestDataValues(t, g)

			// Drop the graph and check that it is empty.
			dropGraph(t, g)
			readCount(t, g, "", 0)
			readCount(t, g, testLabel, 0)
			readCount(t, g, personLabel, 0)
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

	t.Run("Test DriverRemoteConnection Next and HasNext", func(t *testing.T) {
		if runIntegration {
			// Add data
			remote, err := NewDriverRemoteConnection(testHost, testPort)
			assert.Nil(t, err)
			assert.NotNil(t, remote)
			g := Traversal_().WithRemote(remote)
			dropGraph(t, g)
			addTestData(t, g)

			// Run traversal and test Next/HasNext calls
			traversal := g.V().HasLabel(personLabel).Properties(nameKey)
			var names []string
			for i := 0; i < len(getTestNames()); i++ {
				hasN, _ := traversal.HasNext()
				assert.True(t, hasN)
				res, _ := traversal.Next()
				vp, _ := res.GetVertexProperty()
				names = append(names, vp.value.(string))
			}
			hasN, _ := traversal.HasNext()
			assert.False(t, hasN)
			assert.True(t, sortAndCompareTwoStringSlices(names, getTestNames()))
		}
	})
}
