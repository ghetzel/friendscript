package utils

import (
	"github.com/ghetzel/go-stockutil/stringutil"
)

type DefaultExecutor struct {
	from any
}

func NewDefaultExecutor(from any) *DefaultExecutor {
	return &DefaultExecutor{
		from: from,
	}
}

func (self *DefaultExecutor) SetInstance(from any) {
	if from != nil {
		self.from = from
	}
}

func (self *DefaultExecutor) FormatCommandName(name string) string {
	return stringutil.Camelize(name)
}

func (self *DefaultExecutor) ExecuteCommand(name string, arg any, objargs map[string]any) (any, error) {
	return CallCommandFunction(self.from, self.FormatCommandName(name), arg, objargs)
}
