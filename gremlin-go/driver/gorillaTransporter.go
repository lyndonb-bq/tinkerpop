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
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

type gorillaTransporter struct {
	host       string
	port       int
	connection websocketConn
	isClosed   bool
}

func (transporter *gorillaTransporter) Connect() (err error) {
	if transporter.connection != nil {
		return
	}

	u := url.URL{
		Scheme: scheme,
		Host:   transporter.host + ":" + strconv.Itoa(transporter.port),
		Path:   path,
	}

	dialer := websocket.DefaultDialer
	// TODO: make this configurable from client; this currently does nothing since 4096 is the default
	dialer.WriteBufferSize = 4096
	conn, _, err := dialer.Dial(u.String(), nil)
	if err == nil {
		transporter.connection = conn
	}
	return
}

func (transporter *gorillaTransporter) Write(data []byte) (err error) {
	if transporter.connection == nil {
		err = transporter.Connect()
		if err != nil {
			return
		}
	}

	err = transporter.connection.WriteMessage(websocket.BinaryMessage, data)
	return err
}

func (transporter *gorillaTransporter) Read() (bytes []byte, err error) {
	if transporter.connection == nil {
		err = transporter.Connect()
		if err != nil {
			return
		}
	}

	_, bytes, err = transporter.connection.ReadMessage()
	return
}

func (transporter *gorillaTransporter) Close() (err error) {
	if transporter.connection != nil && !transporter.isClosed {
		transporter.isClosed = true
		return transporter.connection.Close()
	}
	return
}

func (transporter *gorillaTransporter) IsClosed() bool {
	return transporter.isClosed
}