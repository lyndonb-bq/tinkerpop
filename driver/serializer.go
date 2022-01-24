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
	"bytes"

	"github.com/google/uuid"
)

// serializer interface for serializers
type serializer interface {
	serializeMessage(request *Request) ([]byte, error)
	deserializeMessage(message *[]byte) (response, error)
}

// graphBinarySerializer serializes/deserializes message to/from GraphBinary
type graphBinarySerializer struct {
	readerClass *GraphBinaryReader
	writerClass *GraphBinaryWriter
	mimeType    string `default:"application/vnd.graphbinary-v1.0"`
}

func newGraphBinarySerializer() serializer {
	return graphBinarySerializer{}
}

const versionByte byte = 0x81

// serializeMessage serializes a request message into GraphBinary
func (gs graphBinarySerializer) serializeMessage(request *Request) ([]byte, error) {
	gs.mimeType = "application/vnd.graphbinary-v1.0"
	finalMessage, err := gs.buildMessage(request, 0x20, gs.mimeType)
	if err != nil {
		return nil, err
	}
	return finalMessage, nil
}

func (gs *graphBinarySerializer) buildMessage(request *Request, mimeLen byte, mimeType string) ([]byte, error) {
	buffer := bytes.Buffer{}

	// mime header
	buffer.WriteByte(mimeLen)
	buffer.WriteString(mimeType)
	// version
	buffer.WriteByte(versionByte)
	// requestID
	_, err := gs.writerClass.writeValue(request.RequestID, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Op
	_, err = gs.writerClass.writeValue(request.Op, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Processor
	_, err = gs.writerClass.writeValue(request.Processor, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Args
	_, err = gs.writerClass.writeValue(request.Args, &buffer, false)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// deserializeMessage deserializes a response message
func (gs graphBinarySerializer) deserializeMessage(responseMessage *[]byte) (response, error) {
	var msg response
	buffer := bytes.Buffer{}
	buffer.Write(*responseMessage)
	// version
	_, err := buffer.ReadByte()
	if err != nil {
		return msg, err
	}
	// UUID
	msgUUID, err := gs.readerClass.readValue(&buffer, byte(UUIDType), true)
	if err != nil {
		return msg, err
	}
	// Status Code
	msgCode, err := gs.readerClass.readValue(&buffer, byte(IntType), false)
	if err != nil {
		return msg, err
	}
	// Nullable Status message
	msgMsg, err := gs.readerClass.readValue(&buffer, byte(StringType), true)
	if err != nil {
		return msg, err
	}
	// Status Attribute
	msgAttr, err := gs.readerClass.readValue(&buffer, byte(MapType), false)
	if err != nil {
		return msg, err
	}
	// Result meta
	msgMeta, err := gs.readerClass.readValue(&buffer, byte(MapType), false)
	if err != nil {
		return msg, err
	}
	// Result data
	msgData, err := gs.readerClass.read(&buffer)
	if err != nil {
		return msg, err
	}

	msg.requestID = msgUUID.(uuid.UUID)
	msg.responseStatus.code = uint16(msgCode.(int32))
	msg.responseStatus.message = msgMsg.(string)
	msg.responseStatus.attributes = msgAttr.(map[interface{}]interface{})
	msg.responseResult.meta = msgMeta.(map[interface{}]interface{})
	msg.responseResult.data = msgData

	return msg, nil
}

// private function for deserializing a request message for testing purposes
func (gs *graphBinarySerializer) deserializeRequestMessage(requestMessage *[]byte) (Request, error) {
	buffer := bytes.Buffer{}
	var msg Request
	buffer.Write(*requestMessage)
	// skip headers
	buffer.Next(33)
	// version
	_, err := buffer.ReadByte()
	if err != nil {
		return msg, err
	}
	msgUUID, err := gs.readerClass.readValue(&buffer, byte(UUIDType), false)
	if err != nil {
		return msg, err
	}
	msgOp, err := gs.readerClass.readValue(&buffer, byte(StringType), false)
	if err != nil {
		return msg, err
	}
	msgProc, err := gs.readerClass.readValue(&buffer, byte(StringType), false)
	if err != nil {
		return msg, err
	}
	msgArgs, err := gs.readerClass.readValue(&buffer, byte(MapType), false)
	if err != nil {
		return msg, err
	}

	msg.RequestID = msgUUID.(uuid.UUID)
	msg.Op = msgOp.(string)
	msg.Processor = msgProc.(string)
	msg.Args = msgArgs.(map[interface{}]interface{})

	return msg, nil
}

// private function for serializing a response message for testing purposes
func (gs *graphBinarySerializer) serializeResponseMessage(response *response) ([]byte, error) {
	buffer := bytes.Buffer{}

	// version
	buffer.WriteByte(versionByte)

	// requestID
	_, err := gs.writerClass.writeValue(response.requestID, &buffer, true)
	if err != nil {
		return nil, err
	}
	// Status Code
	_, err = gs.writerClass.writeValue(response.responseStatus.code, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Status message
	_, err = gs.writerClass.writeValue(response.responseStatus.message, &buffer, true)
	if err != nil {
		return nil, err
	}
	// Status attributes
	_, err = gs.writerClass.writeValue(response.responseStatus.attributes, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Result meta
	_, err = gs.writerClass.writeValue(response.responseResult.meta, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Result
	_, err = gs.writerClass.write(response.responseResult.data, &buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
