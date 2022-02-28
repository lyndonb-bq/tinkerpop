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
	"github.com/lyndonb-bq/tinkerpop/gremlin-go/driver"
)

//****************************************************************************************
//* THIS IS THE TEST FILE - will be replaced by file generated with build/generate.groovy
//****************************************************************************************

var translationMap = map[string][]func(g *gremlingo.GraphTraversalSource, p map[string]interface{}) *gremlingo.GraphTraversal{
	"g_injectXnull_nullX": {func(g *gremlingo.GraphTraversalSource, p map[string]interface{}) *gremlingo.GraphTraversal {
		return g.Inject(nil, nil)
	}},
	"g_VX1X_valuesXageX_injectXnull_nullX": {func(g *gremlingo.GraphTraversalSource, p map[string]interface{}) *gremlingo.GraphTraversal {
		return g.V(p["xx1"]).Values("age").Inject(nil, nil)
	}},
}

func GetTraversal(scenarioName string, g *gremlingo.GraphTraversalSource, parameters map[string]interface{}) (*gremlingo.GraphTraversal, error) {
	if traversalFns, ok := translationMap[scenarioName]; ok {
		traversal := traversalFns[0]
		traversalFns = traversalFns[1:]
		return traversal(g, parameters), nil
	} else {
		return nil, errors.New("scenario for traversal not recognized")
	}
}
