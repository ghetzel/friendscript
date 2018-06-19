package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileOperations(t *testing.T) {
	assert := require.New(t)
	cmd := New(nil)

	temp, err := cmd.Temp(nil)
	assert.NoError(err)

	defer os.Remove(temp.Name())

	assert.NoError(cmd.Write(temp, &WriteArgs{
		Value: `Testing`,
	}))

	reopen, err := cmd.Open(temp.Name())
	assert.NoError(err)

	readback, err := cmd.Read(reopen)
	assert.NoError(err)

	assert.Equal(`Testing`, readback)
}
