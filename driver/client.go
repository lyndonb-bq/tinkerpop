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

package driver

import (
	"gremlin-go/driver/transport"
)

// Client is used to connect and interact with a Gremlin-supported server.
type Client struct {
	host            string
	port            int
	transporterType transport.TransporterType
}

// NewClient creates a Client and configures it with the given parameters.
func NewClient(host string, port int, transporterType transport.TransporterType) *Client {
	client := &Client{host, port, transporterType}
	return client
}

// Submit a Gremlin traversal string to execute.
func (client *Client) Submit(traversalString string) (string, error) {
	// TODO AN-982: Obtain connection from pool of connections held by the client.
	connection := &connection{client.host, client.port, client.transporterType}
	return connection.submit(traversalString)
}
