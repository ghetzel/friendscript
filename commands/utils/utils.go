// Contains miscellaneous utility commands.
package utils

import (
	"github.com/ghetzel/friendscript/utils"
)

type Commands struct {
	utils.Module
	scopeable utils.Scopeable
}

func New(scopeable utils.Scopeable) *Commands {
	cmd := &Commands{
		scopeable: scopeable,
	}

	cmd.Module = utils.NewDefaultExecutor(cmd)
	return cmd
}
