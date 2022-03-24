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
	"errors"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

// DriverRemoteConnectionSettings are used to configure the DriverRemoteConnection.
type DriverRemoteConnectionSettings struct {
	TraversalSource string
	TransporterType TransporterType
	LogVerbosity    LogVerbosity
	Logger          Logger
	Language        language.Tag
	AuthInfo        *AuthInfo
	TlsConfig       *tls.Config
	Session         string

	// TODO: Figure out exact extent of configurability for these and expose appropriate types/helpers
	Protocol   protocol
	Serializer serializer
}

// DriverRemoteConnection is a remote connection.
type DriverRemoteConnection struct {
	client          *Client
	spawnedSessions []*DriverRemoteConnection
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
		Session:         "",

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
		traversalSource: settings.TraversalSource,
		transporterType: settings.TransporterType,
		logHandler:      logHandler,
		connection:      connection,
		session:         settings.Session,
	}

	return &DriverRemoteConnection{client: client}, nil
}

// Close closes the DriverRemoteConnection.
func (driver *DriverRemoteConnection) Close() error {
	// If DriverRemoteConnection has spawnedSessions then they must be closed as well.
	if len(driver.spawnedSessions) > 0 {
		driver.client.logHandler.logf(Info, closingSpawnedSessions, driver.client.url)
		flag := false
		for _, session := range driver.spawnedSessions {
			err := session.Close()
			if err != nil {
				driver.client.logHandler.logf(Error, closeSessionError, err)
				flag = true
			}
		}
		if flag {
			return errors.New("there was at least one error while closing spawned sessions. Check error log for more info")
		}
		return nil
	}

	if driver.isSession() {
		driver.client.logHandler.logf(Info, closeSession, driver.client.url, driver.client.session)
	} else {
		driver.client.logHandler.logf(Info, closeDriverRemoteConnection, driver.client.url)
	}
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

// isSessionBound returns True if a DriverRemoteConnection is a Session
func (driver *DriverRemoteConnection) isSession() bool {
	return driver.client.session != ""
}

// CreateSession generates a new Session. sessionId stores the optional UUID param. It can be used to create a Session with a specific UUID.
func (driver *DriverRemoteConnection) CreateSession(sessionId ...string) (*DriverRemoteConnection, error) {
	if len(sessionId) > 1 {
		return nil, errors.New("more than one Session ID specified. Cannot create Session with multiple UUIDs")
	} else if driver.isSession() {
		return nil, errors.New("connection is already bound to a Session - child sessions are not allowed")
	}

	driver.client.logHandler.logger.Log(Info, "creating Session based connection")
	sid := uuid.New().String()
	// If there is a sessionId provided, overwrite uuid.
	if len(sessionId) == 1 {
		sid = sessionId[0]
	}
	drc, err := NewDriverRemoteConnection(driver.client.url, func(settings *DriverRemoteConnectionSettings) {
		settings.Session = sid
	})
	if err != nil {
		return nil, err
	}
	driver.spawnedSessions = append(driver.spawnedSessions, drc)
	return drc, nil
}

func (driver *DriverRemoteConnection) GetSessionId() string {
	return driver.client.session
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
