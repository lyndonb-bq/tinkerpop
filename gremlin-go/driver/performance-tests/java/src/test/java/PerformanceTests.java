/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// CAUTION!!!
// This file is intended for generating Performance test results only.
// It is not to be used in a production setting.
// It is being made available as a reference to the test results for
// anyone that wants see the code used to obtain the performance metrics.

package org.apache.tinkerpop.gremlin.driver;

import lombok.AllArgsConstructor;
import org.apache.tinkerpop.gremlin.driver.remote.DriverRemoteConnection;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversal;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.GraphTraversalSource;
import org.apache.tinkerpop.gremlin.process.traversal.dsl.graph.__;
import org.junit.Test;
import java.time.Duration;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicReference;

import static org.apache.tinkerpop.gremlin.process.traversal.AnonymousTraversalSource.traversal;

class PerformanceTest {
    final int SAMPLE_SIZE = 101;
    final Cluster cluster = Cluster.build("localhost").port(8182).create();
    final Client client = cluster.connect();
    final Client g1client = client.alias("g");
    final DriverRemoteConnection connection = DriverRemoteConnection.using(client, "ggrateful");
    final GraphTraversalSource g = traversal().withRemote(connection);
    final int REPEATS = 50;

    public static GraphTraversal<?, ?> getTraversal(final String name, final GraphTraversalSource g) {
        final Map<String, GraphTraversal<?, ?>> traversals = new HashMap<>();
        traversals.put("getSongVertices", g.V().hasLabel("song"));
        traversals.put("getNumVertices", g.V().count());
        return traversals.get(name);

    }

    public static String[] getArgs(final int repeats) {
        final String[] args = new String[repeats];
        for (int i = 0; i < repeats; i++) {
            args[i] = Integer.toString(i);
        }
        return args;
    }

    public static GraphTraversal<?, ?> getProjectTraversal(final GraphTraversalSource g, final int repeats, final String[] args) {
        GraphTraversal<?, ?> traversal = g.V().project(args[0], Arrays.copyOfRange(args, 1, repeats));
        for (int i = 0; i < repeats; i++) {
            traversal = traversal.by(__.valueMap(true));
        }
        return traversal;
    }

    @Test
    public void sample() {
        System.out.println("~~~~~~~ RUNNING ONE ITEM PERFORMANCE TEST ~~~~~~~");
        executeGetOneItemPerformanceTest();
        System.out.println("~~~~~~~ RUNNING LIST PERFORMANCE TEST ~~~~~~~");
        executeGetItemListPerformanceTest();
        System.out.println("~~~~~~~ RUNNING TRANSFER PERFORMANCE TEST ~~~~~~~");
        executeTransferPerformanceTest();
        System.out.println("~~~~~~~ PERFORMANCE TESTS COMPLETE ~~~~~~~\"");
    }

    public timingData executeGetOneItemPerformanceTest() {
        final List<Duration> durations = new ArrayList<>();
        for (int i = 0; i < SAMPLE_SIZE; i++) {
            try {
                final String[] args = getArgs(REPEATS);
                final Instant starts = Instant.now();
                final Object temp = getProjectTraversal(g, REPEATS, args).next();
                final Instant ends = Instant.now();
                durations.add(Duration.between(starts, ends));
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        timingData data = getTimingDataFromDurationList(durations);
        System.out.println("Retrieve One Item (.next) performance test results:");
        System.out.println(data.toStringMillis("One Item"));
        return data;
    }

    public timingData executeGetItemListPerformanceTest() {
        final List<Duration> durations = new ArrayList<>();
        for (int i = 0; i < SAMPLE_SIZE; i++) {
            try {
                final String[] args = getArgs(REPEATS);
                final Instant starts = Instant.now();
                final List<?> list = getProjectTraversal(g, REPEATS, args).toList();
                for (Object item : list) {
                    final Object temp = item;
                }
                final Instant ends = Instant.now();
                durations.add(Duration.between(starts, ends));
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        timingData data = getTimingDataFromDurationList(durations);
        System.out.println("Retrieve Item List (.toList) performance test results:");
        System.out.println(data.toStringMillis("List"));
        return data;
    }

    public timingData executeTransferPerformanceTest() {
        final List<Duration> durations = new ArrayList<>();
        for (int i = 0; i < SAMPLE_SIZE; i++) {
            try {
                final String[] args = getArgs(REPEATS);
                final GraphTraversal<?, ?> traversal = getProjectTraversal(g, REPEATS, args);
                final Instant starts = Instant.now();
                while (traversal.hasNext()) {
                    Object obj = traversal.next();
                }
                final Instant ends = Instant.now();
                durations.add(Duration.between(starts, ends));
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        timingData data = getTimingDataFromDurationList(durations);
        System.out.println("Retrieve transfer traversal result (.next --> variable) performance test results:");
        System.out.println(data.toStringMillis("Transfer"));
        return data;
    }

    private timingData getTimingDataFromDurationList(final List<Duration> durations) {;
        Collections.sort(durations);
        durations.remove(0);
        return new timingData(
                durations.stream().reduce(Duration.ZERO, Duration::plus).dividedBy(SAMPLE_SIZE),
                durations.get(durations.size() / 2),
                durations.get((int)(durations.size() * 0.90)),
                durations.get((int)(durations.size() * 0.95)),
                durations.get(0),
                durations.get(durations.size() - 1)
        );
    }
}

@AllArgsConstructor
class timingData {
    final Duration AVG;
    final Duration MEDIAN;
    final Duration PERCENTILE_90;
    final Duration PERCENTILE_95;
    final Duration MIN;
    final Duration MAX;

    @Override
    public String toString() {
        return toStringMillis("unknown");
    }

    public String toStringMillis(final String testType) {
        return "Test Type: " + testType + "\n"+
                "\tAVG=" + AVG.toMillis() + "ms \n" +
                "\tMEDIAN=" + MEDIAN.toMillis() + "ms \n" +
                "\tPERCENTILE_90=" + PERCENTILE_90.toMillis() + "ms \n" +
                "\tPERCENTILE_95=" + PERCENTILE_95.toMillis() + "ms \n" +
                "\tMIN=" + MIN.toMillis() + "ms \n" +
                "\tMAX=" + MAX.toMillis() + "ms \n";
    }

    public String toStringNanos(final String testType) {
        return "Test Type " + testType + "\n"+
                "\tAVG=" + AVG.toNanos() + "ms \n" +
                "\tMEDIAN=" + MEDIAN.toNanos() + "ms \n" +
                "\tPERCENTILE_90=" + PERCENTILE_90.toNanos() + "ms \n" +
                "\tPERCENTILE_95=" + PERCENTILE_95.toNanos() + "ms \n" +
                "\tMIN=" + MIN.toNanos() + "ms \n" +
                "\tMAX=" + MAX.toNanos() + "ms \n";
    }
}