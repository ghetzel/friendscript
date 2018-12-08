package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsErrorStatus(t *testing.T) {
	assert := require.New(t)

	for i := 0; i < 999; i++ {
		v := isErrorStatus(i, `402-999,0-400`)

		if i == 401 {
			assert.True(v)
		} else {
			assert.False(v)
		}
	}

	assert.True(isErrorStatus(199, `200`))
	assert.False(isErrorStatus(200, `200`))
	assert.True(isErrorStatus(201, `200`))
	assert.True(isErrorStatus(199, `200-299`))
	assert.False(isErrorStatus(200, `200-299`))
	assert.False(isErrorStatus(299, `200-299`))
	assert.True(isErrorStatus(300, `200-299`))

	assert.True(isErrorStatus(199, `200,204,404`))
	assert.False(isErrorStatus(200, `200,204,404`))
	assert.True(isErrorStatus(201, `200,204,404`))

	assert.True(isErrorStatus(203, `200,204,404`))
	assert.False(isErrorStatus(204, `200,204,404`))
	assert.True(isErrorStatus(205, `200,204,404`))

	assert.True(isErrorStatus(403, `200,204,404`))
	assert.False(isErrorStatus(404, `200,204,404`))
	assert.True(isErrorStatus(405, `200,204,404`))

}
