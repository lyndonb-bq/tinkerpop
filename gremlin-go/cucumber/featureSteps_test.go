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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/lyndonb-bq/tinkerpop/gremlin-go/driver"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// TODO proper error handling
type tinkerPopGraph struct {
	*TinkerPopWorld
}

var parsers map[*regexp.Regexp]func(string, string) interface{}

func init() {
	parsers = map[*regexp.Regexp]func(string, string) interface{}{
		regexp.MustCompile("d\\[(.*)]\\.[ilfdm]"): toNumeric,
		regexp.MustCompile("v\\[(.+)]"):           toVertex,
		regexp.MustCompile("v\\[(.+)]\\.id"):      toVertexId,
		regexp.MustCompile("e\\[(.+)]"):           toEdge,
		regexp.MustCompile("v\\[(.+)]\\.sid"):     toVertexIdString,
		regexp.MustCompile("e\\[(.+)]\\.id"):      toEdgeId,
		regexp.MustCompile("e\\[(.+)]\\.sid"):     toEdgeIdString,
		regexp.MustCompile("p\\[(.+)]"):           toPath,
		regexp.MustCompile("l\\[(.*)]"):           toList,
		regexp.MustCompile("s\\[(.*)]"):           toSet,
		regexp.MustCompile("m\\[(.+)]"):           toMap,
		regexp.MustCompile("c\\[(.+)]"):           toLambda,
		regexp.MustCompile("t\\[(.+)]"):           toT,
		regexp.MustCompile("null"):                func(string, string) interface{} { return nil },
	}
}

func parseValue(value string, graphName string) interface{} {
	var extractedValue string
	var parser func(string, string) interface{}
	for key, element := range parsers {
		var match = key.FindAllStringSubmatch(value, -1)
		if len(match) > 0 {
			parser = element
			extractedValue = match[0][1]
			break
		}
	}
	if parser == nil {
		return value
	} else {
		return parser(extractedValue, graphName)
	}
}

// parse numeric
func toNumeric(stringVal, graphName string) interface{} {
	val, err := strconv.ParseFloat(stringVal, 64)
	if err != nil {
		return nil
	}
	return val
}

// parse vertex
func toVertex(name, graphName string) interface{} {
	return tg.getDataGraphFromMap(graphName).vertices[name]
}

// parse vertex id
func toVertexId(name, graphName string) interface{} {
	if tg.getDataGraphFromMap(graphName).vertices[name] == nil {
		return nil
	}
	return tg.getDataGraphFromMap(graphName).vertices[name].Id
}

// parse vertex id as string
func toVertexIdString(name, graphName string) interface{} {
	if tg.getDataGraphFromMap(graphName).vertices[name] == nil {
		return nil
	}
	return fmt.Sprint(tg.getDataGraphFromMap(graphName).vertices[name].Id)
}

// parse edge
func toEdge(name, graphName string) interface{} {
	return tg.getDataGraphFromMap(graphName).edges[name]
}

// parse edge id
func toEdgeId(name, graphName string) interface{} {
	if tg.getDataGraphFromMap(graphName).edges[name] == nil {
		return nil
	}
	return tg.getDataGraphFromMap(graphName).edges[name].Id
}

// parse edge id as string
func toEdgeIdString(name, graphName string) interface{} {
	if tg.getDataGraphFromMap(graphName).edges[name] == nil {
		return nil
	}
	return fmt.Sprint(tg.getDataGraphFromMap(graphName).edges[name])
}

// TODO add with updated path implementation
//parse path
func toPath(name, graphName string) interface{} {
	return nil
}

// parse list
func toList(stringList, graphName string) interface{} {
	listVal := make([]interface{}, 0)
	if len(stringList) == 0 {
		return listVal
	}

	for _, str := range strings.Split(stringList, ",") {
		listVal = append(listVal, parseValue(str, graphName))
	}
	return listVal
}

// TODO add with custom set implementation
// parse set
func toSet(name, graphName string) interface{} {
	return nil
}

// parse json as a map
func toMap(name, graphName string) interface{} {
	var jsonValue interface{}
	err := json.Unmarshal([]byte(name), jsonValue)
	if err != nil {
		return nil
	}
	return parseMapValue(jsonValue, graphName)
}

func parseMapValue(value interface{}, graphName string) interface{} {
	if value == nil {
		return nil
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return parseValue(value.(string), graphName)
	case reflect.Float64:
		return toNumeric(value.(string), graphName)
	case reflect.Array, reflect.Slice:
		return parseMapValue(value, graphName)
	case reflect.Map:
		valMap := make(map[interface{}]interface{})
		v := reflect.ValueOf(value)
		keys := v.MapKeys()
		for _, k := range keys {
			convKey := k.Convert(v.Type().Key())
			val := v.MapIndex(convKey)
			valMap[parseMapValue(convKey, graphName)] = parseMapValue(val, graphName)
		}
		return valMap
	default:
		// not supported types
		return nil
	}
}

// TODO add with lambda implementation
// parse lambda
func toLambda(name, graphName string) interface{} {
	return nil
}

// TODO add with T(label) implementation
// parse instance of T enum
func toT(name, graphName string) interface{} {
	return nil
}

func (tg *tinkerPopGraph) anUnsupportedTest() error {
	return nil
}

// TODO add with .Next() implementation
func (tg *tinkerPopGraph) iteratedNext() error {
	return godog.ErrPending
}

func (tg *tinkerPopGraph) iteratedToList() error {
	if tg.traversal == nil {
		return errors.New("nil traversal, feature need to be implemented in go")
	}
	results, err := tg.traversal.ToList()
	if err != nil {
		return err
	}
	tg.result = results
	fmt.Println("RESULT LENGTH", len(results))
	count, err := tg.traversal.Count().ToList()
	for _, c := range count {
		fmt.Println("COUNT() LENGTH", c.GetInterface())
	}
	return nil
}

func (tg *tinkerPopGraph) nothingShouldHappenBecause(arg1 *godog.DocString) error {
	return nil
}

// choose the graph
func (tg *tinkerPopGraph) chooseGraph(graphName string) error {
	tg.graphName = graphName
	data := tg.graphDataMap[graphName]
	tg.g = gremlingo.Traversal_().WithRemote(data.connection)
	if graphName == "empty" {
		err := tg.cleanEmptyDataGraph()
		if err != nil {
			return err
		}
	}
	return nil
}

func (tg *tinkerPopGraph) theGraphInitializerOf(arg1 *godog.DocString) error {
	traversal, err := GetTraversal(tg.scenario.Name, tg.g, tg.parameters)
	if err != nil {
		return err
	}
	_, future, err := traversal.Iterate()
	if err != nil {
		return err
	}
	<-future
	// We may have modified the so-called `empty` graph
	if tg.graphName == "empty" {
		tg.reloadEmptyData()
	}
	return nil
}

func (tg *tinkerPopGraph) theGraphShouldReturnForCountOf(expectedCount int, traversalText string) error {
	traversal, err := GetTraversal(tg.scenario.Name, tg.g, tg.parameters)
	if err != nil {
		return err
	}
	results, err := traversal.ToList()
	if err != nil {
		return err
	}
	if len(results) != expectedCount {
		return errors.New("graph did not return the correct count")
	}
	return nil
}
func (tg *tinkerPopGraph) theResultShouldBe(characterizedAs string, table *godog.Table) error {
	fmt.Println("===GOT TO RESULT SHOULD BE===", characterizedAs)
	ordered := characterizedAs == "ordered"
	switch characterizedAs {
	case "empty":
		if len(tg.result) != 0 {
			return errors.New("actual result is not empty as expected")
		}
		return nil
	case "ordered", "unordered", "of":
		var expectedResult []interface{}
		for idx, row := range table.Rows {
			if idx == 0 {
				// skip the header line
				continue
			}
			val := parseValue(row.Cells[0].Value, tg.graphName)
			v, ok := val.(gremlingo.Path)
			if ok {
				// clear the labels since we don't define them in .feature files
				v.Labels = [][]string{}
				val = v
			}
			expectedResult = append(expectedResult, val)
		}
		var actualResult []interface{}
		for _, res := range tg.result {
			actualResult = append(actualResult, res.GetInterface())
		}
		fmt.Println("EXPECTED RESULTS", expectedResult)
		fmt.Println("ACTUAL RESULTS", actualResult)
		if characterizedAs != "of" && len(actualResult) != len(expectedResult) {
			err := fmt.Sprintf("actual result length %d does not equal to expected result length %d.", len(actualResult), len(expectedResult))
			return errors.New(err)
		}
		if ordered {
			for idx, res := range actualResult {
				if !reflect.DeepEqual(expectedResult[idx], res) {
					return errors.New("actual result is not ordered")
				}
			}
		} else {
			for _, res := range actualResult {
				if !contains(expectedResult, res) {
					return errors.New("actual result does not match expected result")
				}
			}
		}
		break
	default:
		return errors.New("scenario not supported")
	}
	return nil
}

func contains(list []interface{}, item interface{}) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

func (tg *tinkerPopGraph) theResultShouldHaveACountOf(expectedCount int) error {
	actualCount := len(tg.result)
	if len(tg.result) != expectedCount {
		err := fmt.Sprintf("result should return %d for count, but returned %d.", expectedCount, actualCount)
		return errors.New(err)
	}
	return nil
}

func (tg *tinkerPopGraph) theTraversalOf(arg1 *godog.DocString) error {
	traversal, err := GetTraversal(tg.scenario.Name, tg.g, tg.parameters)
	if err != nil {
		return err
	}
	tg.traversal = traversal
	return nil
}

func (tg *tinkerPopGraph) usingTheParameterDefined(name string, params string) error {
	if tg.graphName == "empty" {
		tg.reloadEmptyData()
	}
	tg.parameters[name] = parseValue(strings.Replace(params, "\\\"", "\"", -1), tg.graphName)
	return godog.ErrPending
}

func (tg *tinkerPopGraph) usingTheParameterOfP(paramName, pVal, stringVal string) error {
	// TODO after implementing P class
	return godog.ErrPending
}

// safe?
var tg = &tinkerPopGraph{
	NewTinkerPopWorld(),
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		tg.loadAllDataGraph()
	})
	ctx.AfterSuite(func() {
		tg.closeAllDataGraphConnection()
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tg.scenario = sc
		//tg.recreateAllDataGraphConnection()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		//tg.closeAllDataGraphConnection()
		return ctx, nil
	})

	ctx.Step(`^an unsupported test$`, tg.anUnsupportedTest)
	ctx.Step(`^iterated next$`, tg.iteratedNext)
	ctx.Step(`^iterated to list$`, tg.iteratedToList)
	ctx.Step(`^nothing should happen because$`, tg.nothingShouldHappenBecause)
	ctx.Step(`^the (.+) graph$`, tg.chooseGraph)
	ctx.Step(`^the graph initializer of$`, tg.theGraphInitializerOf)
	ctx.Step(`^the graph should return (\d+) for count of "(.+)"$`, tg.theGraphShouldReturnForCountOf)
	ctx.Step(`^the result should be (\w+)$`, tg.theResultShouldBe)
	ctx.Step(`^the result should have a count of (\d+)$`, tg.theResultShouldHaveACountOf)
	ctx.Step(`^the traversal of$`, tg.theTraversalOf)
	ctx.Step(`^using the parameter (.+) defined as "(.+)"$`, tg.usingTheParameterDefined)

	// TODO depends on implementation of class P (comparators)
	ctx.Step(`^using the parameter (.+) of P\.(.+)\("(.+)"\)$`, tg.usingTheParameterOfP)
}
