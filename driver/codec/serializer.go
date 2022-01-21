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

package codec

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
)

type graphBinMessageSerializer struct {
	readerClass graphBinaryReader
	writerClass graphBinaryWriter
	mimeType    string `default:"application/vnd.graphbinary-v1.0"`
}

const versionByte byte = 0x81

func (gs *graphBinMessageSerializer) getProcessor(processor string) string {
	return processor
}

func (gs *graphBinMessageSerializer) serializeMessage(request *request) []byte {
	finalMessage, _ := gs.buildMessage(request, 0x20, gs.mimeType)
	return finalMessage
}

func (gs *graphBinMessageSerializer) buildMessage(message *request, mimeLen byte, mimeType string) ([]byte, error) {
	byteBuffer := bytes.Buffer{}

	// header
	byteBuffer.WriteByte(mimeLen)
	byteBuffer.WriteString(mimeType)
	byteBuffer.WriteByte(versionByte)

	// RequestID
	binary.Write(&byteBuffer, binary.BigEndian, message.RequestID)

	// Op
	gs.writerClass.writeValue(message.Op, &byteBuffer, false)
	//binary.Write(&byteBuffer, binary.BigEndian, len(message.Op))
	//binary.Write(&byteBuffer, binary.BigEndian, message.Op)

	// Processor
	gs.writerClass.writeValue(message.Processor, &byteBuffer, false)
	//binary.Write(&byteBuffer, binary.BigEndian, len(message.Processor))
	//binary.Write(&byteBuffer, binary.BigEndian, message.Processor)

	// Args
	args := message.Args
	binary.Write(&byteBuffer, binary.BigEndian, int32(len(args)))
	//binary.Write(&byteBuffer, binary.BigEndian, args)
	for k, v := range args {
		gs.writerClass.writeObject(k, &byteBuffer)
		gs.writerClass.writeObject(v, &byteBuffer)
	}

	return byteBuffer.Bytes(), nil
}

func (gs *graphBinMessageSerializer) deserializeMessage(message []byte) response {
	buff := bytes.Buffer{}
	buff.Write(message)
	msgUUID, _ := gs.readerClass.readValue(&buff, byte(uuidType), true)
	msgCode, _ := gs.readerClass.readValue(&buff, byte(intType), false)
	msgMsg, _ := gs.readerClass.readValue(&buff, byte(stringType), true)
	msgAttr, _ := gs.readerClass.readValue(&buff, byte(mapType), false)
	msgMeta, _ := gs.readerClass.readValue(&buff, byte(mapType), true)
	msgData := gs.readerClass.read(&buff)
	var msg response
	msg.RequestID = msgUUID.(uuid.UUID)
	msg.ResponseStatus.Code = msgCode.(uint32)
	msg.ResponseStatus.Message = msgMsg.(string)
	msg.ResponseStatus.Attributes = msgAttr.(map[interface{}]interface{})
	msg.ResponseResult.Meta = msgMeta.(map[interface{}]interface{})
	msg.ResponseResult.Data = msgData

	//msg := map[string]interface{}{
	//	"requestId": request_id,
	//	"status": map[string]interface{}{"code": status_code,
	//		"message":    status_msg,
	//		"attributes": status_attrs},
	//	"result": map[string]interface{}{"meta": meta_attrs,
	//		"data": result}}
	return msg
}
