package pop3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListParser(ts *testing.T) {
	ts.Run("Normal string", func(t *testing.T) {
		index, size, err := parseListResponse("15 100020")
		assert.NoError(t, err)
		assert.Equal(t, 15, index, "test index")
		assert.Equal(t, int64(100020), size, "test size")
	})

	ts.Run("Bad string", func(t *testing.T) {
		_, _, err := parseListResponse("TEST STRING")
		assert.EqualError(t, err,
			`parse "TEST STRING": expected integer`)
	})

	ts.Run("Not enough ints", func(t *testing.T) {
		_, _, err := parseListResponse("1")
		assert.EqualError(t, err, `parse "1": EOF`)
	})
}
