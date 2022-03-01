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

// AnonymousTraversalSource struct used to generate anonymous traversals.
type AnonymousTraversalSource struct {
}

// WithRemote used to set the DriverRemoteConnection within the AnonymousTraversalSource.
func (ats *AnonymousTraversalSource) WithRemote(drc *DriverRemoteConnection) *GraphTraversalSource {
	return NewDefaultGraphTraversalSource().WithRemote(drc)
}

func Traversal_() *AnonymousTraversalSource {
	return &AnonymousTraversalSource{}
}

// AnonymousTraversal struct used for anonymous traversals
type AnonymousTraversal struct {
	graphTraversal func() *GraphTraversal
}

var T__ = &AnonymousTraversal{
	func() *GraphTraversal {
		return NewGraphTraversal(nil, nil, newBytecode(nil), nil)
	},
}

// T__ creates an empty GraphTraversal
func (anonymousTraversal *AnonymousTraversal) T__(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.Inject(args...)
}

// V adds the v step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) V(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().V(args...)
}

// AddE adds the addE step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) AddE(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().AddE(args...)
}

// AddV adds the addV step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) AddV(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().AddV(args...)
}

// Aggregate adds the aggregate step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Aggregate(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Aggregate(args...)
}

// And adds the and step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) And(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().And(args...)
}

// As adds the as step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) As(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().As(args...)
}

// Barrier adds the barrier step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Barrier(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Barrier(args...)
}

// Both adds the both step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Both(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Both(args...)
}

// BothE adds the bothE step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) BothE(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().BothE(args...)
}

// BothV adds the bothV step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) BothV(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().BothV(args...)
}

// Branch adds the branch step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Branch(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Branch(args...)
}

// By adds the by step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) By(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().By(args...)
}

// Cap adds the cap step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Cap(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Cap(args...)
}

// Choose adds the choose step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Choose(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Choose(args...)
}

// Coalesce adds the coalesce step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Coalesce(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Coalesce(args...)
}

// Coin adds the coint step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Coin(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Coin(args...)
}

// ConnectedComponent adds the connectedComponent step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) ConnectedComponent(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().ConnectedComponent(args...)
}

// Constant adds the constant step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Constant(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Constant(args...)
}

// Count adds the count step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Count(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Count(args...)
}

// CyclicPath adds the cyclicPath step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) CyclicPath(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().CyclicPath(args...)
}

// Dedup adds the dedup step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Dedup(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Dedup(args...)
}

// Drop adds the drop step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Drop(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Drop(args...)
}

// ElementMap adds the elementMap step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) ElementMap(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().ElementMap(args...)
}

// Emit adds the emit step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Emit(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Emit(args...)
}

// Filter adds the filter step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Filter(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Filter(args...)
}

// FlatMap adds the flatMap step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) FlatMap(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().FlatMap(args...)
}

// Fold adds the fold step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Fold(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Fold(args...)
}

// From adds the from step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) From(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().From(args...)
}

// Group adds the group step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Group(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Group(args...)
}

// GroupCount adds the groupCount step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) GroupCount(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().GroupCount(args...)
}

// Has adds the has step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Has(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Has(args...)
}

// HasId adds the hasId step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) HasId(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().HasId(args...)
}

// HasKey adds the hasKey step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) HasKey(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().HasKey(args...)
}

// HasLabel adds the hasLabel step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) HasLabel(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().HasLabel(args...)
}

// HasNot adds the hasNot step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) HasNot(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().HasNot(args...)
}

// HasValue adds the hasValue step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) HasValue(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().HasValue(args...)
}

// Id adds the id step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Id(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Id(args...)
}

// Identity adds the identity step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Identity(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Identity(args...)
}

// InE adds the inE step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) InE(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().InE(args...)
}

// InV adds the inV step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) InV(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().InV(args...)
}

// In adds the in step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) In(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().In(args...)
}

// Index adds the index step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Index(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Index(args...)
}

// Inject adds the inject step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Inject(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Inject(args...)
}

// Is adds the is step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Is(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Is(args...)
}

// Key adds the key step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Key(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Key(args...)
}

// Label adds the label step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Label(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Label(args...)
}

// Limit adds the limit step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Limit(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Limit(args...)
}

// Local adds the local step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Local(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Local(args...)
}

// Loops adds the loops step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Loops(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Loops(args...)
}

// Map adds the map step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Map(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Map(args...)
}

// Match adds the match step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Match(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Match(args...)
}

// Math adds the math step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Math(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Math(args...)
}

// Max adds the max step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Max(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Max(args...)
}

// Mean adds the mean step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Mean(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Mean(args...)
}

// Min adds the min step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Min(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Min(args...)
}

// None adds the none step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) None(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().None(args...)
}

// Not adds the not step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Not(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Not(args...)
}

// Option adds the option step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Option(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Option(args...)
}

// Optional adds the optional step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Optional(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Optional(args...)
}

// Or adds the or step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Or(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Or(args...)
}

// Order adds the order step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Order(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Order(args...)
}

// OtherV adds the otherV step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) OtherV(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().OtherV(args...)
}

// Out adds the out step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Out(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Out(args...)
}

// OutE adds the outE step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) OutE(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().OutE(args...)
}

// OutV adds the outV step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) OutV(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().OutV(args...)
}

// PageRank adds the pageRank step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) PageRank(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().PageRank(args...)
}

// Path adds the path step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Path(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Path(args...)
}

// PeerPressure adds the peerPressure step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) PeerPressure(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().PeerPressure(args...)
}

// Profile adds the profile step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Profile(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Profile(args...)
}

// Program adds the program step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Program(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Program(args...)
}

// Project adds the project step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Project(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Project(args...)
}

// Properties adds the properties step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Properties(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Properties(args...)
}

// Property adds the property step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Property(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Property(args...)
}

// PropertyMap adds the propertyMap step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) PropertyMap(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().PropertyMap(args...)
}

// Range adds the range step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Range(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Range(args...)
}

// Read adds the read step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Read(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Read(args...)
}

// Repeat adds the repeat step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Repeat(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Repeat(args...)
}

// Sack adds the sack step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Sack(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Sack(args...)
}

// Sample adds the sample step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Sample(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Sample(args...)
}

// Select adds the select step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Select(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Select(args...)
}

// ShortestPath adds the shortestPath step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) ShortestPath(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().ShortestPath(args...)
}

// SideEffect adds the sideEffect step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) SideEffect(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().SideEffect(args...)
}

// SimplePath adds the simplePath step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) SimplePath(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().SimplePath(args...)
}

// Skip adds the skip step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Skip(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Skip(args...)
}

// Store adds the store step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Store(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Store(args...)
}

// Subgraph adds the subgraph step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Subgraph(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Subgraph(args...)
}

// Sum adds the sum step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Sum(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Sum(args...)
}

// Tail adds the tail step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Tail(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Tail(args...)
}

// TimeLimit adds the timeLimit step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) TimeLimit(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().TimeLimit(args...)
}

// Times adds the times step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Times(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Times(args...)
}

// To adds the to step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) To(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().To(args...)
}

// ToE adds the toE step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) ToE(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().ToE(args...)
}

// ToV adds the toV step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) ToV(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().ToV(args...)
}

// Tree adds the tree step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Tree(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Tree(args...)
}

// Unfold adds the unfold step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Unfold(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Unfold(args...)
}

// Union adds the union step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Union(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Union(args...)
}

// Until adds the until step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Until(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Until(args...)
}

// Value adds the value step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Value(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Value(args...)
}

// ValueMap adds the valueMap step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) ValueMap(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().ValueMap(args...)
}

// Values adds the values step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Values(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Values(args...)
}

// Where adds the where step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Where(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Where(args...)
}

// With adds the with step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) With(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().With(args...)
}

// Write adds the write step to the GraphTraversal
func (anonymousTraversal *AnonymousTraversal) Write(args ...interface{}) *GraphTraversal {
	return anonymousTraversal.graphTraversal().Write(args...)
}
