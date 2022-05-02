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
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/big"
	"reflect"
	"time"
)

var deserializers map[dataType]func(data *[]byte, i *int) interface{}

func readTemp(data *[]byte, i *int, len int) *[]byte {
	tmp := make([]byte, len)
	for j := 0; j < len; j++ {
		tmp[j] = (*data)[j+*i]
	}
	*i += len
	return &tmp
}

// Primitive
func readBoolean(data *[]byte, i *int) interface{} {
	return readByte(data, i) != 0
}

func readByte(data *[]byte, i *int) interface{} {
	*i++
	return (*data)[*i-1]
}

func readShort(data *[]byte, i *int) interface{} {
	*i += 2
	return int16((*data)[*i+1-2]) | int16((*data)[*i-2])<<8
}

func readInt(data *[]byte, i *int) interface{} {
	*i += 4
	return int32((*data)[*i+3-4]) | int32((*data)[*i+2-4])<<8 | int32((*data)[*i+1-4])<<16 | int32((*data)[*i-4])<<24
}

func readLong(data *[]byte, i *int) interface{} {
	*i += 8
	return int64((*data)[*i+7-8]) | int64((*data)[*i+6-8])<<8 | int64((*data)[*i+5-8])<<16 | int64((*data)[*i+4-8])<<24 |
		int64((*data)[*i+3-8])<<32 | int64((*data)[*i+2-8])<<40 | int64((*data)[*i+1-8])<<48 | int64((*data)[*i-8])<<56
}

func readUint32(data *[]byte, i *int) interface{} {
	*i += 4
	return uint32((*data)[*i+3-4]) | uint32((*data)[*i+2-4])<<8 | uint32((*data)[*i+1-4])<<16 | uint32((*data)[*i-4])<<24
}

func readUint64(data *[]byte, i *int) interface{} {
	*i += 8
	return uint64((*data)[*i+7-8]) | uint64((*data)[*i+6-8])<<8 | uint64((*data)[*i+5-8])<<16 | uint64((*data)[*i+4-8])<<24 |
		uint64((*data)[*i+3-8])<<32 | uint64((*data)[*i+2-8])<<40 | uint64((*data)[*i+1-8])<<48 | uint64((*data)[*i-8])<<56
}

func readFloat(data *[]byte, i *int) interface{} {
	return math.Float32frombits(readUint32(data, i).(uint32))
}

func readDouble(data *[]byte, i *int) interface{} {
	return math.Float64frombits(readUint64(data, i).(uint64))
}

func readBigInt(data *[]byte, i *int) interface{} {
	sz := readInt(data, i).(int32)
	b := readTemp(data, i, int(sz))

	var newBigInt = big.NewInt(0).SetBytes(*b)
	var one = big.NewInt(1)
	if len(*b) == 0 {
		return newBigInt
	}
	// If the first bit in the first element of the byte array is a 1, we need to interpret the byte array as a two's complement representation
	if (*b)[0]&0x80 == 0x00 {
		newBigInt.SetBytes(*b)
		return newBigInt
	}
	// Undo two's complement to byte array and set negative boolean to true
	length := uint((len(*b)*8)/8+1) * 8
	b2 := new(big.Int).Sub(newBigInt, new(big.Int).Lsh(one, length)).Bytes()

	// Strip the resulting 0xff byte at the start of array
	b2 = b2[1:]

	// Strip any redundant 0x00 byte at the start of array
	if b2[0] == 0x00 {
		b2 = b2[1:]
	}
	newBigInt = big.NewInt(0)
	newBigInt.SetBytes(b2)
	newBigInt.Neg(newBigInt)
	return newBigInt
}

func readString2(data *[]byte, i *int) interface{} {
	sz := int(readUint32(data, i).(uint32))
	if sz == 0 {
		return ""
	}

	tmp := make([]byte, sz)
	copy((*data)[*i:*i+sz], tmp[:])
	*i += sz
	return string(tmp)
}

func readDataType(data *[]byte, i *int) dataType {
	return dataType(readByte(data, i).(byte))
}

// Composite
func readList(data *[]byte, i *int) interface{} {
	// listEnd := time.Now()
	sz := readInt(data, i).(int32)
	var valList []interface{}
	for j := int32(0); j < sz; j++ {
		valList = append(valList, readFullyQualifiedNullable(data, i, true))
	}
	// listEnd := time.Now()
	// println("list: ", mapEnd.Sub(mapStart).Seconds())
	return valList
}

func readMap2(data *[]byte, i *int) interface{} {
	// mapStart := time.Now()
	sz := readUint32(data, i).(uint32)
	var mapData = make(map[interface{}]interface{})
	for j := uint32(0); j < sz; j++ {
		k := readFullyQualifiedNullable(data, i, true)
		v := readFullyQualifiedNullable(data, i, true)
		if k == nil {
			mapData[nil] = v
		} else {
			switch reflect.TypeOf(k).Kind() {
			case reflect.Map:
				mapData[&k] = v
				break
			case reflect.Slice:
				mapData[fmt.Sprint(k)] = v
				break
			default:
				mapData[k] = v
				break
			}
		}
	}
	// mapEnd := time.Now()
	// println("map: ", mapEnd.Sub(mapStart).Seconds())
	return mapData //, nil
}

func readMapUnqualified2(data *[]byte, i *int) interface{} {
	sz := readUint32(data, i).(uint32)
	var mapData = make(map[string]interface{})
	for j := uint32(0); j < sz; j++ {
		keyDataType := readDataType(data, i)
		if keyDataType != stringType {
			return nil // , newError(err0703ReadMapNonStringKeyError)
		}

		// Skip nullable, key must be present
		*i++

		k := readString2(data, i).(string)
		mapData[k] = readFullyQualifiedNullable(data, i, true)
	}
	return mapData //, nil
}

func readSet(data *[]byte, i *int) interface{} {
	return NewSimpleSet(readList(data, i).([]interface{})...)
}

func readUuid(data *[]byte, i *int) interface{} {
	id, _ := uuid.FromBytes(*readTemp(data, i, 16))
	return id
}

func timeReader2(data *[]byte, i *int) interface{} {
	return time.UnixMilli(readLong(data, i).(int64))
}

func durationReader2(data *[]byte, i *int) interface{} {
	return time.Duration(readLong(data, i).(int64)*int64(time.Second) + int64(readInt(data, i).(int32)))
}

// Graph

// {fully qualified id}{unqualified label}
func vertexReader2(data *[]byte, i *int) interface{} {
	return vertexReaderNullByte(data, i, true)
}

// {fully qualified id}{unqualified label}{[unused null byte]}
func vertexReaderNullByte(data *[]byte, i *int, unusedByte bool) interface{} {
	v := new(Vertex)
	v.Id = readFullyQualifiedNullable(data, i, true)
	v.Label = readUnqualified(data, i, stringType, false).(string)
	if unusedByte {
		*i++
	}
	return v
}

// {fully qualified id}{unqualified label}{in vertex w/o null byte}{out vertex}{unused null byte}{unused null byte}
func edgeReader2(data *[]byte, i *int) interface{} {
	e := new(Edge)
	e.Id = readFullyQualifiedNullable(data, i, true)
	e.Label = readUnqualified(data, i, stringType, false).(string)
	e.InV = *vertexReaderNullByte(data, i, false).(*Vertex)
	e.OutV = *vertexReaderNullByte(data, i, false).(*Vertex)
	*i += 2
	return e
}

// {unqualified key}{fully qualified value}{null byte}
func propertyReader2(data *[]byte, i *int) interface{} {
	p := new(Property)
	p.Key = readUnqualified(data, i, stringType, false).(string)
	p.Value = readFullyQualifiedNullable(data, i, true)
	*i++
	return p
}

// {fully qualified id}{unqualified label}{fully qualified value}{null byte}{null byte}
func vertexPropertyReader2(data *[]byte, i *int) interface{} {
	vp := new(VertexProperty)
	vp.Id = readFullyQualifiedNullable(data, i, true)
	vp.Label = readUnqualified(data, i, stringType, false).(string)
	vp.Value = readFullyQualifiedNullable(data, i, true)
	*i += 2
	return vp
}

// {list of set of strings}{list of fully qualified objects}
func pathReader2(data *[]byte, i *int) interface{} {
	path := new(Path)
	newLabels := readFullyQualifiedNullable(data, i, true)
	for _, param := range newLabels.([]interface{}) {
		path.Labels = append(path.Labels, param.(*SimpleSet))
	}
	path.Objects = readFullyQualifiedNullable(data, i, true).([]interface{})
	return path
}

// {bulk int}{fully qualified value}
func traverserReader(data *[]byte, i *int) interface{} {
	traverser := new(Traverser)
	traverser.bulk = readLong(data, i).(int64)
	traverser.value = readFullyQualifiedNullable(data, i, true)
	return traverser
}

// {int32 length}{fully qualified item_0}{int64 repetition_0}...{fully qualified item_n}{int64 repetition_n}
func bulkSetReader2(data *[]byte, i *int) interface{} {
	sz := int(readInt(data, i).(int32))
	var valList []interface{}
	for j := 0; j < sz; j++ {
		val := readFullyQualifiedNullable(data, i, true)
		rep := readLong(data, i).(int64)
		for k := 0; k < int(rep); k++ {
			valList = append(valList, val)
		}
	}
	return valList
}

// {type code (always string so ignore)}{nil code (always false so ignore)}{int32 size}{string enum}
func enumReader2(data *[]byte, i *int) interface{} {
	*i += 2
	return readString2(data, i)
}

// {unqualified key}{fully qualified value}
func bindingReader2(data *[]byte, i *int) interface{} {
	b := new(Binding)
	b.Key = readUnqualified(data, i, stringType, false).(string)
	b.Value = readFullyQualifiedNullable(data, i, true)
	return b
}

func deserializeMessage2(message []byte) (response, error) {
	// Skip version and nullable byte.

	if deserializers == nil || len(deserializers) == 0 {
		deserializers = map[dataType]func(data *[]byte, i *int) interface{}{
			// Primitive
			booleanType:    readBoolean,
			byteType:       readByte,
			shortType:      readShort,
			intType:        readInt,
			longType:       readLong,
			bigIntegerType: readBigInt,
			floatType:      readFloat,
			doubleType:     readDouble,
			stringType:     readString2,

			// Composite
			listType: readList,
			mapType:  readMap2,
			setType:  readSet,
			uuidType: readUuid,

			// Date Time
			dateType:      timeReader2,
			timestampType: timeReader2,
			durationType:  durationReader2,

			// Graph
			traverserType:      traverserReader,
			vertexType:         vertexReader2,
			edgeType:           edgeReader2,
			propertyType:       propertyReader2,
			vertexPropertyType: vertexPropertyReader2,
			pathType:           pathReader2,
			bulkSetType:        bulkSetReader2,
			tType:              enumReader2,
			directionType:      enumReader2,
			bindingType:        bindingReader2,
		}
	}
	i := 2
	var msg response
	msg.responseID = readUuid(&message, &i).(uuid.UUID)
	msg.responseStatus.code = uint16(readUint32(&message, &i).(uint32) & 0xFF)
	if isMessageValid := readBoolean(&message, &i).(bool); isMessageValid {
		msg.responseStatus.message = readString2(&message, &i).(string)
	}
	msg.responseStatus.attributes = readMapUnqualified2(&message, &i).(map[string]interface{})
	msg.responseResult.meta = readMapUnqualified2(&message, &i).(map[string]interface{})
	msg.responseResult.data = readFullyQualifiedNullable(&message, &i, true)
	return msg, nil
}

func readUnqualified(data *[]byte, i *int, dataTyp dataType, nullable bool) interface{} {
	if nullable && readBoolean(data, i).(bool) {
		return nil
	}
	return deserializers[dataTyp](data, i)
}

func readFullyQualifiedNullable(data *[]byte, i *int, nullable bool) interface{} {
	dataTyp := readDataType(data, i)
	if dataTyp == nullType {
		return nil
	} else if nullable {
		if readByte(data, i).(byte) != byte(0) {
			return nil
		}
	}
	return deserializers[dataTyp](data, i)
}
