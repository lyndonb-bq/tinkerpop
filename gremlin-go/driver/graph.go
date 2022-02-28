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
	"fmt"
	"strings"
)

type Graph struct {
}

// Element is the base structure for both Vertex and Edge.
// The inherited identifier must be unique to the inheriting classes.
type Element struct {
	Id    interface{}
	Label string
}

// Vertex contains a single Vertex which has a Label and an Id.
type Vertex struct {
	Element
}

// Edge links two Vertex structs along with its Property Objects. An edge has both a direction and a Label.
type Edge struct {
	Element
	OutV Vertex
	InV  Vertex
}

// VertexProperty is similar to propery in that it denotes a Key/Value pair associated with a Vertex, but is different
// in that it also represents an entity that is an Element and can have properties of its own.
type VertexProperty struct {
	Element
	Key    string // This is the Label of Vertex.
	Value  interface{}
	Vertex Vertex // Vertex that owns the property.
}

// Property denotes a Key/Value pair associated with an Edge. A property can be empty.
type Property struct {
	Key     string
	Value   interface{}
	Element Element
}

// Path denotes a particular walk through a Graph as defined by a traversal.
// A list of Labels and a list of Objects is maintained in the path.
// The list of Labels are the Labels of the steps traversed, and the Objects are the Objects that are traversed.
// TODO: change Labels to be []<set of string> after implementing set in AN-1022 and update the GetPathObject accordingly
type Path struct {
	Labels  [][]string
	Objects []interface{}
}

func (v *Vertex) String() string {
	return fmt.Sprintf("v[%s]", v.Id)
}

func (e *Edge) String() string {
	return fmt.Sprintf("e[%s][%s-%s->%s]", e.Id, e.OutV.Id, e.Label, e.InV.Id)
}

func (vp *VertexProperty) String() string {
	return fmt.Sprintf("vp[%s->%v]", vp.Label, vp.Value)
}

func (p *Property) String() string {
	return fmt.Sprintf("p[%s->%v]", p.Key, p.Value)
}

func (p *Path) String() string {
	return fmt.Sprintf("path[%s]", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(p.Objects)), ", "), "[]"))
}

// GetPathObject returns the Value that corresponds to the Key for the Path and error if the Value is not present or cannot be retrieved.
func (p *Path) GetPathObject(key string) (interface{}, error) {
	if len(p.Objects) != len(p.Labels) {
		return nil, errors.New("path is invalid because it does not contain an equal number of Labels and Objects")
	}
	var objectList []interface{}
	var object interface{}
	for i := 0; i < len(p.Labels); i++ {
		for j := 0; j < len(p.Labels[i]); j++ {
			if p.Labels[i][j] == key {
				if object == nil {
					object = p.Objects[i]
				} else if objectList != nil {
					objectList = append(objectList, p.Objects[i])
				} else {
					objectList = []interface{}{object, p.Objects[i]}
				}
			}
		}
	}
	if objectList != nil {
		return objectList, nil
	} else if object != nil {
		return object, nil
	} else {
		return nil, fmt.Errorf("path does not contain a Label of '%s'", key)
	}
}
