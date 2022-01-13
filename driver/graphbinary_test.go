package gremlingo

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestGraphBinaryV1(t *testing.T) {
	t.Run("test long", func(t *testing.T) {
		var x = 100
		var y int32 = 100
		var z int64 = 100
		var s = "serialize this!"
		var m = map[interface{}]interface{}{
			"marko": 666,
			"noone": "blah",
		}
		var u, _ = uuid.Parse("41d2e28a-20a4-4ab0-b379-d810dede3786")
		writer := graphBinaryWriter{}
		fmt.Println(writer.write(x))
		fmt.Println(writer.write(y))
		fmt.Println(writer.write(z))
		fmt.Println(writer.write(s))
		fmt.Println(writer.write(m))
		fmt.Println(writer.write(u))
	})
}
