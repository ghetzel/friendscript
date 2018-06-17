package utils

import (
	"github.com/ghetzel/go-stockutil/stringutil"
)

type DefaultExecutor struct {
	from interface{}
}

func NewDefaultExecutor(from interface{}) *DefaultExecutor {
	return &DefaultExecutor{
		from: from,
	}
}

func (self *DefaultExecutor) SetInstance(from interface{}) {
	if from != nil {
		self.from = from
	}
}

func (self *DefaultExecutor) FormatCommandName(name string) string {
	return stringutil.Camelize(name)
}

func (self *DefaultExecutor) ExecuteCommand(name string, arg interface{}, objargs map[string]interface{}) (interface{}, error) {
	return CallCommandFunction(self.from, self.FormatCommandName(name), arg, objargs)
}
