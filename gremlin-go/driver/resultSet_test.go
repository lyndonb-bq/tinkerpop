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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannelResultSet(t *testing.T) {
	const mockID = "mockID"

	t.Run("Test ResultSet test getter/setters.", func(t *testing.T) {
		r := newChannelResultSet(mockID)
		testStatusAttribute := map[string]interface{}{
			"1": 1234,
			"2": "foo",
		}
		testAggregateTo := "test2"
		r.setStatusAttributes(testStatusAttribute)
		assert.Equal(t, r.GetStatusAttributes(), testStatusAttribute)
		r.setAggregateTo(testAggregateTo)
		assert.Equal(t, r.GetAggregateTo(), testAggregateTo)
	})

	t.Run("Test ResultSet close.", func(t *testing.T) {
		channelResultSet := newChannelResultSet(mockID)
		assert.NotPanics(t, func() { channelResultSet.Close() })
	})

	t.Run("Test ResultSet one.", func(t *testing.T) {
		channelResultSet := newChannelResultSet(mockID)
		AddResults(&channelResultSet, 10)
		idx := 0
		for i := 0; i < 10; i++ {
			result := channelResultSet.one()
			assert.Equal(t, result.GetString(), fmt.Sprintf("%v", idx))
			idx++
		}
		go closeAfterTime(500, &channelResultSet)
		assert.Nil(t, channelResultSet.one())
	})

	t.Run("Test ResultSet one Paused.", func(t *testing.T) {
		channelResultSet := newChannelResultSet(mockID)
		go AddResultsPause(&channelResultSet, 10, 500)
		idx := 0
		for i := 0; i < 10; i++ {
			result := channelResultSet.one()
			assert.Equal(t, result.GetString(), fmt.Sprintf("%v", idx))
			idx++
		}
		go closeAfterTime(500, &channelResultSet)
		assert.Nil(t, channelResultSet.one())
	})

	t.Run("Test ResultSet one close.", func(t *testing.T) {
		channelResultSet := newChannelResultSet(mockID)
		channelResultSet.Close()
	})

	t.Run("Test ResultSet All.", func(t *testing.T) {
		channelResultSet := newChannelResultSet(mockID)
		AddResults(&channelResultSet, 10)
		go closeAfterTime(500, &channelResultSet)
		results := channelResultSet.All()
		for idx, result := range results {
			assert.Equal(t, (*result).GetString(), fmt.Sprintf("%v", idx))
		}
	})

	t.Run("Test ResultSet All close before.", func(t *testing.T) {
		channelResultSet := newChannelResultSet(mockID)
		AddResults(&channelResultSet, 10)
		channelResultSet.Close()
		results := channelResultSet.All()
		assert.Equal(t, len(results), 10)
		for idx, result := range results {
			assert.Equal(t, (*result).GetString(), fmt.Sprintf("%v", idx))
		}
	})
}

func AddResultsPause(resultSet *ResultSet, count int, ticks time.Duration) {
	rs := *resultSet
	for i := 0; i < count/2; i++ {
		rs.addResult(&Result{i})
	}
	time.Sleep(ticks * time.Millisecond)
	for i := count / 2; i < count; i++ {
		rs.addResult(&Result{i})
	}
}

func AddResults(resultSet *ResultSet, count int) {
	rs := *resultSet
	for i := 0; i < count; i++ {
		rs.addResult(&Result{i})
	}
}

func closeAfterTime(ticks time.Duration, resultSet *ResultSet) {
	time.Sleep(ticks * time.Millisecond)
	(*resultSet).Close()
}