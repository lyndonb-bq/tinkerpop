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
	"fmt"
	"github.com/cucumber/godog"
	"github.com/lyndonb-bq/tinkerpop/gremlin-go/driver"
)

type TinkerPopWorld struct {
	scenario     *godog.Scenario
	g            *gremlingo.GraphTraversalSource
	graphName    string
	traversal    *gremlingo.GraphTraversal
	result       []*gremlingo.Result
	graphDataMap map[string]DataGraph
	parameters   map[string]interface{}
}

type DataGraph struct {
	name       string
	connection *gremlingo.DriverRemoteConnection
	vertices   map[string]*gremlingo.Vertex
	edges      map[string]*gremlingo.Edge
}

const scenarioHost string = "localhost"
const scenarioPort int = 8182

func NewTinkerPopWorld() *TinkerPopWorld {
	return &TinkerPopWorld{
		scenario:     nil,
		g:            nil,
		graphName:    "",
		traversal:    nil,
		result:       nil,
		graphDataMap: make(map[string]DataGraph),
		parameters:   make(map[string]interface{}),
	}
}

func getGraphNames() []string {
	return []string{"modern", "classic", "crew", "grateful", "sink", "empty"}
}

func (t *TinkerPopWorld) getDataGraph(name string) *DataGraph {
	if val, ok := t.graphDataMap[name]; ok {
		return &val
	} else {
		return nil
	}
}

// May be redundant, currently not needed
func (t *TinkerPopWorld) loadEmptyDataGraph() {
	connection, _ := gremlingo.NewDriverRemoteConnection(scenarioHost, scenarioPort, func(settings *gremlingo.DriverRemoteConnectionSettings) {
		settings.TraversalSource = "ggraph"
	})
	t.graphDataMap["empty"] = DataGraph{connection: connection}
}

func (t *TinkerPopWorld) reloadEmptyData() {
	graphData := t.getDataGraph("empty")
	g := gremlingo.Traversal_().WithRemote(graphData.connection)
	graphData.vertices = getVertices(g)
	graphData.edges = getEdges(g)
}

func (t *TinkerPopWorld) cleanEmptyDataGraph() error {
	connection := t.graphDataMap["empty"].connection
	g := gremlingo.Traversal_().WithRemote(connection)
	_, future, err := g.V().Drop().Iterate()
	if err != nil {
		return err
	}
	<-future
	return nil
}

func (t *TinkerPopWorld) loadAllDataGraph() {
	for _, name := range getGraphNames() {
		fmt.Println(t.graphDataMap)
		if name == "empty" {
			connection, _ := gremlingo.NewDriverRemoteConnection(scenarioHost, scenarioPort, func(settings *gremlingo.DriverRemoteConnectionSettings) {
				settings.TraversalSource = "ggraph"
			})
			t.graphDataMap["empty"] = DataGraph{connection: connection}
		} else {
			connection, _ := gremlingo.NewDriverRemoteConnection(scenarioHost, scenarioPort, func(settings *gremlingo.DriverRemoteConnectionSettings) {
				settings.TraversalSource = "g" + name
			})
			g := gremlingo.Traversal_().WithRemote(connection)
			t.graphDataMap["hi"] = DataGraph{}
			t.graphDataMap[name] = DataGraph{
				name:       name,
				connection: connection,
				vertices:   getVertices(g),
				edges:      getEdges(g),
			}
		}
	}
}

// TODO implement after .Next() implementation
func getVertices(g *gremlingo.GraphTraversalSource) map[string]*gremlingo.Vertex {
	// need .Next() at the end
	//__ := gremlingo.AnonTrav__{}
	//g.V().Group().By("name").By(__.Tail())
	vertexMap := make(map[string]*gremlingo.Vertex)
	return vertexMap
}

// TODO implement after .Next() implementation
func getEdges(g *gremlingo.GraphTraversalSource) map[string]*gremlingo.Edge {
	edgeMap := make(map[string]*gremlingo.Edge)
	return edgeMap
}
