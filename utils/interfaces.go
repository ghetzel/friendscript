package utils

import (
	"io"

	"github.com/ghetzel/friendscript/scripting"
)

type PathWriterFunc = func(string) (string, io.Writer, error)
type PathReaderFunc = func(string) (io.ReadCloser, error)

type Runtime interface {
	Scope() *scripting.Scope
	Run(scriptName string, options *RunOptions) (any, error)
	GetReaderForPath(path string) (io.ReadCloser, error)
	GetWriterForPath(path string) (string, io.Writer, error)
	RegisterPathWriter(handler PathWriterFunc)
	RegisterPathReader(handler PathReaderFunc)
	Open(fileOrReader any) (io.ReadCloser, error)
}

type Module interface {
	ExecuteCommand(name string, arg any, objargs map[string]any) (any, error)
	FormatCommandName(string) string
	SetInstance(any)
}
