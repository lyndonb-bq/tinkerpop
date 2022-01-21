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

	codec            codec.Codec
	maxContentLength int
	username         string
	password         string
}

func (protocol *AbstractProtocol) ConnectionMade(transporter *transport.Transporter) {
	protocol.transporter = *transporter
}

func (protocol *GremlinServerWSProtocol) DataReceived(message *[]byte, resultSets map[string]results.ResultSet) (uint32, error) {
	if message == nil {
		return 0, errors.New("malformed ws or wss URL")
	}

	// TOOD: Update so pointer on other side of DeserializeMessage
	response := (*protocol.codec.Deserializer).DeserializerMessage(*message)
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
		// TODO: AN-968: Insert empty list into resultSet.
		return statusCode, nil
	} else if statusCode == http.StatusOK || statusCode == http.StatusPartialContent {
		resultSet.AddResult(results.NewResult(data))
		if statusCode == http.StatusOK {
			resultSet.SetStatusAttributes(response.ResponseStatus.Attributes)
		}
		return statusCode, nil
	} else {
		return 0, errors.New(fmt.Sprint("statusCode: ", statusCode))
	}
}

func CheckAndSet(key string, funcToSet func(interface{}), dictionary map[interface{}]interface{}) {
	if value, ok := dictionary[key]; ok {
		funcToSet(value.(string))
	}
}

func (protocol *GremlinServerWSProtocol) ParseResult(message map[string]interface{}, resultSets map[string]results.ResultSet) (int, error) {
	if message == nil {
		return 0, errors.New("malformed ws or wss URL")
	}

	// Deserialize

	msg := map[string]interface{}{
		"requestId": "request_id",
		"status": map[interface{}]interface{}{
			"code":       "status_code",
			"message":    "status_msg",
			"attributes": "status_attrs",
		},
		"result": map[interface{}]interface{}{
			"meta": "",
			"data": "result",
		},
	}

	requestId, statusCode, aggregateTo, data := unpackMessage(msg)
	resultSet := resultSets[requestId]
	if resultSet == nil {
		resultSet = results.NewChannelResultSet()
		// TODO: Add resultset to resultSets
	}
	resultSet.SetAggregateTo(aggregateTo)
	if statusCode == http.StatusProxyAuthRequired {
		// TODO AN-989: Implement authentication (including handshaking).
		return 0, errors.New("authentication is not currently supported")
	} else if statusCode == http.StatusNoContent {
		// TODO: AN-968: Insert empty list into resultSet.
		return statusCode, nil
	} else if statusCode == http.StatusOK || statusCode == http.StatusPartialContent {
		// TODO AN-968: Insert data into resultSet.
		if statusCode == http.StatusOK {
			resultSet.SetStatusAttributes(msg["status"].(map[string]interface{})["attributes"])
			removeResultSetsKey(resultSets, requestId)
		}
		return statusCode, nil
	} else {
		return 0, errors.New(fmt.Sprint("statusCode: ", statusCode, ", message: ", msg))
	}
}

func removeResultSetsKey(resultSets map[string]results.ResultSet, key string) {
	// Need to check if exists before removal, otherwise we might cause a panic.
	if _, ok := resultSets[key]; ok {
		delete(resultSets, key)
	}
}

func unpackMessage(message map[string]interface{}) (string, int, string, interface{}) {
	requestId := message["requestId"].(string)
	statusCode := message["status"].(map[string]interface{})["code"].(int)
	aggregateTo := message["result"].(map[string]interface{})["meta"].(map[string]interface{})["aggregateTo"]
	if aggregateTo == nil {
		aggregateTo = "list"
	}
	data := message["result"].(map[string]interface{})["data"]
	return requestId, statusCode, aggregateTo.(string), data
}

func (protocol *GremlinServerWSProtocol) Write(requestId int, requestMessage *string) {
	// TODO: Fix.
	// var message string = protocol.serializer.serialize_message(requestId, requestMessage)
	// protocol.transport.write(message)
}

func NewGremlinServerWSProtocol() *GremlinServerWSProtocol {
	ap := &AbstractProtocol{}
	protocol := &GremlinServerWSProtocol{ap, nil, 1, "", ""}
	return protocol
}
