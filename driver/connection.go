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
	"github.com/gorilla/websocket"
	"net/url"
	"strconv"
)

// TODO: make sure these are constants
const scheme = "ws"
const path = "gremlin"

type connection struct {
	host string
	port int
}

// TODO: refactor this when implementing full connection
func (connection *connection) submit(traversalString string) (response string, err error) {
	u := url.URL{
		Scheme: scheme,
		Host:   connection.host + ":" + strconv.Itoa(connection.port),
		Path:   path,
	}

	dialer := websocket.DefaultDialer
	// TODO: make this configurable from client; this currently does nothing since 4096 is the default
	dialer.WriteBufferSize = 4096
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Connecting failed!")
		return
	}
	defer conn.Close()

	err = conn.WriteJSON(makeStringRequest(traversalString))
	if err != nil {
		fmt.Println("Writing request failed!")
		return
	}

	_, responseMessage, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Reading message failed!")
		return
	}

	response = string(responseMessage)
	return
}
