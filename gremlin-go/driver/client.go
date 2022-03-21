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
	"golang.org/x/text/language"
)

// ClientSettings is used to modify a Client's settings on initialization.
type ClientSettings struct {
	TraversalSource string
	TransporterType TransporterType
	LogVerbosity    LogVerbosity
	Logger          Logger
	Language        language.Tag
	AuthInfo        *AuthInfo
	TlsConfig       *tls.Config
	session         string
	closed          bool
}

// Client is used to connect and interact with a Gremlin-supported server.
type Client struct {
	url             string
	traversalSource string
	logHandler      *logHandler
	transporterType TransporterType
	connection      *connection
	session         string
	closed          bool
}

// NewClient creates a Client and configures it with the given parameters.
func NewClient(url string, configurations ...func(settings *ClientSettings)) (*Client, error) {
	settings := &ClientSettings{
		TraversalSource: "g",
		TransporterType: Gorilla,
		LogVerbosity:    Info,
		Logger:          &defaultLogger{},
		Language:        language.English,
		AuthInfo:        &AuthInfo{},
		TlsConfig:       &tls.Config{},
		session:         "",
		closed:          false,
	}
	for _, configuration := range configurations {
		configuration(settings)
	}

	logHandler := newLogHandler(settings.Logger, settings.LogVerbosity, settings.Language)
	conn, err := createConnection(url, settings.AuthInfo, settings.TlsConfig, logHandler)
	if err != nil {
		return nil, err
	}
	client := &Client{
		url:             url,
		traversalSource: "g",
		logHandler:      logHandler,
		transporterType: settings.TransporterType,
		connection:      conn,
	}
	// TODO: PoolSize must be 1 on session mode. Implement after AN-980
	return client, nil
}

// Close closes the client via connection.
func (client *Client) Close() error {
	if client.IsClosed() {
		return nil
	}
	// If it is a Session, call closeSession
	if client.session != "" {
		_, err := client.closeSession()
		if err != nil {
			return err
		}
	}
	client.logHandler.logger.Logf(Info, "Closing Client with url %s", client.url)
	return client.connection.close()
}

func (client *Client) IsClosed() bool {
	return client.closed
}

// Submit submits a Gremlin script to the server and returns a ResultSet.
func (client *Client) Submit(message interface{}) (ResultSet, error) {
	// TODO AN-982: Obtain connection from pool of connections held by the client.
	client.logHandler.logf(Debug, submitStarted, message)
	args := map[string]interface{}{
		"gremlin": message,
		"aliases": map[string]interface{}{
			"g": client.traversalSource,
		},
	}
	var processor string
	var op string
	switch message.(type) {
	case bytecode:
		client.logHandler.logf(Debug, bytecodeReceived, message)
		op = "bytecode"
		processor = "traversal"
	case string:
		client.logHandler.logf(Debug, stringReceived, message)
		// TODO: implement after bindings (AN-1018).
		// args['bindings'] = bindings (Add Argument to func)
		processor = ""
		op = "eval"
	default:
		return nil, errors.New("message must either be a string or bytecode, neither was passed")
	}
	// If session
	if client.session != "" {
		args["session"] = client.session
		processor = "session"
	}
	// TODO: Get connection from pool after AN-982
	client.logHandler.logf(Debug, "processor='%s', op='%s', args='%s'", processor, op, args)
	request := makeRequest(op, processor, args)
	return client.connection.write(&request)
}

// submitBytecode submits bytecode to the server to execute and returns a ResultSet.
func (client *Client) submitBytecode(bytecode *bytecode) (ResultSet, error) {
	client.logHandler.logf(Debug, submitStartedBytecode, *bytecode)
	request := makeBytecodeRequest(bytecode, client.traversalSource)
	return client.connection.write(&request)
}

func (client *Client) closeSession() (ResultSet, error) {
	message := makeRequest("close", "session", map[string]interface{}{
		"session": client.session,
	})
	return client.connection.write(&message)
}
