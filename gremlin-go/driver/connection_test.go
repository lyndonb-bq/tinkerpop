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
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

const personLabel = "Person"
const testLabel = "Test"
const nameKey = "name"
const failureHost = "invalidHost"
const failurePort = 9999

func dropGraph(t *testing.T, g *GraphTraversalSource) {
	// Drop vertices that were added.
	_, promise, err := g.V().Drop().Iterate()
	assert.Nil(t, err)
	assert.NotNil(t, promise)
	<-promise
}

func getTestNames() []string {
	return []string{"Lyndon", "Yang", "Simon", "Rithin", "Alexey", "Valentyn"}
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

func initializeGraph(t *testing.T, testHost string, testPort int) *GraphTraversalSource {
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

	return g
}

func resetGraph(t *testing.T, g *GraphTraversalSource) {
	// Drop the graph and check that it is empty.
	dropGraph(t, g)
	readCount(t, g, "", 0)
	readCount(t, g, testLabel, 0)
	readCount(t, g, personLabel, 0)
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

func readUsingAnonymousTraversal(t *testing.T, g *GraphTraversalSource) {
	results, err := g.V().Fold().
		Project(testLabel, personLabel).
		By(T__.Unfold().HasLabel(testLabel).Count()).
		By(T__.Unfold().HasLabel(personLabel).Count()).
		ToList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))
	resultMap := results[0].GetInterface().(map[interface{}]interface{})
	assert.Equal(t, int64(0), resultMap[testLabel])
	assert.Equal(t, int64(len(getTestNames())), resultMap[personLabel])
}

func readWithNextAndHasNext(t *testing.T, g *GraphTraversalSource) {
	traversal := g.V().HasLabel(personLabel).Properties(nameKey)
	var names []string
	for i := 0; i < len(getTestNames()); i++ {
		hasN, err := traversal.HasNext()
		assert.Nil(t, err)
		assert.True(t, hasN)
		res, err := traversal.Next()
		assert.Nil(t, err)
		assert.NotNil(t, res)
		vp, err := res.GetVertexProperty()
		assert.Nil(t, err)
		names = append(names, vp.value.(string))
	}
	hasN, _ := traversal.HasNext()
	assert.False(t, hasN)
	// Check for Next error when no more elements left
	res, err := traversal.Next()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.True(t, sortAndCompareTwoStringSlices(names, getTestNames()))
}

func getEnvOrDefaultString(key string, defaultValue string) string {
	// Missing value is returned as "".
	value := os.Getenv(key)
	if len(value) != 0 {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultBool(key string, defaultValue bool) bool {
	value := getEnvOrDefaultString(key, "")
	if len(value) != 0 {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	value := getEnvOrDefaultString(key, "")
	if len(value) != 0 {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func TestConnection(t *testing.T) {
	testHost := getEnvOrDefaultString("GREMLIN_SERVER_HOSTNAME", "localhost")
	testPort := getEnvOrDefaultInt("GREMLIN_SERVER_PORT", 8182)
	runIntegration := getEnvOrDefaultBool("RUN_INTEGRATION_TESTS", true)

	t.Run("Test DriverRemoteConnection GraphTraversal", func(t *testing.T) {
		if runIntegration {
			// Initialize graph
			g := initializeGraph(t, testHost, testPort)
			// Read test data out of the graph and check that it is correct.
			readTestDataVertexProperties(t, g)
			readTestDataValues(t, g)
			// Reset Graph
			resetGraph(t, g)
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
			assert.Equal(t, "0", result.GetString())
			err = connection.close()
			assert.Nil(t, err)
		}
	})

	t.Run("Test connection.close() failure", func(t *testing.T) {
		connection, err := createConnection(testHost, testPort, newLogHandler(&defaultLogger{}, Info, language.English))
		err = connection.close()
		assert.Nil(t, err)
		err = connection.close()
		assert.NotNil(t, err)
		err = connection.close()
		assert.NotNil(t, err)
	})

	t.Run("Test createConnection failure", func(t *testing.T) {
		_, err := createConnection(failureHost, failurePort, newLogHandler(&defaultLogger{}, Info, language.English))
		assert.NotNil(t, err)
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
			assert.Equal(t, "0", result.GetString())
			err = client.Close()
			assert.Nil(t, err)
		}
	})

	t.Run("Test Traversal. Next and HasNext", func(t *testing.T) {
		if runIntegration {
			g := initializeGraph(t, testHost, testPort)
			readWithNextAndHasNext(t, g)
			resetGraph(t, g)
		}
	})

	t.Run("Test anonymousTraversal", func(t *testing.T) {
		if runIntegration {
			g := initializeGraph(t, testHost, testPort)
			readUsingAnonymousTraversal(t, g)
			resetGraph(t, g)
		}
	})

	t.Run("Test Traversal.ToString fail", func(t *testing.T) {
		anonTrav := T__.Unfold().HasLabel(testLabel)
		slice, err := anonTrav.ToList()
		assert.Nil(t, slice)
		assert.NotNil(t, err)
	})

	t.Run("Test Traversal.Iterate fail", func(t *testing.T) {
		anonTrav := T__.Unfold().HasLabel(testLabel)
		traversal, channel, err := anonTrav.Iterate()
		assert.Nil(t, traversal)
		assert.Nil(t, channel)
		assert.NotNil(t, err)
	})
}
