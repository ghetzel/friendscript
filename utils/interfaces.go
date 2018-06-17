package utils

import (
	"github.com/ghetzel/friendscript/scripting"
)

type Scopeable interface {
	Scope() *scripting.Scope
}

type Module interface {
	ExecuteCommand(name string, arg interface{}, objargs map[string]interface{}) (interface{}, error)
	FormatCommandName(string) string
	SetInstance(interface{})
}
