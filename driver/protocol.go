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
	"net/http"

	"github.com/google/uuid"
)

type protocol interface {
	connectionMade(transport transporter)
	dataReceived(message []byte, resultSets map[string]ResultSet) (protocolStatus, error)
	write(message string, results map[string]ResultSet) error
}

type protocolBase struct {
	protocol

	transporter transporter
}

type gremlinServerWSProtocol struct {
	*protocolBase

	serializer       serializer
	logHandler       *logHandler
	maxContentLength int
	username         string
	password         string
}

func (protocol *protocolBase) connectionMade(transporter transporter) {
	protocol.transporter = transporter
}

type protocolStatus = uint16

func (protocol *gremlinServerWSProtocol) dataReceived(message []byte, resultSets map[string]ResultSet) (protocolStatus, error) {
	if message == nil {
		protocol.logHandler.log(Error, malformedURL)
		return 0, errors.New("malformed ws or wss URL")
	}
	response, err := protocol.serializer.deserializeMessage(message)
	if err != nil {
		return 0, err
	}

	requestID, statusCode, metadata, data := response.requestID, response.responseStatus.code,
		response.responseResult.meta, response.responseResult.data

	resultSet := resultSets[requestID.String()]
	if resultSet == nil {
		resultSet = newChannelResultSet(uuid.New().String())
	}
	if aggregateTo, ok := metadata["aggregateTo"]; ok {
		resultSet.setAggregateTo(aggregateTo.(string))
	}
	if statusCode == http.StatusProxyAuthRequired {
		// TODO AN-989: Implement authentication (including handshaking).
		return 0, errors.New("authentication is not currently supported")
	} else if statusCode == http.StatusNoContent {
		// Add empty slice to result.
		resultSet.addResult(newResult(make([]interface{}, 0)))
		return statusCode, nil
	} else if statusCode == http.StatusOK || statusCode == http.StatusPartialContent {
		// Add data to the ResultSet.
		resultSet.addResult(newResult(data))
		if statusCode == http.StatusOK {
			resultSet.setStatusAttributes(response.responseStatus.attributes)
		}
		return statusCode, nil
	} else {
		return 0, errors.New(fmt.Sprint("statusCode: ", statusCode))
	}
}

func (protocol *gremlinServerWSProtocol) write(message string, results map[string]ResultSet) error {
	request := makeStringRequest(message)
	bytes, err := protocol.serializer.serializeMessage(&request)
	if err == nil {
		err = protocol.transporter.Write(bytes)
	}
	return err
}

func newGremlinServerWSProtocol(handler *logHandler) *gremlinServerWSProtocol {
	ap := &protocolBase{}

	protocol := &gremlinServerWSProtocol{ap, newGraphBinarySerializer(handler), handler, 1, "", ""}
	return protocol
}
