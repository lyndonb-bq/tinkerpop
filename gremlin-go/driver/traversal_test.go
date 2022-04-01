/*
Licensed to the Apache Software Foundation (ASF) Under one
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
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraversal(t *testing.T) {
	t.Run("Test clone traversal", func(t *testing.T) {
		g := NewGraphTraversalSource(&Graph{}, &TraversalStrategies{}, newBytecode(nil), nil)
		original := g.V().Out("created")
		clone := original.Clone().Out("knows")
		cloneClone := clone.Clone().Out("created")

		assert.Equal(t, 2, len(original.bytecode.stepInstructions))
		assert.Equal(t, 3, len(clone.bytecode.stepInstructions))
		assert.Equal(t, 4, len(cloneClone.bytecode.stepInstructions))

		original.Has("person", "name", "marko")
		clone.V().Out()

		assert.Equal(t, 3, len(original.bytecode.stepInstructions))
		assert.Equal(t, 5, len(clone.bytecode.stepInstructions))
		assert.Equal(t, 4, len(cloneClone.bytecode.stepInstructions))
	})

	t.Run("Test traversal with bindings", func(t *testing.T) {
		g := NewGraphTraversalSource(&Graph{}, &TraversalStrategies{}, newBytecode(nil), nil)
		bytecode := g.V((&Bindings{}).Of("a", []int32{1, 2, 3})).
			Out((&Bindings{}).Of("b", "created")).
			Where(T__.In((&Bindings{}).Of("c", "created"), (&Bindings{}).Of("d", "knows")).
				Count().Is((&Bindings{}).Of("e", P.Gt(2)))).bytecode
		assert.Equal(t, 5, len(bytecode.bindings))
		assert.Equal(t, []int32{1, 2, 3}, bytecode.bindings["a"])
		assert.Equal(t, "created", bytecode.bindings["b"])
		assert.Equal(t, "created", bytecode.bindings["c"])
		assert.Equal(t, "knows", bytecode.bindings["d"])
		assert.Equal(t, P.Gt(2), bytecode.bindings["e"])
		assert.Equal(t, &Binding{
			Key:   "b",
			Value: "created",
		}, bytecode.stepInstructions[1].arguments[0])
		assert.Equal(t, "binding[b=created]", bytecode.stepInstructions[1].arguments[0].(*Binding).String())
	})

	t.Run("Test Transaction commit", func(t *testing.T) {
		// skipTestsIfNotEnabled(t, integrationTestSuiteName, testNoAuthEnable)
		// Start a transaction traversal.
		remote := newConnection(t)
		g := Traversal_().WithRemote(remote)
		startCount := getCount(t, g)
		tx := g.Tx()

		// Except transaction to not be open until begin is called.
		assert.False(t, tx.IsOpen())
		gtx, _ := tx.Begin()
		assert.True(t, tx.IsOpen())

		addV(t, gtx)
		addV(t, gtx)
		assert.Equal(t, startCount, getCount(t, g))
		assert.Equal(t, startCount+2, getCount(t, gtx))

		// Commit the transaction, this should close it.
		// Our vertex count outside the transaction should be 2 + the start count.
		_, err := tx.Commit()
		assert.Nil(t, err)

		assert.False(t, tx.IsOpen())
		// todo: assert.Equal(t, startCount+2, getCount(t, g))

		// dropGraphCheckCount(t, g)
		verifyGtxClosed(t, gtx)
	})
	t.Run("Test Transaction rollback", func(t *testing.T) {
		// skipTestsIfNotEnabled(t, integrationTestSuiteName, testNoAuthEnable)
		// Start a transaction traversal.
		remote := newConnection(t)
		g := Traversal_().WithRemote(remote)
		startCount := getCount(t, g)
		tx := g.Tx()

		// Except transaction to not be open until begin is called.
		assert.False(t, tx.IsOpen())
		gtx, _ := tx.Begin()
		assert.True(t, tx.IsOpen())

		addV(t, gtx)
		addV(t, gtx)
		assert.Equal(t, startCount, getCount(t, g))
		assert.Equal(t, startCount+2, getCount(t, gtx))

		// Rollback the transaction, this should close it.
		// Our vertex count outside the transaction should be the start count.
		_, err := tx.Rollback()
		assert.Nil(t, err)

		assert.False(t, tx.IsOpen())
		assert.Equal(t, startCount, getCount(t, g))
		assert.Equal(t, closed, gtx.remoteConnection.client.connection.state)

		// dropGraphCheckCount(t, g)
		verifyGtxClosed(t, gtx)
	})

}

func newConnection(t *testing.T) *DriverRemoteConnection {
	testNoAuthWithAliasUrl := getEnvOrDefaultString("GREMLIN_SERVER_URL", "ws://localhost:8182/gremlin")
	testNoAuthWithAliasAuthInfo := &AuthInfo{}
	testNoAuthWithAliasTlsConfig := &tls.Config{}

	remote, err := NewDriverRemoteConnection(testNoAuthWithAliasUrl,
		func(settings *DriverRemoteConnectionSettings) {
			settings.TlsConfig = testNoAuthWithAliasTlsConfig
			settings.AuthInfo = testNoAuthWithAliasAuthInfo
			settings.TraversalSource = "gtx"
		})
	assert.Nil(t, err)
	assert.NotNil(t, remote)
	return remote
}

func addV(t *testing.T, g *GraphTraversalSource) {
	_, promise, err := g.AddV("person").Property("name", "lyndon").Iterate()
	assert.Nil(t, err)
	assert.Nil(t, <-promise)
}

/*func addNodeValidateTransactionState(t *testing.T, g, gAddTo *GraphTraversalSource,
	gStartCount, gAddToStartCount int32, txVerifyList ...*transaction) {
	// Add a single node to g_add_to, but not g.
	// Check that vertex count in g is gStartCount and vertex count in gAddTo is gAddToStartCount + 1.
	_, promise, err := gAddTo.AddV("person").Property("name", "lyndon").Iterate()
	assert.Nil(t, err)
	assert.Nil(t, <-promise)
	fmt.Printf("--addNodeValidateTransactionState %v-%v\n", gStartCount, getCount(t, g))
	assert.Equal(t, gAddToStartCount+1, getCount(t, gAddTo))
	assert.Equal(t, gStartCount, getCount(t, g))
	verifyTxState(t, txVerifyList, true)
}*/

func verifyTxState(t *testing.T, gtxList []*transaction, value bool) {
	for _, tx := range gtxList {
		assert.Equal(t, value, tx.IsOpen())
	}
}

func dropGraphCheckCount(t *testing.T, g *GraphTraversalSource) {
	//time.Sleep(200 * time.Millisecond)

	g1, _ := g.V().ToList()

	//time.Sleep(200 * time.Millisecond)

	dropGraph(t, g)

	g2, err := g.V().ToList()
	count, err := g.V().Count().ToList()
	println("-- XXX:", g1, g2, err, count)

	assert.Equal(t, int32(0), getCount(t, g))
}

func verifyGtxClosed(t *testing.T, gtx *GraphTraversalSource) {
	// Attempt to add an additional vertex to the transaction. This should return an error since it
	// has been closed.
	_, _, err := gtx.AddV("failure").Iterate()
	assert.NotNil(t, err)
}

func getCount(t *testing.T, g *GraphTraversalSource) int32 {
	count, err := g.V().Count().ToList()
	assert.Nil(t, err)
	assert.NotNil(t, count)
	assert.Equal(t, 1, len(count))
	val, err := count[0].GetInt32()
	assert.Nil(t, err)
	return val
}
