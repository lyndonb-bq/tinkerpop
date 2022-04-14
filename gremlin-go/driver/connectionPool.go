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
	"sync"
)

type connectionPool interface {
	write(*request) (ResultSet, error)
	close()
}

const defaultNewConnectionThreshold = 4

// loadBalancingPool has two configurations: maximumConcurrentConnections/cap(connections) and newConnectionThreshold.
// maximumConcurrentConnections denotes the maximum amount of active connections at any given time.
// newConnectionThreshold specifies the minimum amount of concurrent active traversals on the least used connection
// which will trigger creation of a new connection if maximumConcurrentConnections has not been reached.
// loadBalancingPool will use the least-used connection, and as a part of the process, getLeastUsedConnection(), will
// remove any unusable connections from the pool and ensure that the returned connection is usable. If there are
// multiple active connections with no active traversals on them, one will be used and the others will be closed and
// removed from the pool.
type loadBalancingPool struct {
	url          string
	logHandler   *logHandler
	connSettings *connectionSettings

	newConnectionThreshold int
	connections            []*connection
	loadBalanceLock        sync.Mutex
}

func (pool *loadBalancingPool) close() {
	for _, connection := range pool.connections {
		err := connection.close()
		if err != nil {
			pool.logHandler.logf(Warning, errorClosingConnection, err.Error())
		}
	}
}

func (pool *loadBalancingPool) write(request *request) (ResultSet, error) {
	connection, err := pool.getLeastUsedConnection()
	if err != nil {
		return nil, err
	}
	return connection.write(request)
}

func (pool *loadBalancingPool) getLeastUsedConnection() (*connection, error) {
	pool.loadBalanceLock.Lock()
	defer pool.loadBalanceLock.Unlock()
	if len(pool.connections) == 0 {
		return pool.newConnection()
	} else {
		var leastUsed *connection = nil
		validIndex := 0
		for _, connection := range pool.connections {
			// Purge dead connections from pool
			if connection.state == established {
				// Close and purge connections from pool if there is more than one being unused
				if leastUsed != nil && (leastUsed.activeResults() == 0 && connection.activeResults() == 0) {
					// Close the connection asynchronously since it is a high-latency method
					go func() {
						pool.logHandler.log(Debug, closeUnusedPoolConnection)
						err := connection.close()
						if err != nil {
							pool.logHandler.logf(Warning, errorClosingConnection, err.Error())
						}
					}()

					continue
				}

				// Mark connection as valid to keep
				pool.connections[validIndex] = connection
				validIndex++

				// Set the least used connection
				if leastUsed == nil || connection.activeResults() < leastUsed.activeResults() {
					leastUsed = connection
				}
			} else {
				pool.logHandler.log(Debug, purgingDeadConnection)
			}
		}

		// Deallocate truncated dead connections to prevent memory leak
		for invalidIndex := validIndex; invalidIndex < len(pool.connections); invalidIndex++ {
			pool.connections[invalidIndex] = nil
		}
		pool.connections = pool.connections[:validIndex]

		// Create new connection if no valid connections were found in the pool or the least used connection exceeded
		// the concurrent usage threshold while the pool still has capacity for a new connection
		if leastUsed == nil ||
			(leastUsed.activeResults() >= pool.newConnectionThreshold && len(pool.connections) < cap(pool.connections)) {
			return pool.newConnection()
		} else {
			return leastUsed, nil
		}
	}
}

func (pool *loadBalancingPool) newConnection() (*connection, error) {
	connection, err := createConnection(pool.url, pool.logHandler, pool.connSettings)
	if err != nil {
		return nil, err
	}
	pool.connections = append(pool.connections, connection)
	return connection, nil
}

func newLoadBalancingPool(url string, logHandler *logHandler, connSettings *connectionSettings, newConnectionThreshold int,
	maximumConcurrentConnections int) (connectionPool, error) {
	pool := make([]*connection, 0, maximumConcurrentConnections)
	initialConnection, err := createConnection(url, logHandler, connSettings)
	if err != nil {
		return nil, err
	}
	pool = append(pool, initialConnection)
	return &loadBalancingPool{
		url:                    url,
		logHandler:             logHandler,
		connSettings:           connSettings,
		newConnectionThreshold: newConnectionThreshold,
		connections:            pool,
	}, nil
}
