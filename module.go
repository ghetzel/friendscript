package friendscript

import (
	"github.com/ghetzel/friendscript/utils"
)

type Module = utils.Module

func CreateModule(from interface{}) Module {
	return utils.NewDefaultExecutor(from)
}
