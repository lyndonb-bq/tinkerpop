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

import "github.com/google/uuid"

// TODO: remove these constants
const op = "eval"
const processor = ""
const graphType = "g:Map"

type Request struct {
	RequestID uuid.UUID                   `json:"requestId"`
	Op        string                      `json:"op"`
	Processor string                      `json:"processor"`
	Args      map[interface{}]interface{} `json:"args"`
}

func makeStringRequest(requestString string) (req Request) {
	req.RequestID = uuid.New()
	req.Op = op
	req.Processor = processor
	req.Args = make(map[interface{}]interface{})
	req.Args["@type"] = graphType
	value := make([]string, 2)
	value[0] = "gremlin"
	value[1] = requestString
	req.Args["@value"] = value
	return
}
