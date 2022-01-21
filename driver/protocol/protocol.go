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

package protocol

import (
	"errors"
	"fmt"
	"gremlin-go/driver/codec"
	"gremlin-go/driver/results"
	"gremlin-go/driver/transport"
	"net/http"
)

type Protocol interface {
	ConnectionMade(transport *string)
	DataReceived(requestId int, requestMessage *string)
	Write(message *string, results map[string]interface{})
}

type AbstractProtocol struct {
	Protocol

	transporter transport.Transporter
}

type GremlinServerWSProtocol struct {
	*AbstractProtocol

	serializer       *codec.Serializer
	maxContentLength int
	username         string
	password         string
}

func (protocol *AbstractProtocol) ConnectionMade(transporter *transport.Transporter) {
	protocol.transporter = *transporter
}

func (protocol *GremlinServerWSProtocol) DataReceived(message *[]byte, resultSets map[string]results.ResultSet) (uint16, error) {
	if message == nil {
		return 0, errors.New("malformed ws or wss URL")
	}
	response, err := (*protocol.serializer).DeserializeMessage(*message)
	if err != nil {
		return 0, err
	}

	requestId, statusCode, metadata, data := response.RequestID, response.ResponseStatus.Code,
		response.ResponseResult.Meta, response.ResponseResult.Data

	resultSet := resultSets[requestId.String()]
	if resultSet == nil {
		resultSet = results.NewChannelResultSet()
	}
	if aggregateTo, ok := metadata["aggregateTo"]; ok {
		resultSet.SetAggregateTo(aggregateTo.(string))
	}
	if statusCode == http.StatusProxyAuthRequired {
		// TODO AN-989: Implement authentication (including handshaking).
		return 0, errors.New("authentication is not currently supported")
	} else if statusCode == http.StatusNoContent {
		// Add empty slice to result.
		resultSet.AddResult(results.NewResult(make([]interface{}, 0)))
		return statusCode, nil
	} else if statusCode == http.StatusOK || statusCode == http.StatusPartialContent {
		// Add data to the ResultSet.
		resultSet.AddResult(results.NewResult(data))
		if statusCode == http.StatusOK {
			resultSet.SetStatusAttributes(response.ResponseStatus.Attributes)
		}
		return statusCode, nil
	} else {
		return 0, errors.New(fmt.Sprint("statusCode: ", statusCode))
	}
}

func (protocol *GremlinServerWSProtocol) Write(requestMessage *codec.Request) error {
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
