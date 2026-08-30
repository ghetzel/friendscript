// Contains miscellaneous utility commands.
package utils

import (
	"go.gary.cool/friendscript/utils"
)

type Commands struct {
	utils.Module
	env utils.Runtime
}

func New(env utils.Runtime) *Commands {
	cmd := &Commands{
		env: env,
	}

	cmd.Module = utils.NewDefaultExecutor(cmd)
	return cmd
}
