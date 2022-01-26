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
	"golang.org/x/text/language"

	"github.com/google/uuid"
)

const graphBinaryMimeType = "application/vnd.graphbinary-v1.0"

// serializer interface for serializers
type serializer interface {
	serializeMessage(request *request) ([]byte, error)
	deserializeMessage(message []byte) (response, error)
}

// graphBinarySerializer serializes/deserializes message to/from GraphBinary
type graphBinarySerializer struct {
	readerClass *graphBinaryReader
	writerClass *graphBinaryWriter
	mimeType    string `default:"application/vnd.graphbinary-v1.0"`
}

func newGraphBinarySerializer(handler *logHandler) serializer {
	reader := graphBinaryReader{handler}
	writer := graphBinaryWriter{handler}
	return graphBinarySerializer{&reader, &writer, graphBinaryMimeType}
}

const versionByte byte = 0x81

// serializeMessage serializes a request message into GraphBinary
func (gs graphBinarySerializer) serializeMessage(request *request) ([]byte, error) {
	gs.mimeType = graphBinaryMimeType
	finalMessage, err := gs.buildMessage(request, 0x20, gs.mimeType)
	if err != nil {
		return nil, err
	}
	return finalMessage, nil
}

func (gs *graphBinarySerializer) buildMessage(request *request, mimeLen byte, mimeType string) ([]byte, error) {
	buffer := bytes.Buffer{}

	// mime header
	buffer.WriteByte(mimeLen)
	buffer.WriteString(mimeType)
	// version
	buffer.WriteByte(versionByte)
	// requestID
	logHandler := newLogHandler(&defaultLogger{}, Info, language.English)
	logHandler.logger.Logf(Error, "requestID")
	_, err := gs.writerClass.writeValue(request.requestID, &buffer, false)
	if err != nil {
		return nil, err
	}
	// op
	logHandler.logger.Logf(Error, "op")
	_, err = gs.writerClass.writeValue(request.op, &buffer, false)
	if err != nil {
		return nil, err
	}
	// processor
	logHandler.logger.Logf(Error, "processor")
	_, err = gs.writerClass.writeValue(request.processor, &buffer, false)
	if err != nil {
		return nil, err
	}
	// args
	logHandler.logger.Logf(Error, "args")
	_, err = gs.writerClass.writeValue(request.args, &buffer, false)
	if err != nil {
		return nil, err
	}

	logHandler.logger.Logf(Error, "Done")
	return buffer.Bytes(), nil
}

// deserializeMessage deserializes a response message
func (gs graphBinarySerializer) deserializeMessage(responseMessage []byte) (response, error) {
	var msg response
	buffer := bytes.Buffer{}
	buffer.Write(responseMessage)
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
func (gs *graphBinarySerializer) deserializeRequestMessage(requestMessage *[]byte) (request, error) {
	buffer := bytes.Buffer{}
	var msg request
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

	msg.requestID = msgUUID.(uuid.UUID)
	msg.op = msgOp.(string)
	msg.processor = msgProc.(string)
	msg.args = msgArgs.(map[string]interface{})

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
