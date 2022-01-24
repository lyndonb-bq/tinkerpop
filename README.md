<!--

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

-->

# Go Gremlin Language Variant

[Apache TinkerPop™][tk] is a graph computing framework for both graph databases (OLTP) and graph analytic systems
(OLAP). [Gremlin][gremlin] is the graph traversal language of TinkerPop. It can be described as a functional,
data-flow language that enables users to succinctly express complex traversals on (or queries of) their application's
property graph.

<!--
TODO: Add gremlin specific details for following paragraph.
Python readme example:
Gremlin-Python implements Gremlin within the Python language and can be used on any Python virtual machine including 
the popular CPython machine. Python’s syntax has the same constructs as Java including "dot notation" for function 
chaining (a.b.c), round bracket function arguments (a(b,c)), and support for global namespaces (a(b()) vs a(__.b())). 
As such, anyone familiar with Gremlin-Java will immediately be able to work with Gremlin-Python. 
Moreover, there are a few added constructs to Gremlin-Python that make traversals a bit more succinct.
-->
Gremlin-Go implements Gremlin within the Go language.

Gremlin-Go is designed to connect to a "server" that is hosting a TinkerPop-enabled graph system. That "server"
could be [Gremlin Server][gs] or a [remote Gremlin provider][rgp] that exposes protocols by which Gremlin-Javascript
can connect.

A typical connection to a server running on "localhost" that supports the Gremlin Server protocol using websockets
looks like this:

<!--
TODO: Add Go code example of connection to server.
Javascript readme example:
```javascript
const gremlin = require('gremlin');
const traversal = gremlin.process.AnonymousTraversalSource.traversal;
const DriverRemoteConnection = gremlin.driver.DriverRemoteConnection;

const g = traversal().withRemote(new DriverRemoteConnection('ws://localhost:8182/gremlin'));
```
-->

Once "g" has been created using a connection, it is then possible to start writing Gremlin traversals to query the
remote graph:

<!--
TODO: Add Go code example of a Gremlin traversal query.
Javascript readme example:
```javascript
const gremlin = require('gremlin');
const traversal = gremlin.process.AnonymousTraversalSource.traversal;
const DriverRemoteConnection = gremlin.driver.DriverRemoteConnection;

const g = traversal().withRemote(new DriverRemoteConnection('ws://localhost:8182/gremlin'));
```
-->

# The following material is currently Work-in-progress: 

## Sample Traversals

<!--
TODO: Add Go specific changes to a paragraph such as this
javascript example:
The Gremlin language allows users to write highly expressive graph traversals and has a broad list of functions that 
cover a wide body of features. The [Reference Documentation][steps] describes these functions and other aspects of the 
TinkerPop ecosystem including some specifics on [Gremlin in Javascript][docs] itself. Most of the examples found in the 
documentation use Groovy language syntax in the [Gremlin Console][console]. For the most part, these examples
should generally translate to Javascript with [little modification][differences]. Given the strong correspondence 
between canonical Gremlin in Java and its variants like Javascript, there is a limited amount of Javascript-specific 
documentation and examples. This strong correspondence among variants ensures that the general Gremlin reference 
documentation is applicable to all variants and that users moving between development languages can easily adopt the 
Gremlin variant for that language.
-->

### Create Vertex
<!--
TODO: Add Go code to create a vertex.
javascript example:
```javascript
/* if we want to assign our own ID and properties to this vertex */
const { t: { id } } = gremlin.process;
const { cardinality: { single } } = gremlin.process;

/**
 * Create a new vertex with Id, Label and properties
 * @param {String,Number} vertexId Vertex Id (assuming the graph database allows id assignment)
 * @param {String} vlabel Vertex Label
 */
const createVertex = async (vertexId, vlabel) => {
  const vertex = await g.addV(vlabel)
    .property(id, vertexId)
    .property(single, 'name', 'Apache')
    .property('lastname', 'Tinkerpop') // default database cardinality
    .next();

  return vertex.value;
};
```
-->

### Find Vertices

<!--
TODO: Add Go code for Find Vertices.
javascript code example:
```javascript
/**
 * List all vertexes in db
 * @param {Number} limit
 */
const listAll = async (limit = 500) => {
  return g.V().limit(limit).elementMap().toList();
};
/**
 * Find unique vertex with id
 * @param {Object} vertexId Vertex Id
 */
const findVertex = async (vertexId) => {
  const vertex = await g.V(vertexId).elementMap().next();
  return vertex.value;
};
/**
 * Find vertices by label and 'name' property
 * @param {String} vlabel Vertex label
 * @param {String} name value of 'name' property
 */
const listByLabelAndName = async (vlabel, name) => {
  return g.V().has(vlabel, 'name', name).elementMap().toList();
};
```
-->

### Update Vertex

<!--
TODO: Add Go code for Update Vertex.
javascript code example:
```javascript
const { cardinality: { single } } = gremlin.process;

/**
 * Update Vertex Properties
 * @param {String,Number} vertexId Vertex Id
 * @param {String} name Vertex Name Property
 */
const updateVertex = async (vertexId, label, name) => {
  const vertex = await g.V(vertexId).property(single, 'name', name).next();
  return vertex.value;
};
```
-->

NOTE that versions suffixed with "-rc" are considered release candidates (i.e. pre-alpha, alpha, beta, etc.) and thus 
for early testing purposes only.

## Test Coverage

[![codecov](https://codecov.io/gh/Bit-Quill/gremlin-go/branch/main/graph/badge.svg?token=lzavk3wBTi)](https://codecov.io/gh/Bit-Quill/gremlin-go)