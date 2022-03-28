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

type errorCode string

var count int = 34

const (
	// connection.go errors
	err0101ConnectionCloseError errorCode = "E0101_CONNECTION_CLOSE_ERROR"

	// driverRemoteConnection.go errors
	err0201CreateSessionMultipleIdsError errorCode = "E0201_DRIVER_REMOTE_CONNECTION_CREATESESSION_MULTIPLE_UUIDS_ERROR"
	err0202CreateSessionFromSessionError errorCode = "E0202_DRIVER_REMOTE_CONNECTION_CREATESESSION_SESSION_FROM_SESSION_ERROR"

	// graph.go errors
	err0301GetPathObjectInvalidPathUnequalLengthsError errorCode = "E0301_GRAPH_GETPATHOBJECT_UNEQUAL_LABELS_OBJECTS_LENGTH_ERROR"
	err0302GetPathObjectInvalidPathNonStringLabelError errorCode = "E0302_GRAPH_GETPATHOBJECT_NON_STRING_VALUE_IN_LABELS_ERROR"
	err0303GetPathNoLabelFoundError                    errorCode = "E0303_GRAPH_GETPATHOBJECT_NO_LABEL_ERROR"

	// graphBinary.go errors
	err0401WriteTypeValueUnexpectedNullError    errorCode = "E0401_GRAPH_BINARY_WRITETYPEVALUE_UNEXPECTED_NULL_ERROR"
	err0402BytecodeWriterError                  errorCode = "E0402_GRAPH_BINARY_WRITER_BYTECODE_ERROR"
	err0403WriteTypeUnexpectedNullError         errorCode = "E0403_GRAPH_BINARY_WRITETYPE_UNEXPECTED_NULL_ERROR"
	err0404ReadNullTypeError                    errorCode = "E0404_GRAPH_BINARY_READ_NULLTYPE_ERROR"
	err0405ReadValueInvalidNullInputError       errorCode = "E0405_GRAPH_BINARY_READ_NULLTYPE_ERROR"
	err0406EnumReaderInvalidTypeError           errorCode = "E0406_GRAPH_BINARY_INVALID_TYPE_ERROR"
	err0407GetSerializerToWriteUnknownTypeError errorCode = "E0407_GRAPH_BINARY_GETSERIALIZERTOWRITE_UNKNOWN_TYPE_ERROR"
	err0407GetSerializerToReadUnknownTypeError  errorCode = "E0408_GRAPH_BINARY_GETSERIALIZERTOREAD_UNKNOWN_TYPE_ERROR"

	// protocol.go errors
	err0501ResponseHandlerResultSetNotCreatedError errorCode = "E0501_PROTOCOL_RESPONSEHANDLER_NO_RESULTSET_ON_DATA_RECEIVE"
	err0502ResponseHandlerReadLoopError            errorCode = "E0502_PROTOCOL_RESPONSEHANDLER_READ_LOOP_ERROR"
	err0502ResponseHandlerAuthError                errorCode = "E0503_PROTOCOL_RESPONSEHANDLER_AUTH_ERROR"

	// result.go errors
	err0601ResultNotVertexError         errorCode = "E0601_RESULT_NOT_VERTEX_ERROR"
	err0602ResultNotEdgeError           errorCode = "E0602_RESULT_NOT_EDGE_ERROR"
	err0603ResultNotElementError        errorCode = "E0603_RESULT_NOT_ELEMENT_ERROR"
	err0604ResultNotPathError           errorCode = "E0604_RESULT_NOT_PATH_ERROR"
	err0605ResultNotPropertyError       errorCode = "E0605_RESULT_NOT_PROPERTY_ERROR"
	err0606ResultNotVertexPropertyError errorCode = "E0606_RESULT_NOT_VERTEX_PROPERTY_ERROR"
	err0607ResultNotTraverserError      errorCode = "E0607_RESULT_NOT_TRAVERSER_ERROR"
	err0608ResultNotSliceError          errorCode = "E0608_RESULT_NOT_SLICE_ERROR"

	// serializer.go errors
	err0701ReadMapNullKeyError          errorCode = "E0701_SERIALIZER_READMAP_NULL_KEY_ERROR"
	err0702ReadMapNullKeyError          errorCode = "E0702_SERIALIZER_READMAP_NULL_KEY_ERROR"
	err0703ReadMapNonStringKeyError     errorCode = "E0703_SERIALIZER_READMAP_NON_STRING_KEY_ERROR"
	err0704ConvertArgsNoSerializerError errorCode = "E0704_SERIALIZER_CONVERTARGS_NO_SERIALIZER_ERROR"

	// transporterFactory.go errors
	err0801GetTransportLayerNoTypeError errorCode = "E0801_TRANSPORTERFACTORY_GETTRANSPORTLAYER_NO_TYPE_ERROR"

	// traversal.go errors
	err0901ToListAnonTraversalError  errorCode = "E0901_TRAVERSAL_TOLIST_ANON_TRAVERSAL_ERROR"
	err0902IterateAnonTraversalError errorCode = "E0902_TRAVERSAL_ITERATE_ANON_TRAVERSAL_ERROR"
	err0903NextNoResultsLeftError    errorCode = "E0903_TRAVERSAL_NEXT_NO_RESULTS_LEFT_ERROR"
)
