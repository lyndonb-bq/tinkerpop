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
)

type Protocol interface {
	connectionMade(transport *string)
	dataReceived(requestId int, requestMessage *string)
	write(message *string, results map[string]interface{})
}

type AbstractProtocol struct {
	Protocol

	transporter Transporter
}

type GremlinServerWSProtocol struct {
	*AbstractProtocol

	serializer       *serializer
	maxContentLength int
	username         string
	password         string
}

func (protocol *AbstractProtocol) ConnectionMade(transporter *Transporter) {
	protocol.transporter = *transporter
}

func (protocol *GremlinServerWSProtocol) DataReceived(message *[]byte, resultSets map[string]ResultSet) (uint16, error) {
	if message == nil {
		return 0, errors.New("malformed ws or wss URL")
	}
	response, err := (*protocol.serializer).DeserializeMessage(*message)
	if err != nil {
		return 0, err
	}

	requestId, statusCode, metadata, data := response.requestID, response.responseStatus.code,
		response.responseResult.meta, response.responseResult.data

	resultSet := resultSets[requestId.String()]
	if resultSet == nil {
		resultSet = NewChannelResultSet()
	}
	if aggregateTo, ok := metadata["aggregateTo"]; ok {
		resultSet.SetAggregateTo(aggregateTo.(string))
	}
	if statusCode == http.StatusProxyAuthRequired {
		// TODO AN-989: Implement authentication (including handshaking).
		return 0, errors.New("authentication is not currently supported")
	} else if statusCode == http.StatusNoContent {
		// Add empty slice to result.
		resultSet.AddResult(NewResult(make([]interface{}, 0)))
		return statusCode, nil
	} else if statusCode == http.StatusOK || statusCode == http.StatusPartialContent {
		// Add data to the ResultSet.
		resultSet.AddResult(NewResult(data))
		if statusCode == http.StatusOK {
			resultSet.SetStatusAttributes(response.responseStatus.attributes)
		}
		return statusCode, nil
	} else {
		return 0, errors.New(fmt.Sprint("statusCode: ", statusCode))
	}
}

func (protocol *GremlinServerWSProtocol) Write(requestMessage *Request) error {
	message, err := (*protocol.serializer).SerializeMessage(requestMessage)
	if err == nil {
		err = protocol.transporter.Write(message)
	}
	return err
}

func NewGremlinServerWSProtocol() *GremlinServerWSProtocol {
	ap := &AbstractProtocol{}
	protocol := &GremlinServerWSProtocol{ap, nil, 1, "", ""}
	return protocol
}
