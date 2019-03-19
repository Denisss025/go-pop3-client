package pop3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(ts *testing.T) {
	ts.Run("nested errors", func(t *testing.T) {
		err := newError("level 1", newError("level 2",
			fmt.Errorf("error")))
		assert.EqualError(t, err, "pop3: level 2: level 1: error")
	})
}
