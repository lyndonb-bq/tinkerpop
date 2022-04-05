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
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"runtime"
	"time"
)

// DriverRemoteConnectionSettings are used to configure the DriverRemoteConnection.
type DriverRemoteConnectionSettings struct {
	TraversalSource   string
	TransporterType   TransporterType
	LogVerbosity      LogVerbosity
	Logger            Logger
	Language          language.Tag
	AuthInfo          *AuthInfo
	TlsConfig         *tls.Config
	KeepAliveInterval time.Duration
	WriteDeadline     time.Duration
	// Minimum amount of concurrent active traversals on a connection to trigger creation of a new connection
	NewConnectionThreshold int
	// Maximum number of concurrent connections. Default: number of runtime processors
	MaximumConcurrentConnections int
	Session                      string

	// TODO: Figure out exact extent of configurability for these and expose appropriate types/helpers
	Protocol   protocol
	Serializer serializer
}

// DriverRemoteConnection is a remote connection.
type DriverRemoteConnection struct {
	client          *Client
	spawnedSessions []*DriverRemoteConnection
	isClosed        bool
}

// NewDriverRemoteConnection creates a new DriverRemoteConnection.
// If no custom connection settings are passed in, a connection will be created with "g" as the default TraversalSource,
// Gorilla as the default Transporter, Info as the default LogVerbosity, a default logger stuct, and English and as the
// default language
func NewDriverRemoteConnection(
	url string,
	configurations ...func(settings *DriverRemoteConnectionSettings)) (*DriverRemoteConnection, error) {
	settings := &DriverRemoteConnectionSettings{
		TraversalSource:              "g",
		TransporterType:              Gorilla,
		LogVerbosity:                 Info,
		Logger:                       &defaultLogger{},
		Language:                     language.English,
		AuthInfo:                     &AuthInfo{},
		TlsConfig:                    &tls.Config{},
		KeepAliveInterval:            keepAliveIntervalDefault,
		WriteDeadline:                writeDeadlineDefault,
		NewConnectionThreshold:       defaultNewConnectionThreshold,
		MaximumConcurrentConnections: runtime.NumCPU(),
		Session:                      "",

		// TODO: Figure out exact extent of configurability for these and expose appropriate types/helpers
		Protocol:   nil,
		Serializer: nil,
	}
	for _, configuration := range configurations {
		configuration(settings)
	}

	logHandler := newLogHandler(settings.Logger, settings.LogVerbosity, settings.Language)
	if settings.Session != "" {
		logHandler.log(Info, sessionDetected)
		settings.MaximumConcurrentConnections = 1
	}

	pool, err := newLoadBalancingPool(url, logHandler, settings.AuthInfo, settings.TlsConfig, settings.KeepAliveInterval,
		settings.WriteDeadline, settings.NewConnectionThreshold, settings.MaximumConcurrentConnections)
	if err != nil {
		if err != nil {
			logHandler.logf(Error, logErrorGeneric, "NewDriverRemoteConnection", err.Error())
		}
		return nil, err
	}

	client := &Client{
		url:             url,
		traversalSource: settings.TraversalSource,
		transporterType: settings.TransporterType,
		logHandler:      logHandler,
		connections:     pool,
		session:         settings.Session,
	}

	return &DriverRemoteConnection{client: client, isClosed: false}, nil
}

// Close closes the DriverRemoteConnection.
// Errors if any will be logged
func (driver *DriverRemoteConnection) Close() {
	// If DriverRemoteConnection has spawnedSessions then they must be closed as well.
	if len(driver.spawnedSessions) > 0 {
		driver.client.logHandler.logf(Info, closingSpawnedSessions, driver.client.url)
		for _, session := range driver.spawnedSessions {
			session.Close()
		}
		driver.spawnedSessions = driver.spawnedSessions[:0]
	}

	if driver.isSession() {
		driver.client.logHandler.logf(Info, closeSession, driver.client.url, driver.client.session)
	} else {
		driver.client.logHandler.logf(Info, closeDriverRemoteConnection, driver.client.url)
	}
	driver.client.Close()
	driver.isClosed = true
}

// Submit sends a string traversal to the server.
func (driver *DriverRemoteConnection) Submit(traversalString string) (ResultSet, error) {
	result, err := driver.client.Submit(traversalString)
	if err != nil {
		driver.client.logHandler.logf(Error, logErrorGeneric, "Driver.Submit()", err.Error())
	}
	return result, err
}

// submitBytecode sends a bytecode traversal to the server.
func (driver *DriverRemoteConnection) submitBytecode(bytecode *bytecode) (ResultSet, error) {
	if driver.isClosed {
		return nil, newError(err0203SubmitBytecodeToClosedConnectionError)
	}
	return driver.client.submitBytecode(bytecode)
}

func (driver *DriverRemoteConnection) isSession() bool {
	return driver.client.session != ""
}

// CreateSession generates a new Session. sessionId stores the optional UUID param. It can be used to create a Session with a specific UUID.
func (driver *DriverRemoteConnection) CreateSession(sessionId ...string) (*DriverRemoteConnection, error) {
	if len(sessionId) > 1 {
		return nil, newError(err0201CreateSessionMultipleIdsError)
	} else if driver.isSession() {
		return nil, newError(err0202CreateSessionFromSessionError)
	}

	driver.client.logHandler.log(Info, creatingSessionConnection)
	drc, err := NewDriverRemoteConnection(driver.client.url, func(settings *DriverRemoteConnectionSettings) {
		settings.TraversalSource = driver.client.traversalSource
		if len(sessionId) == 1 {
			settings.Session = sessionId[0]
		} else {
			settings.Session = uuid.New().String()
		}
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

func (driver *DriverRemoteConnection) commit() (ResultSet, error) {
	bc := &bytecode{}
	bc.addSource("tx", "commit")
	return driver.submitBytecode(bc)
}

func (driver *DriverRemoteConnection) rollback() (ResultSet, error) {
	bc := &bytecode{}
	bc.addSource("tx", "rollback")
	return driver.submitBytecode(bc)
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
