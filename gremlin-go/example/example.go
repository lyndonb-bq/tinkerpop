package main

import (
	"fmt"
	"github.com/lyndonbauto/tinkerpop/gremlin-go/v3/driver"
)

func main() {
	// Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection("ws://localhost:8182/gremlin",
		func(settings *gremlingo.DriverRemoteConnectionSettings) {
			settings.TraversalSource = "gmodern"
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Cleanup
	defer driverRemoteConnection.Close()

	// Creating graph traversal
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

	// Perform traversal
	result, err := g.V().HasLabel("person").Order().By("age").Values("name").ToList()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range result {
		fmt.Println(r.GetString())
	}
}
