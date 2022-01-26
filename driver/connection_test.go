package gremlingo

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestConnection(t *testing.T) {

	t.Run("Test connect", func(t *testing.T) {
		connection := newConnection("localhost", 8181, Gorilla, newLogHandler(&defaultLogger{}, Info, language.English), nil, nil, nil)
		err := connection.connect()
		assert.Nil(t, err)
	})

	t.Run("Test write", func(t *testing.T) {
		connection := newConnection("localhost", 8181, Gorilla, newLogHandler(&defaultLogger{}, Info, language.English), nil, nil, nil)
		err := connection.connect()
		assert.Nil(t, err)
		resultSet, err := connection.write("g.V().count()")
		assert.Nil(t, err)
		assert.NotNil(t, resultSet)
		results := resultSet.All()
		for result := range results {
			assert.Equal(t, "0", result)
		}
	})

	//t.Run("Sandbox", func(t *testing.T) {
	//	client := NewClient("localhost", 8182)
	//
	//	response, err := client.Submit("1 + 1")
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println(response)
	//})
}
