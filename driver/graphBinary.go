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
	"github.com/google/uuid"
	"reflect"
)

// Version 1.0

type dataType int

const (
	nullType   dataType = 0xFE
	intType    dataType = 0x01
	longType   dataType = 0x02
	doubleType dataType = 0x07
	floatType  dataType = 0x08

	stringType dataType = 0x03
	uuidType   dataType = 0x0c
	mapType    dataType = 0x0a
)

var nullBytes = []byte{nullType.getCodeByte(), 0x01}

func (dataType dataType) getCode() int {
	return int(dataType)
}

func (dataType dataType) getCodeByte() byte {
	return byte(dataType)
}

func (dataType dataType) getCodeBytes() []byte {
	return []byte{dataType.getCodeByte()}
}

type graphBinaryWriter struct {
}

const (
	valueFlagNull byte = 1
	valueFlagNone byte = 0
)

type graphBinarySerializer interface {
	// change type of value for each specific serializer?
	write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error)
	writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error)
	read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error)
	readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error)
	getDataType() dataType
}

// gets the type of the serializer based on the value
func (writer *graphBinaryWriter) getSerializer(val interface{}) (graphBinarySerializer, error) {
	switch val.(type) {
	case string:
		return &stringSerializer{}, nil
	case int64, int, uint32:
		return &longSerializer{}, nil
	case int32, int8, uint16:
		return &intSerializer{}, nil
	case uuid.UUID:
		return &uuidSerializer{}, nil
	default:
		switch reflect.TypeOf(val).Kind() {
		case reflect.Map:
			return &mapSerializer{}, nil
		default:
			return nil, errors.New("unknown data type")
		}
	}
}

// gets the type of the serializer based on the value
func (reader *graphBinaryReader) getSerializer(val byte) (graphBinarySerializer, error) {
	switch val {
	case stringType.getCodeByte():
		return &stringSerializer{}, nil
	case longType.getCodeByte():
		return &longSerializer{}, nil
	case intType.getCodeByte():
		return &intSerializer{}, nil
	case uuidType.getCodeByte():
		return &uuidSerializer{}, nil
	case mapType.getCodeByte():
		return &mapSerializer{}, nil
	default:
		return nil, errors.New("unknown data type")
	}
}

// should return type be void?
func (writer *graphBinaryWriter) write(objectData interface{}) interface{} {
	buffer := bytes.Buffer{}
	fmt.Println("At writer")
	fmt.Println(objectData)
	return writer.writeObject(objectData, &buffer)
}

// Writes an object in fully-qualified format, containing {type_code}{type_info}{value_flag}{value}.
func (writer *graphBinaryWriter) writeObject(valueObject interface{}, buffer *bytes.Buffer) interface{} {
	if valueObject == nil {
		// return Object of type "unspecified object null" with the value flag set to null.
		buffer.Write(nullBytes)
		return buffer.Bytes()
	}

	serializer, _ := writer.getSerializer(valueObject)
	buffer.Write(serializer.getDataType().getCodeBytes())
	message, _ := serializer.write(valueObject, buffer, writer)
	return message
}

// Writes a value without including type information.
func (writer *graphBinaryWriter) writeValue(value interface{}, buffer *bytes.Buffer, nullable bool) (interface{}, error) {
	if value == nil {
		if !nullable {
			return nil, errors.New("unexpected null value when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	serializer, err := writer.getSerializer(value)
	if err != nil {
		return nil, err
	}
	buffer.Write(serializer.getDataType().getCodeBytes())
	message, err := serializer.writeValue(value, buffer, writer, nullable)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (writer *graphBinaryWriter) writeValueFlagNull(buffer *bytes.Buffer) {
	buffer.WriteByte(valueFlagNull)
}

func (writer *graphBinaryWriter) writeValueFlagNone(buffer *bytes.Buffer) {
	buffer.WriteByte(valueFlagNone)
}

type graphBinaryReader struct {
}

// Reads the type code, information and value of a given buffer with fully-qualified format.
func (reader *graphBinaryReader) read(buffer *bytes.Buffer) interface{} {
	fmt.Println("At reader")
	typeCode, _ := buffer.ReadByte()
	if typeCode == nullType.getCodeByte() {
		check, _ := buffer.ReadByte()
		// check this
		if check != 1 {
			return errors.New("read value flag")
		}
		return nil
	}

	fmt.Println(typeCode)
	serializer, err := reader.getSerializer(typeCode)
	if err != nil {
		panic("cannot get serializer for type")
	}
	val, _ := serializer.read(buffer, reader)
	return val
}

func (reader *graphBinaryReader) readValue(buffer *bytes.Buffer, typ byte, nullable bool) (interface{}, error) {
	if buffer == nil {
		panic("input cannot be null")
	}
	serializer, _ := reader.getSerializer(typ)
	val, _ := serializer.readValue(buffer, reader, nullable)
	return val, nil
}

// Format: 4-byte two’s complement integer.
type intSerializer struct {
	graphBinarySerializer
}

func (intSerializer *intSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return intSerializer.writeValue(value, buffer, writer, true)
}

func (intSerializer *intSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			return nil, errors.New("unexpected null value when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	err := binary.Write(buffer, binary.BigEndian, value)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (intSerializer *intSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return intSerializer.readValue(buffer, reader, true)
}

func (intSerializer *intSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if (nullFlag & 1) == 1 {
			return 0, nil
		}
	}
	var val int32
	err := binary.Read(buffer, binary.BigEndian, &val)
	if err != nil {
		panic("read failed")
	}
	return val, nil
}

func (intSerializer *intSerializer) getDataType() dataType {
	return intType
}

// Format: 8-byte two’s complement integer.
type longSerializer struct{}

func (longSerializer *longSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return longSerializer.writeValue(value, buffer, writer, true)
}

func (longSerializer *longSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			return nil, errors.New("unexpected null value when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	switch value.(type) {
	case int:
		val := int64(value.(int))
		err := binary.Write(buffer, binary.BigEndian, val)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return buffer.Bytes(), nil
	default:
		err := binary.Write(buffer, binary.BigEndian, value)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return buffer.Bytes(), nil
	}
}

func (longSerializer *longSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return longSerializer.readValue(buffer, reader, true)
}

func (longSerializer *longSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if (nullFlag & 1) == 1 {
			return 0, nil
		}
	}
	var val int64
	err := binary.Read(buffer, binary.BigEndian, &val)
	if err != nil {
		panic("read failed")
	}
	return val, nil
}

func (longSerializer *longSerializer) getDataType() dataType {
	return longType
}

// Format: {length}{text_value}
type stringSerializer struct{}

func (stringSerializer *stringSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return stringSerializer.writeValue(value, buffer, writer, true)
}

func (stringSerializer *stringSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			return nil, errors.New("unexpected null value when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	val := value.(string)
	binary.Write(buffer, binary.BigEndian, int32(len(val)))
	_, err := buffer.WriteString(value.(string))
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (stringSerializer *stringSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return stringSerializer.readValue(buffer, reader, true)
}

func (stringSerializer *stringSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if (nullFlag & 1) == 1 {
			return "", nil
		}
	}
	// some better way to advance the pointer?
	buffer.Next(3)
	fmt.Println(buffer.Bytes())
	size, _ := binary.ReadUvarint(buffer)
	valBytes := make([]byte, size)
	fmt.Println(size)
	_, err := buffer.Read(valBytes)
	fmt.Println(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(valBytes), nil
}

func (stringSerializer *stringSerializer) getDataType() dataType {
	return stringType
}

// Format: 16 bytes representing the uuid.
type uuidSerializer struct {
}

func (uuidSerializer *uuidSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return uuidSerializer.writeValue(value, buffer, writer, true)
}

func (uuidSerializer *uuidSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			return nil, errors.New("unexpected null value when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	err := binary.Write(buffer, binary.BigEndian, value)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (uuidSerializer *uuidSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return uuidSerializer.readValue(buffer, reader, true)
}

func (uuidSerializer *uuidSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if (nullFlag & 1) == 1 {
			return uuid.Nil, nil
		}
	}

	valBytes := make([]byte, 16)
	_, err := buffer.Read(valBytes)
	if err != nil {
		return uuid.Nil, err
	}
	val, _ := uuid.FromBytes(valBytes)
	return val, nil
}

func (uuidSerializer *uuidSerializer) getDataType() dataType {
	return uuidType
}

// Format: {length}{item_0}...{item_n}
type mapSerializer struct{}

func (mapSerializer *mapSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return mapSerializer.writeValue(value, buffer, writer, true)
}

func (mapSerializer *mapSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			return nil, errors.New("unexpected null value when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	// problem with maps: currently it can only take type map[interface{}]interface{},
	// doesn't have a super object type like java.
	// would type assertion even work?
	val := value.(map[interface{}]interface{})
	binary.Write(buffer, binary.BigEndian, int32(len(val)))
	for k, v := range val {
		writer.writeObject(k, buffer)
		writer.writeObject(v, buffer)
	}
	return buffer.Bytes(), nil
}

func (mapSerializer *mapSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return mapSerializer.readValue(buffer, reader, true)
}

func (mapSerializer *mapSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if (nullFlag & 1) == 1 {
			return nil, nil
		}
	}
	buffer.Next(3)
	size, _ := binary.ReadUvarint(buffer)
	valMap := make(map[interface{}]interface{})
	fmt.Println("map size: ", size)
	for i := 0; i < int(size); i++ {
		valMap[reader.read(buffer)] = reader.read(buffer)
	}
	return valMap, nil
}

func (mapSerializer *mapSerializer) getDataType() dataType {
	return mapType
}
