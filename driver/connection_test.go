package gremlingo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

const runIntegration = false

func TestConnection(t *testing.T) {
	t.Run("Test connect", func(t *testing.T) {
		if runIntegration {
			connection := connection{"localhost", 8181, Gorilla, newLogHandler(&defaultLogger{}, Info, language.English), nil, nil, nil}
			err := connection.connect()
			assert.Nil(t, err)
		}
	})

	t.Run("Test write", func(t *testing.T) {
		if runIntegration {
			connection := connection{"localhost", 8181, Gorilla, newLogHandler(&defaultLogger{}, Info, language.English), nil, nil, nil}
			err := connection.connect()
			assert.Nil(t, err)
			resultSet, err := connection.write("g.V().count()")
			assert.Nil(t, err)
			assert.NotNil(t, resultSet)
			result := resultSet.one()
			assert.NotNil(t, result)
			assert.Equal(t, "[0]", result.AsString())
		}
	})
}
