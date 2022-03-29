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
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"strconv"
	"sync"
	"testing"
)

// Arbitrarily high value to use to not trigger creation of new connections
const newConnectionThreshold = 100

func getPoolForTesting() *loadBalancingPool {
	return &loadBalancingPool{
		url:                    "",
		authInfo:               nil,
		tlsConfig:              nil,
		logHandler:             newLogHandler(&defaultLogger{}, Info, language.English),
		newConnectionThreshold: newConnectionThreshold,
		connections:            nil,
		loadBalanceLock:        sync.Mutex{},
	}
}

func TestConnectionPool(t *testing.T) {
	t.Run("loadBalancingPool", func(t *testing.T) {
		logHandler := newLogHandler(&defaultLogger{}, Info, language.English)
		mockConnection1 := &connection{
			logHandler: logHandler,
			protocol:   nil,
			results:    nil,
			state:      established,
		}
		mockConnection2 := &connection{
			logHandler: logHandler,
			protocol:   nil,
			results:    nil,
			state:      established,
		}
		mockConnection3 := &connection{
			logHandler: logHandler,
			protocol:   nil,
			results:    nil,
			state:      established,
		}

		smallMap := make(map[string]ResultSet)
		bigMap := make(map[string]ResultSet)
		for i := 1; i < 4; i++ {
			bigMap[strconv.Itoa(i)] = nil
			if i < 3 {
				smallMap[strconv.Itoa(i)] = nil
			}
		}
		mockConnection1.results = bigMap
		mockConnection2.results = smallMap
		mockConnection3.results = bigMap

		t.Run("getLeastUsedConnection", func(t *testing.T) {
			t.Run("getting the least used connection", func(t *testing.T) {
				pool := getPoolForTesting()
				connections := []*connection{mockConnection1, mockConnection2, mockConnection3}
				pool.connections = connections

				connection, err := pool.getLeastUsedConnection()
				assert.Nil(t, err)
				assert.Equal(t, mockConnection2, connection)
			})

			t.Run("purge non-established connections", func(t *testing.T) {
				pool := getPoolForTesting()
				nonEstablished := &connection{
					logHandler: logHandler,
					protocol:   nil,
					results:    nil,
					state:      closed,
				}
				connections := []*connection{nonEstablished, mockConnection2}
				pool.connections = connections

				connection, err := pool.getLeastUsedConnection()
				assert.Nil(t, err)
				assert.Equal(t, mockConnection2, connection)
				assert.Len(t, pool.connections, 1)
			})

			t.Run("purge non-used connections", func(t *testing.T) {
				pool := getPoolForTesting()
				empty := make(map[string]ResultSet)
				emptyConn1 := &connection{
					logHandler: logHandler,
					protocol:   nil,
					results:    empty,
					state:      established,
				}
				emptyConn2 := &connection{
					logHandler: logHandler,
					protocol:   nil,
					results:    empty,
					state:      established,
				}
				connections := []*connection{emptyConn1, emptyConn2}
				pool.connections = connections

				connection, err := pool.getLeastUsedConnection()
				assert.Nil(t, err)
				assert.NotNil(t, connection)
				assert.Len(t, pool.connections, 1)
			})
		})

		t.Run("close", func(t *testing.T) {
			pool := getPoolForTesting()
			empty := make(map[string]ResultSet)
			openConn1 := &connection{
				logHandler: logHandler,
				protocol:   nil,
				results:    empty,
				state:      established,
			}
			openConn2 := &connection{
				logHandler: logHandler,
				protocol:   nil,
				results:    empty,
				state:      established,
			}
			connections := []*connection{openConn1, openConn2}
			pool.connections = connections

			pool.close()
			assert.Equal(t, closed, openConn1.state)
			assert.Equal(t, closed, openConn2.state)
		})
	})
}
