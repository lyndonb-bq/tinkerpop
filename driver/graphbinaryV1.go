package gremlingo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type DataType int

const (
	NULL   DataType = 0xFE
	INT    DataType = 0x01
	LONG   DataType = 0x02
	DOUBLE DataType = 0x07
	FLOAT  DataType = 0x08

	STRING DataType = 0x03
	UUID   DataType = 0x0c
	MAP    DataType = 0x0a
)

var nullBytes = []byte{NULL.getCodeByte(), 0x01}

func (dataType DataType) getCode() int {
	return int(dataType)
}

func (dataType DataType) getCodeByte() byte {
	return byte(dataType)
}

func (dataType DataType) getCodeBytes() []byte {
	return []byte{dataType.getCodeByte()}
}

type graphBinarySerializer interface {
	// change type of value for each specific serializer?
	write(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter) ([]byte, error)
	writeValue(value interface{}, buffer *bytes.Buffer, writer *graphBinaryWriter, nullable bool) ([]byte, error)
	getDataType() DataType
	// read()
	// readValue()
}

// map for different types of serializers used in writing to graphbinary
var graphBinSerializerMap = map[interface{}]graphBinarySerializer{
	// simple types
	"int32":  &intSerializer{},
	"int64":  &longSerializer{},
	"int":    &longSerializer{}, // go int can be 32 or 64-bit
	"string": &stringSerializer{},
	"UUID":   &uuidSerializer{},

	// complex types
	reflect.Map: &mapSerializer{},
}

type graphBinaryWriter struct {
}

const (
	valueFlagNull byte = 1
	valueFlagNone byte = 0
)

// instead of using the map, use type assertion?
// these are basic types
func (writer *graphBinaryWriter) getSerializer(t reflect.Type) graphBinarySerializer {
	switch t.Kind() {
	case reflect.String:
		return &stringSerializer{}
	case reflect.Int64, reflect.Int, reflect.Uint32:
		return &longSerializer{}
	case reflect.Int32, reflect.Int8, reflect.Uint16:
		return &intSerializer{}
	case reflect.Array:
		if t.Name() == "UUID" {
			return &uuidSerializer{}
		} else {
			return nil
		}
	case reflect.Map:
		return &mapSerializer{}
	default:
		return nil
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
	//objectTypeName := reflect.TypeOf(valueObject).Name()
	//objectTypeKind := reflect.TypeOf(valueObject).Kind()
	//if serializer, found := graphBinSerializerMap[objectTypeName]; found {
	//	buffer.Write(serializer.getDataType().getCodeBytes())
	//	message, _ := serializer.write(valueObject, buffer, writer)
	//	return message
	//} else if serializer, found := graphBinSerializerMap[objectTypeKind]; found {
	//	buffer.Write(serializer.getDataType().getCodeBytes())
	//	message, _ := serializer.write(valueObject, buffer, writer)
	//	return message
	//} else {
	//	return nil
	//}

	serializer := writer.getSerializer(reflect.TypeOf(valueObject))
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

	//objectTypeName := reflect.TypeOf(value).Name()
	//objectTypeKind := reflect.TypeOf(value).Kind()
	//if serializer, ok := graphBinSerializerMap[objectTypeName]; ok {
	//	buffer.Write(serializer.getDataType().getCodeBytes())
	//	message, _ := serializer.writeValue(value, buffer, writer, nullable)
	//	return message, nil
	//} else if serializer, found := graphBinSerializerMap[objectTypeKind]; found {
	//	buffer.Write(serializer.getDataType().getCodeBytes())
	//	message, _ := serializer.writeValue(value, buffer, writer, nullable)
	//	return message, nil
	//} else {
	//	return nil, errors.New("encountered unknown data type")
	//}

	serializer := writer.getSerializer(reflect.TypeOf(value))
	buffer.Write(serializer.getDataType().getCodeBytes())
	message, _ := serializer.writeValue(value, buffer, writer, nullable)
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

type simpleTypeSerializer struct{}

// Format: 4-byte two’s complement integer.
type intSerializer struct{}

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

func (intSerializer *intSerializer) getDataType() DataType {
	return INT
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

func (longSerializer *longSerializer) getDataType() DataType {
	return LONG
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

func (stringSerializer *stringSerializer) getDataType() DataType {
	return STRING
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

func (uuidSerializer *uuidSerializer) getDataType() DataType {
	return UUID
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

func (mapSerializer *mapSerializer) getDataType() DataType {
	return MAP
}
