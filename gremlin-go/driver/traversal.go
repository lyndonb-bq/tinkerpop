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

type Traverser struct {
	bulk  int64
	value interface{}
}

type Traversal struct {
	graph               *Graph
	traversalStrategies *TraversalStrategies
	bytecode            *bytecode
	traverser           *Traverser
	remote              *DriverRemoteConnection
}

// ToList returns the result in a list
// TODO use TraversalStrategies instead of direct remote after they are implemented
func (t *Traversal) ToList() ([]*Result, error) {
	results, err := t.remote.SubmitBytecode(t.bytecode)
	if err != nil {
		return nil, err
	}
	return results.All(), nil
}

// ToSet returns the results in a set.
// TODO Go doesn't have sets, determine the best structure for this
func (t *Traversal) ToSet() (map[*Result]bool, error) {
	set := map[*Result]bool{}
	results, err := t.remote.SubmitBytecode(t.bytecode)
	if err != nil {
		return nil, err
	}

	for _, r := range results.All() {
		set[r] = true
	}
	return set, nil
}

// Iterate all the Traverser instances in the traversal and returns the empty traversal
func (t *Traversal) Iterate() (*Traversal, error) {
	err := t.bytecode.addStep("none")
	if err != nil {
		return nil, err
	}

	_, err = t.remote.SubmitBytecode(t.bytecode)
	if err != nil {
		return nil, err
	}
	return t, nil
}
