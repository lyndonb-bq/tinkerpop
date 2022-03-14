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
	"crypto/tls"
	"golang.org/x/text/language"
	"net/http"
)

// AuthInfo is an option struct that allows authentication information to be specified.
// Authentication can be provided via http.Header Header directly.
// Basic authentication can also be used via the BasicAuthInfo function.
type AuthInfo struct {
	Header   http.Header
	Username string
	Password string
}

// getHeader provides a safe way to get a header from the AuthInfo even if it is nil.
// This way we don't need any additional logic in the transport layer.
func (authInfo *AuthInfo) getHeader() http.Header {
	if authInfo == nil {
		return nil
	} else {
		return authInfo.Header
	}
}

// getUseBasicAuth provides a safe way to get a if basic auth info is available from the AuthInfo even if it is nil.
// This way we don't need any additional logic in the transport layer.
func (authInfo *AuthInfo) getUseBasicAuth() bool {
	return authInfo != nil && authInfo.Username != "" && authInfo.Password != ""
}

// BasicAuthInfo provides a way to generate AuthInfo. Enter username and password and get the AuthInfo back.
func BasicAuthInfo(username string, password string) *AuthInfo {
	return &AuthInfo{Username: username, Password: password}
}

// DriverRemoteConnectionSettings are used to configure the DriverRemoteConnection.
type DriverRemoteConnectionSettings struct {
	TraversalSource string
	TransporterType TransporterType
	LogVerbosity    LogVerbosity
	Logger          Logger
	Language        language.Tag
	AuthInfo        *AuthInfo
	TlsConfig       *tls.Config

	// TODO: Figure out exact extent of configurability for these and expose appropriate types/helpers
	Protocol   protocol
	Serializer serializer
}

// DriverRemoteConnection is a remote connection.
type DriverRemoteConnection struct {
	client *Client
}

// NewDriverRemoteConnection creates a new DriverRemoteConnection.
// If no custom connection settings are passed in, a connection will be created with "g" as the default TraversalSource,
// Gorilla as the default Transporter, Info as the default LogVerbosity, a default logger stuct, and English and as the
// default language
func NewDriverRemoteConnection(
	url string,
	configurations ...func(settings *DriverRemoteConnectionSettings)) (*DriverRemoteConnection, error) {
	settings := &DriverRemoteConnectionSettings{
		TraversalSource: "g",
		TransporterType: Gorilla,
		LogVerbosity:    Info,
		Logger:          &defaultLogger{},
		Language:        language.English,
		AuthInfo:        &AuthInfo{},
		TlsConfig:       &tls.Config{},

		// TODO: Figure out exact extent of configurability for these and expose appropriate types/helpers
		Protocol:   nil,
		Serializer: nil,
	}
	for _, configuration := range configurations {
		configuration(settings)
	}

	logHandler := newLogHandler(settings.Logger, settings.LogVerbosity, settings.Language)
	connection, err := createConnection(url, settings.AuthInfo, settings.TlsConfig, logHandler)
	if err != nil {
		return nil, err
	}

	client := &Client{
		url:             url,
		authInfo:        settings.AuthInfo,
		tlsConfig:       settings.TlsConfig,
		transporterType: settings.TransporterType,
		logHandler:      logHandler,
		connection:      connection,
	}

	return &DriverRemoteConnection{client: client}, nil
}

// Close closes the DriverRemoteConnection.
func (driver *DriverRemoteConnection) Close() error {
	return driver.client.Close()
}

// Submit sends a string traversal to the server.
func (driver *DriverRemoteConnection) Submit(traversalString string) (ResultSet, error) {
	return driver.client.Submit(traversalString)
}

// submitBytecode sends a bytecode traversal to the server.
func (driver *DriverRemoteConnection) submitBytecode(bytecode *bytecode) (ResultSet, error) {
	return driver.client.submitBytecode(bytecode)
}

// TODO: Bytecode, OptionsStrategy, RequestOptions
//func extractRequestOptions(bytecode Bytecode) RequestOptions {
//	var optionsStrategy OptionsStrategy = nil
//	for _, instruction := range bytecode.sourceInstructions {
//		if instruction[0] == "withStrategies" {
//			_, isOptionsStrategy := instruction[1].(OptionsStrategy)
//			if isOptionsStrategy {
//				optionsStrategy = instruction
//				break
//			}
//		}
//	}
//
//	var requestOptions RequestOptions = nil
//	if optionsStrategy != nil {
//		allowedKeys := []string{'evaluationTimeout', 'scriptEvaluationTimeout', 'batchSize', 'requestId', 'userAgent'}
//		requestOptions := make(map[string]string)
//		for _, allowedKey := range allowedKeys {
//			if isAllowedKeyInConfigs(allowedKey, optionsStrategy[1].configuration) {
//				requestOptions[allowedKey] = optionsStrategy[1].configuration[allowedKey]
//			}
//		}
//	}
//	return requestOptions
//}

//func isAllowedKeyInConfigs(allowedKey string, configs []string) bool {
//	for _, config := range configs {
//		if allowedKey == config {
//			return true
//		}
//	}
//	return false
//}
