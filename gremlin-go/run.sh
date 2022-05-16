#!/bin/bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

# Parses current gremlin server version from project pom.xml file using perl regex
GREMLIN_SERVER_VERSION=$(grep tinkerpop -A2 pom.xml | grep -Po '(?<=<version>)([0-9]+\.?){3}(-SNAPSHOT)?(?=<)')
echo "$GREMLIN_SERVER_VERSION"

# Passes current gremlin server version into docker compose as environment variable
GREMLIN_SERVER="$GREMLIN_SERVER_VERSION"  docker-compose up --exit-code-from gremlin-go-integration-tests
docker-compose down
