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
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

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
	ser *graphBinaryTypeSerializer
}

func newGraphBinarySerializer(handler *logHandler) serializer {
	ser := graphBinaryTypeSerializer{NullType, nil, nil, nil, handler}
	return graphBinarySerializer{&ser}
}

const versionByte byte = 0x81

func convertArgs(request *request, gs graphBinarySerializer) (map[string]interface{}, error) {
	// TODO: Remote transaction session processor is same as bytecode
	if request.processor == bytecodeProcessor {
		// Convert to format:
		// args["gremlin"]: <serialized args["gremlin"]
		gremlin := request.args["gremlin"]
		switch gremlin.(type) {
		case bytecode:
			buffer := bytes.Buffer{}
			gremlinBuffer, err := gs.ser.write(gremlin, &buffer)
			if err != nil {
				return nil, err
			}
			request.args["gremlin"] = gremlinBuffer
			return request.args, nil
		default:
			return nil, errors.New("<<<Error message here>>>")
		}
	} else {
		// Use standard processor, which effectively does nothing.
		return request.args, nil
	}
}

// serializeMessage serializes a request message into GraphBinary
func (gs graphBinarySerializer) serializeMessage(request *request) ([]byte, error) {
	args, err := convertArgs(request, gs)
	if err != nil {
		return nil, err
	}
	finalMessage, err := gs.buildMessage(request.requestID, byte(len(graphBinaryMimeType)), request.op, request.processor, args)
	if err != nil {
		return nil, err
	}
	return finalMessage, nil
}

func writeStr(buffer bytes.Buffer, str string) error {
	err := binary.Write(&buffer, binary.BigEndian, int64(len(str)))
	if err != nil {
		return err
	}
	_, err = buffer.WriteString(str)
	return err
}

func (gs *graphBinarySerializer) buildMessage(id uuid.UUID, mimeLen byte, op string, processor string, args map[string]interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}

	// mime header
	buffer.WriteByte(mimeLen)
	buffer.WriteString(graphBinaryMimeType)

	// Version
	buffer.WriteByte(versionByte)

	// Request uuid
	bigIntUUID := uuidToBigInt(id)
	lower := bigIntUUID.Uint64()
	upperBigInt := bigIntUUID.Rsh(&bigIntUUID, 64)
	upper := upperBigInt.Uint64()
	err := binary.Write(&buffer, binary.BigEndian, upper)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buffer, binary.BigEndian, lower)
	if err != nil {
		return nil, err
	}

	// op
	err = binary.Write(&buffer, binary.BigEndian, uint32(len(op)))
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(op)
	if err != nil {
		return nil, err
	}

	// processor
	err = binary.Write(&buffer, binary.BigEndian, uint32(len(processor)))
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(processor)
	if err != nil {
		return nil, err
	}

	// args
	err = binary.Write(&buffer, binary.BigEndian, uint32(len(args)))
	for k, v := range args {

		// Python:
		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129  54   6 193  42 112 245 79 138 169  36 151 150 34 33 237 190 0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 2 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 3 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]
		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 178 142 53 137 183 117 64 25 165 170 74 181 172 244 93 186 0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108   0 0 0 2 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 2 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 3 0 0 0 0 7 97 108 105 97 115 101 115 ?? ? 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]464 <nil>

		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 178 142 53 137 183 117 64 25 165 170 74 181 172 244 93 186 0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108   0 0 0 2 3 0 0 0 0 7 97 108 105 97 115 101 115 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 2 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0]464 <nil>

		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 178 142 53 137 183 117 64 25 165 170 74 181 172 244 93 186 0 0 0   8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 97 108 105 97 115 101 115 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 2 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0]
		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 231 219  57 118 230  91 67  94 129 216 250  84 57 41 102 153 0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 2 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]
		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 60 126 3 108 39 219 79 62 181 126 242 146 238 210 8 87       0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 2 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 0 0 0 0 7 97 108 105 97 115 101 115 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]

		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 71 225 169 72 3 97 78 34 172 195 198 124 59 219 245 12        0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 3 0 0 0 1 86 0 0 0 0 0 0 0 8 104 97 115 76 97 98 101 108 0 0 0 1 3 0 0 0 0 4 84 101 115 116 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 3 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]
		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 20 203 174 148 119 209 74 130 190 153 100 232 208 148 122 142 0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 3 0 0 0 1 86 0 0 0 0 0 0 0 8 104 97 115 76 97 98 101 108 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 4 84 101 115 116 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 3 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]
		// [32 97 112 112 108 105 99 97 116 105 111 110 47 118 110 100 46 103 114 97 112 104 98 105 110 97 114 121 45 118 49 46 48 129 104 219 81 107 98 88 77 57 150 62 54 87 186 60 145 236        0 0 0 8 98 121 116 101 99 111 100 101 0 0 0 9 116 114 97 118 101 114 115 97 108 0 0 0 2 3 0 0 0 0 7 103 114 101 109 108 105 110 21 0 0 0 0 5 0 0 0 1 86 0 0 0 0 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 1 86 0 0 0 0 0 0 0 8 104 97 115 76 97 98 101 108 0 0 0 1 0 0 0 1 9 0 0 0 0 1 3 0 0 0 0 4 84 101 115 116 0 0 0 5 99 111 117 110 116 0 0 0 0 0 0 0 0 3 0 0 0 0 7 97 108 105 97 115 101 115 10 0 0 0 0 1 3 0 0 0 0 1 103 3 0 0 0 0 1 103]
		_, err = gs.ser.write(k, &buffer)
		if err != nil {
			return nil, err
		}

		switch v.(type) {
		case []byte:
			_, err = buffer.Write(v.([]byte))
			break
		default:
			_, err = gs.ser.write(v, &buffer)
			break
		}
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(fmt.Printf("%v", buffer.Bytes()))
	return buffer.Bytes(), nil
}

/**
bytearray(b' application/vnd.graphbinary-v1.0\x814gZ\xca\xe2\x9eG\x1a\xa9@\x99\x97q\x85\xac\x92\x00\x00\x00\x08bytecode\x00\x00\x00\ttraversal\x00\x00\x00\x02\x03\x00\x00\x00\x00\x07gremlin\x15\x00\x00\x00\x00\x02\x00\x00\x00\x01V\x00\x00\x00\x00\x00\x00\x00\x05count\x00\x00\x00\x00\x00\x00\x00\x00\x03\x00\x00\x00\x00\x07aliases\n\x00\x00\x00\x00\x01\x03\x00\x00\x00\x00\x01g\x03\x00\x00\x00\x00\x01g')
*/
func uuidToBigInt(requestID uuid.UUID) big.Int {
	var bigInt big.Int
	bigInt.SetString(strings.Replace(requestID.String(), "-", "", 4), 16)
	return bigInt
}

func readUUID(buffer *bytes.Buffer) (uuid.UUID, error) {
	var nullable byte
	err := binary.Read(buffer, binary.LittleEndian, &nullable)
	if err != nil {
		return uuid.UUID{}, err
	}
	uuidBytes := make([]byte, 16)
	err = binary.Read(buffer, binary.LittleEndian, uuidBytes)
	return uuid.FromBytes(uuidBytes)
}

func readMap(buffer *bytes.Buffer, gs *graphBinarySerializer) (map[string]interface{}, error) {
	var mapSize uint32
	err := binary.Read(buffer, binary.BigEndian, &mapSize)
	if err != nil {
		return nil, err
	}
	var mapData = map[string]interface{}{}
	for i := uint32(0); i < mapSize; i++ {
		var keyType DataType
		err = binary.Read(buffer, binary.BigEndian, &keyType)
		if err != nil {
			return nil, err
		} else if keyType != StringType {
			return nil, errors.New(fmt.Sprintf("expected string key for map, got type='0x%x'", keyType))
		}
		var nullable byte
		err = binary.Read(buffer, binary.BigEndian, &nullable)
		if nullable != 0 {
			return nil, errors.New("expected non-null key for map")
		}

		k, err := readString(buffer)
		if err != nil {
			return nil, err
		}
		mapData[k], err = gs.ser.read(buffer)
		if err != nil {
			return nil, err
		}
	}
	return mapData, nil
}

func readString(buffer *bytes.Buffer) (string, error) {
	var strLength uint32
	err := binary.Read(buffer, binary.BigEndian, &strLength)
	if err != nil {
		return "", err
	}

	strBytes := make([]byte, strLength)
	err = binary.Read(buffer, binary.BigEndian, strBytes)
	if err != nil {
		return "", err
	}
	return string(strBytes[:]), nil
}

// deserializeMessage deserializes a response message
func (gs graphBinarySerializer) deserializeMessage(responseMessage []byte) (response, error) {
	var msg response
	buffer := bytes.Buffer{}
	buffer.Write(responseMessage)

	// Version
	_, err := buffer.ReadByte()
	if err != nil {
		return msg, err
	}

	// Response uuid
	msgUUID, err := readUUID(&buffer)
	if err != nil {
		return msg, err
	}

	// Status Code
	var statusCode uint32
	err = binary.Read(&buffer, binary.BigEndian, &statusCode)
	if err != nil {
		return msg, err
	}
	statusCode = statusCode & 0xFF

	// Nullable Status message
	var statusMessageNull byte
	var statusMessage string
	err = binary.Read(&buffer, binary.LittleEndian, &statusMessageNull)
	if statusMessageNull == 0 {
		statusMessage, err = readString(&buffer)
		if err != nil {
			return msg, err
		}
	}

	// Status Attributes
	statusAttributes, err := readMap(&buffer, &gs)
	if err != nil {
		return msg, err
	}

	// Meta Attributes
	metaAttributes, err := readMap(&buffer, &gs)
	if err != nil {
		return msg, err
	}

	// Result data
	data, err := gs.ser.read(&buffer)
	if err != nil {
		return msg, err
	}

	msg.responseID = msgUUID
	msg.responseStatus.code = uint16(statusCode)
	msg.responseStatus.message = statusMessage
	msg.responseStatus.attributes = statusAttributes
	msg.responseResult.meta = metaAttributes
	msg.responseResult.data = data

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
	msgUUID, err := gs.ser.readValue(&buffer, byte(UUIDType), false)
	if err != nil {
		return msg, err
	}
	msgOp, err := gs.ser.readValue(&buffer, byte(StringType), false)
	if err != nil {
		return msg, err
	}
	msgProc, err := gs.ser.readValue(&buffer, byte(StringType), false)
	if err != nil {
		return msg, err
	}
	msgArgs, err := gs.ser.readValue(&buffer, byte(MapType), false)
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
	_, err := gs.ser.writeValue(response.responseID, &buffer, true)
	if err != nil {
		return nil, err
	}
	// Status Code
	_, err = gs.ser.writeValue(response.responseStatus.code, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Status message
	_, err = gs.ser.writeValue(response.responseStatus.message, &buffer, true)
	if err != nil {
		return nil, err
	}
	// Status attributes
	_, err = gs.ser.writeValue(response.responseStatus.attributes, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Result meta
	_, err = gs.ser.writeValue(response.responseResult.meta, &buffer, false)
	if err != nil {
		return nil, err
	}
	// Result
	_, err = gs.ser.write(response.responseResult.data, &buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
