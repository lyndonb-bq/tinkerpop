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

package results

const DEFAULT_CAPACITY = 1000

// AN-968 Finish ResultSet implementation.

// ResultSet interface to define the functions of a ResultSet.
type ResultSet interface {
	SetAggregateTo(val string)
	GetAggregateTo() string
	SetStatusAttributes(statusAttributes map[interface{}]interface{})
	GetStatusAttributes() map[interface{}]interface{}
	GetRequestId() int
	IsEmpty() bool
	Close()
	Channel() chan *Result
	AddResult(result *Result)
	One() *Result
	All() []*Result
}

// ChannelResultSet Channel based implementation of ResultSet.
type ChannelResultSet struct {
	channel          chan *Result
	aggregateTo      string
	statusAttributes map[interface{}]interface{}
	closed           bool
}

func (channelResultSet *ChannelResultSet) IsEmpty() bool {
	return channelResultSet.closed && len(channelResultSet.channel) == 0
}

func (channelResultSet *ChannelResultSet) Close() {
	close(channelResultSet.channel)
	channelResultSet.closed = true
}

func (channelResultSet *ChannelResultSet) SetAggregateTo(val string) {
	channelResultSet.aggregateTo = val
}

func (channelResultSet *ChannelResultSet) GetAggregateTo() string {
	return channelResultSet.aggregateTo
}

func (channelResultSet *ChannelResultSet) SetStatusAttributes(val map[interface{}]interface{}) {
	channelResultSet.statusAttributes = val
}

func (channelResultSet *ChannelResultSet) GetStatusAttributes() map[interface{}]interface{} {
	return channelResultSet.statusAttributes
}

func (channelResultSet *ChannelResultSet) GetRequestId() int {
	return -1
}

func (channelResultSet *ChannelResultSet) Channel() chan *Result {
	return channelResultSet.channel
}

func (channelResultSet *ChannelResultSet) One() *Result {
	return <-channelResultSet.channel
}

func (channelResultSet *ChannelResultSet) All() []*Result {
	var results []*Result
	for result := range channelResultSet.channel {
		results = append(results, result)
	}
	return results
}

func (channelResultSet *ChannelResultSet) AddResult(result *Result) {
	channelResultSet.channel <- result
}

func NewChannelResultSetCapacity(channelSize int) ResultSet {
	return &ChannelResultSet{make(chan *Result, channelSize), "", nil, false}
}

func NewChannelResultSet() ResultSet {
	return NewChannelResultSetCapacity(DEFAULT_CAPACITY)
}
