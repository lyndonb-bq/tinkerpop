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

TINKERPOP_HOME=/opt/gremlin-server
cp /opt/test/scripts/* ${TINKERPOP_HOME}/scripts

IP=$(hostname -i)

INCLUDE_NEO4J=

function usage {
  echo -e "\nUsage: `basename $0` <version> [OPTIONS]" \
          "\nStart Gremlin Server instances that match the Maven integration test environment." \
          "\n\nOptions are:\n" \
          "\n\t<version> This value is optional and if unspecified will build the current version" \
          "\n\t-n, --neo4j              include Neo4j to make transactions testable" \
          "\n\t-h, --help               show this message" \
          "\n"
}

while [ ! -z "$1" ]; do
  case "$1" in
    -n | --neo4j ) INCLUDE_NEO4J=true; shift ;;
    -h | --help ) usage; exit 0 ;;
    *) usage 1>&2; exit 1 ;;
  esac
done

echo "#############################################################################"
echo IP is $IP
echo
echo Available Gremlin Server instances:
echo "ws://${IP}:45940/gremlin with anonymous access"
echo "wss://${IP}:45941/gremlin with basic authentication (stephen/password)"
echo
if [ ! -z "${INCLUDE_NEO4J}" ]; then
  echo Installing Neo4j to the environment: transactions are testable on port 45940
  echo
fi
echo "#############################################################################"

cp *.yaml ${TINKERPOP_HOME}/conf/

java -version

# dynamically installs Neo4j libraries so that we can test variants with transactions,
# but only only port 45940 is configured with the neo4j graph as the neo4j-empty.properties
# is statically pointing at a temp directory and that space can only be accessed by one
# graph at a time
if [ ! -z "${INCLUDE_NEO4J}" ]; then
  sed -i 's/graphs: {/graphs: {\n  tx: conf\/neo4j-empty.properties,/' ${TINKERPOP_HOME}/conf/gremlin-server-integration.yaml
  /opt/gremlin-server/bin/gremlin-server.sh install org.apache.tinkerpop neo4j-gremlin ${NEO4J_VERSION}
fi

/opt/gremlin-server/bin/gremlin-server.sh conf/gremlin-server-integration.yaml &

/opt/gremlin-server/bin/gremlin-server.sh conf/gremlin-server-integration-secure.yaml
