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

import "github.com/google/uuid"

type ResponseStatus struct {
	Code       uint16                      `json:"code"`
	Message    string                      `json:"message"`
	Attributes map[interface{}]interface{} `json:"attributes"`
}

type ResponseResult struct {
	Meta map[interface{}]interface{} `json:"meta"`
	Data interface{}                 `json:"data"`
}

type Response struct {
	RequestID      uuid.UUID      `json:"requestId"`
	ResponseStatus ResponseStatus `json:"responseStatus"`
	ResponseResult ResponseResult `json:"responseResult"`
}
