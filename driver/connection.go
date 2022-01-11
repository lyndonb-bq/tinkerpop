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
)

// TODO: make sure these are constants
const scheme = "ws"
const path = "gremlin"

type connection struct {
	host            string
	port            int
	transporterType TransporterType
}

// TODO: refactor this when implementing full connection
func (connection *connection) submit(traversalString string) (response string, err error) {
	transporter := GetTransportLayer(connection.transporterType, connection.host, connection.port)
	defer transporter.Close()

	err = transporter.Write(traversalString)
	if err != nil {
		fmt.Println("Writing request failed!")
		return
	}

	bytes, err := transporter.Read()
	if err != nil {
		fmt.Println("Reading message failed!")
		return
	}

	response = string(bytes)
	transporter.Close()
	return
}
