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
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuthentication(t *testing.T) {

	t.Run("Test BasicAuthInfo.", func(t *testing.T) {
		username := "Lyndon"
		password := "Bauto"
		header := BasicAuthInfo("Lyndon", "Bauto")
		assert.NotNil(t, header.getHeader())

		// fmt.Sprintf("%v, <value>) is used below because it strips off the http.Header type comparison issue.
		// http.Header is technically a map[string][]string underneath.
		usernamePasswordEncoded := map[string][]string{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))}}
		assert.Equal(t, fmt.Sprintf("%v", usernamePasswordEncoded), fmt.Sprintf("%v", header.getHeader()))
	})

	t.Run("Test getHeader.", func(t *testing.T) {
		header := &AuthInfo{}
		assert.Nil(t, header.getHeader())
		header = nil
		assert.Nil(t, header.getHeader())
		httpHeader := http.Header{}
		header = &AuthInfo{Header: httpHeader}
		assert.Equal(t, httpHeader, header.getHeader())
	})
}
