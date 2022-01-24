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

const defaultCapacity = 1000

// ResultSet interface to define the functions of a ResultSet.
type ResultSet interface {
	setAggregateTo(val string)
	GetAggregateTo() string
	setStatusAttributes(statusAttributes map[interface{}]interface{})
	GetStatusAttributes() map[interface{}]interface{}
	GetRequestId() int
	IsEmpty() bool
	Close()
	Channel() chan *Result
	addResult(result *Result)
	one() *Result
	All() []*Result
}

// channelResultSet Channel based implementation of ResultSet.
type channelResultSet struct {
	channel          chan *Result
	aggregateTo      string
	statusAttributes map[interface{}]interface{}
	closed           bool
}

func (channelResultSet *channelResultSet) IsEmpty() bool {
	return channelResultSet.closed && len(channelResultSet.channel) == 0
}

func (channelResultSet *channelResultSet) Close() {
	close(channelResultSet.channel)
	channelResultSet.closed = true
}

func (channelResultSet *channelResultSet) setAggregateTo(val string) {
	channelResultSet.aggregateTo = val
}

func (channelResultSet *channelResultSet) GetAggregateTo() string {
	return channelResultSet.aggregateTo
}

func (channelResultSet *channelResultSet) setStatusAttributes(val map[interface{}]interface{}) {
	channelResultSet.statusAttributes = val
}

func (channelResultSet *channelResultSet) GetStatusAttributes() map[interface{}]interface{} {
	return channelResultSet.statusAttributes
}

func (channelResultSet *channelResultSet) GetRequestId() int {
	return -1
}

func (channelResultSet *channelResultSet) Channel() chan *Result {
	return channelResultSet.channel
}

func (channelResultSet *channelResultSet) one() *Result {
	return <-channelResultSet.channel
}

func (channelResultSet *channelResultSet) All() []*Result {
	var results []*Result
	for result := range channelResultSet.channel {
		results = append(results, result)
	}
	return results
}

func (channelResultSet *channelResultSet) addResult(result *Result) {
	channelResultSet.channel <- result
}

func newChannelResultSetCapacity(channelSize int) ResultSet {
	return &channelResultSet{make(chan *Result, channelSize), "", nil, false}
}

func newChannelResultSet() ResultSet {
	return newChannelResultSetCapacity(defaultCapacity)
}
