# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

from gremlin_python.process.anonymous_traversal import traversal
from gremlin_python.driver.driver_remote_connection import DriverRemoteConnection
import json


to_string = json.dumps


def main():

    # connect to a remote server that is compatible with the Gremlin Server protocol. for those who
    # downloaded and are using Gremlin Server directly be sure that it is running locally with:
    #
    # bin/gremlin-server.sh console
    #
    # which starts it in "console" mode with an empty in-memory TinkerGraph ready to go bound to a
    # variable named "g" as referenced in the following line.

    g = traversal().withRemote(DriverRemoteConnection('ws://localhost:8182/gremlin', 'g'))

    # add some data - be sure to use a terminating step like iterate() so that the traversal
    # "executes". iterate() does not return any data and is used to just generate side-effects
    # (i.e. write data to the database)
    foo = g.V().count().toList()
    bar = g.V().hasLabel("Test").count().toList()
    print(foo)
    print(bar)

if __name__ == "__main__":
    main()
