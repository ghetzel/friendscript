package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileOperations(t *testing.T) {
	assert := require.New(t)

	temp, err := ioutil.TempFile(``, ``)
	assert.NoError(err)
	defer os.Remove(temp.Name())

	cmd := New(nil)
}
