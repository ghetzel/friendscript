package assert

import (
	"testing"

	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/testify/require"
)

func TestAssertExists(t *testing.T) {
	var cmd = New(nil)

	require.NoError(t, cmd.Exists(`hello`, nil))
	require.NoError(t, cmd.Exists(true, nil))
	require.NoError(t, cmd.Exists(false, nil))
	require.NoError(t, cmd.Exists(0, nil))
	require.NoError(t, cmd.Exists(float64(0), nil))

	require.True(t, log.ErrHasPrefix(cmd.Exists(``, nil), `Expected `))
	require.True(t, log.ErrHasPrefix(cmd.Exists(nil, nil), `Expected `))
	require.True(t, log.ErrHasPrefix(cmd.Exists(nil, &AssertArgs{}), `Expected `))

	require.True(t, cmd.Exists(nil, &AssertArgs{
		Message: `nope.`,
	}).Error() == `nope.`)
}
