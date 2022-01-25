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

type connection struct {
	host            string
	port            int
	transporterType TransporterType
	logHandler      *logHandler
	transporter     transporter
	protocol        protocol
	results         map[string]ResultSet
}

func (connection *connection) close() (err error) {
	if connection.transporter != nil {
		err = connection.transporter.Close()
	}
	return
}

func (connection *connection) connect() error {
	if connection.transporter != nil {
		closeErr := connection.transporter.Close()
		connection.logHandler.logf(Warning, transportCloseFailed, closeErr)
	}

	connection.transporter = getTransportLayer(connection.transporterType, connection.host, connection.port)
	err := connection.transporter.Connect()
	if err == nil {
		connection.protocol.connectionMade(connection.transporter)
	}
	return err
}

func (connection *connection) write(traversalString string) ([]byte, error) {
	if connection.transporter == nil || connection.transporter.IsClosed() {
		err := connection.connect()
		if err != nil {
			return nil, err
		}
	}

	requestID := uuid.New().String()
	if connection.results == nil {
		connection.results = map[string]ResultSet{}
	}
	connection.results[requestID] = newChannelResultSet(requestID)

	err := connection.protocol.write(traversalString, connection.results)
	if err != nil {
		return nil, err
	}

	data, err := connection.transporter.Read()
	if err != nil {
		return nil, err
	}

	// TODO: Is checking the status required? GorillaTransporter.Read() abstracts it away
	_, err = connection.protocol.dataReceived(data, connection.results)
	if err != nil {
		return nil, err
	}
	return data, nil
}
