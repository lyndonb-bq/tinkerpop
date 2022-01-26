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
	"reflect"

	"github.com/google/uuid"
)

// Version 1.0

// DataType graphbinary types
type DataType int

// DataType defined as constants
const (
	NullType    DataType = 0xFE
	IntType     DataType = 0x01
	LongType    DataType = 0x02
	StringType  DataType = 0x03
	DoubleType  DataType = 0x07
	FloatType   DataType = 0x08
	ListType    DataType = 0x09
	MapType     DataType = 0x0a
	UUIDType    DataType = 0x0c
	ByteType    DataType = 0x24
	ShortType   DataType = 0x26
	BooleanType DataType = 0x27
)

var nullBytes = []byte{NullType.getCodeByte(), 0x01}

func (dataType DataType) getCodeByte() byte {
	return byte(dataType)
}

func (dataType DataType) getCodeBytes() []byte {
	return []byte{dataType.getCodeByte()}
}

// GraphBinaryTypeSerializer interface for the different types of serializers
type GraphBinaryTypeSerializer interface {
	write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error)
	writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error)
	read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error)
	readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error)
	getDataType() DataType
}

// Format: 4-byte two’s complement integer.
type intSerializer struct{}

// Format: 8-byte two’s complement integer.
type longSerializer struct{}

// Format: {length}{text_value}
type stringSerializer struct{}

// Format: {length}{item_0}...{item_n}
type listSerializer struct{}

// Format: {length}{item_0}...{item_n}
type mapSerializer struct{}

// Format: 16 bytes representing the uuid.
type uuidSerializer struct{}

// graphBinaryWriter writes an object to byte array
type graphBinaryWriter struct {
	logHandler *logHandler
}

// graphBinaryReader reads a byte array into an object
type graphBinaryReader struct {
	logHandler *logHandler
}

const (
	valueFlagNull byte = 1
	valueFlagNone byte = 0
)

// gets the type of the serializer based on the value
func (writer *graphBinaryWriter) getSerializerToWrite(val interface{}) (GraphBinaryTypeSerializer, error) {
	switch val.(type) {
	case string:
		return &stringSerializer{}, nil
	case int64, int, uint32:
		return &longSerializer{}, nil
	case int32, int8, uint16:
		return &intSerializer{}, nil
	case uuid.UUID:
		return &uuidSerializer{}, nil
	case []string, []int, []int8, []int16, []int32, []uint, []uint8, []uint16, []uint32:
		return &listSerializer{}, nil
	default:
		switch reflect.TypeOf(val).Kind() {
		case reflect.Map:
			return &mapSerializer{}, nil
		default:
			writer.logHandler.log(Error, serializeDataTypeError)
			return nil, errors.New("unknown data type to serialize")
		}
	}
}

// gets the type of the serializer based on the DataType byte value
func (reader *graphBinaryReader) getSerializerToRead(typ byte) (GraphBinaryTypeSerializer, error) {
	switch typ {
	case StringType.getCodeByte():
		return &stringSerializer{}, nil
	case LongType.getCodeByte():
		return &longSerializer{}, nil
	case IntType.getCodeByte():
		return &intSerializer{}, nil
	case UUIDType.getCodeByte():
		return &uuidSerializer{}, nil
	case MapType.getCodeByte():
		return &mapSerializer{}, nil
	default:
		reader.logHandler.log(Error, deserializeDataTypeError)
		return nil, errors.New("unknown data type to deserialize")
	}
}

// Writes an object in fully-qualified format, containing {type_code}{type_info}{value_flag}{value}.
func (writer *graphBinaryWriter) write(valueObject interface{}, buffer *bytes.Buffer) (interface{}, error) {
	if valueObject == nil {
		// return Object of type "unspecified object null" with the value flag set to null.
		buffer.Write(nullBytes)
		return buffer.Bytes(), nil
	}

	serializer, err := writer.getSerializerToWrite(valueObject)
	if err != nil {
		return nil, err
	}
	buffer.Write(serializer.getDataType().getCodeBytes())
	message, err := serializer.write(valueObject, buffer, writer)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Writes a value without including type information.
func (writer *graphBinaryWriter) writeValue(value interface{}, buffer *bytes.Buffer, nullable bool) (interface{}, error) {
	if value == nil {
		if !nullable {
			writer.logHandler.log(Error, unexpectedNull)
			return nil, errors.New("unexpected null value to write when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	serializer, err := writer.getSerializerToWrite(value)
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

// Reads the type code, information and value of a given buffer with fully-qualified format.
func (reader *graphBinaryReader) read(buffer *bytes.Buffer) (interface{}, error) {
	typeCode, _ := buffer.ReadByte()
	if typeCode == NullType.getCodeByte() {
		check, _ := buffer.ReadByte()
		// check this
		if check != 1 {
			return nil, errors.New("read value flag")
		}
		return nil, nil
	}

	serializer, err := reader.getSerializerToRead(typeCode)
	if err != nil {
		return nil, err
	}
	val, err := serializer.read(buffer, reader)
	return val, err
}

func (reader *graphBinaryReader) readValue(buffer *bytes.Buffer, typ byte, nullable bool) (interface{}, error) {
	if buffer == nil {
		reader.logHandler.log(Error, nullInput)
		return nil, errors.New("input cannot be null")
	}
	typeCode, err := buffer.ReadByte()
	if err != nil {
		return nil, err
	}
	if typeCode != typ {
		reader.logHandler.logf(Error, unmatchedDataType)
		return nil, errors.New("datatype read from input buffer different from requested datatype")
	}
	serializer, _ := reader.getSerializerToRead(typ)
	val, _ := serializer.readValue(buffer, reader, nullable)
	return val, nil
}

func (intSerializer *intSerializer) getDataType() DataType {
	return IntType
}

func (intSerializer *intSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return intSerializer.writeValue(value, buffer, writer, true)
}

func (intSerializer *intSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			writer.logHandler.log(Error, unexpectedNull)
			return nil, errors.New("unexpected null value to write when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	//int8, uint16, int32
	var val int32
	switch value := value.(type) {
	case int8:
		val = int32(value)
	case uint16:
		val = int32(value)
	case int32:
		val = value
	default:
		writer.logHandler.log(Error, unmatchedDataType)
		return nil, errors.New("datatype read from input buffer different from requested datatype")
	}

	err := binary.Write(buffer, binary.BigEndian, val)
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
		if nullFlag == valueFlagNull {
			return 0, nil
		}
	}
	var val int32
	err := binary.Read(buffer, binary.BigEndian, &val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (longSerializer *longSerializer) getDataType() DataType {
	return LongType
}

func (longSerializer *longSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return longSerializer.writeValue(value, buffer, writer, true)
}

func (longSerializer *longSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			writer.logHandler.log(Error, unexpectedNull)
			return nil, errors.New("unexpected null value to write when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	// int, uint32, int64
	var val int64
	switch value := value.(type) {
	case int:
		val = int64(value)
	case uint32:
		val = int64(value)
	case int64:
		val = value
	default:
		writer.logHandler.log(Error, unmatchedDataType)
		return nil, errors.New("datatype read from input buffer different from requested datatype")
	}
	err := binary.Write(buffer, binary.BigEndian, val)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil

}

func (longSerializer *longSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return longSerializer.readValue(buffer, reader, true)
}

func (longSerializer *longSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if nullFlag == valueFlagNull {
			return 0, nil
		}
	}
	var val int64
	err := binary.Read(buffer, binary.BigEndian, &val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (stringSerializer *stringSerializer) getDataType() DataType {
	return StringType
}

func (stringSerializer *stringSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return stringSerializer.writeValue(value, buffer, writer, true)
}

func (stringSerializer *stringSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			writer.logHandler.log(Error, unexpectedNull)
			return nil, errors.New("unexpected null value to write when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	val := value.(string)
	err := binary.Write(buffer, binary.BigEndian, int32(len(val)))
	if err != nil {
		return nil, err
	}
	_, err = buffer.WriteString(value.(string))
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
		if nullFlag == valueFlagNull {
			return "", nil
		}
	}
	var size int32
	err := binary.Read(buffer, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	valBytes := make([]byte, size)
	_, err = buffer.Read(valBytes)
	if err != nil {
		return "", err
	}
	return string(valBytes), nil
}

func (mapSerializer *mapSerializer) getDataType() DataType {
	return MapType
}

func (mapSerializer *mapSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return mapSerializer.writeValue(value, buffer, writer, true)
}

func (mapSerializer *mapSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			writer.logHandler.log(Error, unexpectedNull)
			return nil, errors.New("unexpected null value to write when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Map {
		writer.logHandler.log(Error, notMap)
		return buffer.Bytes(), errors.New("did not get the expected map type as input")
	}

	keys := v.MapKeys()
	err := binary.Write(buffer, binary.BigEndian, int32(len(keys)))
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		convKey := k.Convert(v.Type().Key())
		// serialize k
		_, err := writer.write(k.Interface(), buffer)
		if err != nil {
			return nil, err
		}
		// serialize v.MapIndex(c_key)
		val := v.MapIndex(convKey)
		_, err = writer.write(val.Interface(), buffer)
		if err != nil {
			return nil, err
		}

	}
	return buffer.Bytes(), nil
}

func (mapSerializer *mapSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return mapSerializer.readValue(buffer, reader, true)
}

func (mapSerializer *mapSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if nullFlag == valueFlagNull {
			return nil, nil
		}
	}
	var size int32
	err := binary.Read(buffer, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	valMap := make(map[interface{}]interface{})
	for i := 0; i < int(size); i++ {
		key, err := reader.read(buffer)
		if err != nil {
			return nil, err
		}
		val, err := reader.read(buffer)
		if err != nil {
			return nil, err
		}
		valMap[key] = val
	}
	return valMap, nil
}

func (uuidSerializer *uuidSerializer) getDataType() DataType {
	return UUIDType
}

func (uuidSerializer *uuidSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return uuidSerializer.writeValue(value, buffer, writer, true)
}

func (uuidSerializer *uuidSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			writer.logHandler.log(Error, unexpectedNull)
			return nil, errors.New("unexpected null value to write when nullable is false")
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
		if nullFlag == valueFlagNull {
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

func (listSerializer *listSerializer) getDataType() DataType {
	return ListType
}

func (listSerializer *listSerializer) write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error) {
	return listSerializer.writeValue(value, buffer, writer, true)
}

func (listSerializer *listSerializer) writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error) {
	if value == nil {
		if !nullable {
			writer.logHandler.log(Error, unexpectedNull)
			return nil, errors.New("unexpected null value to write when nullable is false")
		}
		writer.writeValueFlagNull(buffer)
		return buffer.Bytes(), nil
	}

	if nullable {
		writer.writeValueFlagNone(buffer)
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		writer.logHandler.log(Error, notSlice)
		return buffer.Bytes(), errors.New("did not get the expected slice type as input")
	}

	// Need to convert to underlying type. Right now hardcode string.
	values := v.Interface().([]string)
	err := binary.Write(buffer, binary.BigEndian, int32(len(values)))
	if err != nil {
		return nil, err
	}
	if len(values) < 1 {
		return buffer.Bytes(), nil
	}
	subWriter, err := writer.getSerializerToWrite(values[0])
	if err != nil {
		return nil, err
	}
	for val := range values {
		_, err := subWriter.write(val, buffer, writer)
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func (listSerializer *listSerializer) read(buffer *bytes.Buffer, reader *graphBinaryReader) (interface{}, error) {
	return listSerializer.readValue(buffer, reader, true)
}

func (listSerializer *listSerializer) readValue(buffer *bytes.Buffer, reader *graphBinaryReader, nullable bool) (interface{}, error) {
	if nullable {
		nullFlag, _ := buffer.ReadByte()
		if nullFlag == valueFlagNull {
			return nil, nil
		}
	}
	var size int32
	err := binary.Read(buffer, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	valMap := make(map[interface{}]interface{})
	for i := 0; i < int(size); i++ {
		key, err := reader.read(buffer)
		if err != nil {
			return nil, err
		}
		val, err := reader.read(buffer)
		if err != nil {
			return nil, err
		}
		valMap[key] = val
	}
	return valMap, nil
}
