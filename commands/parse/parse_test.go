package parse

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/ghetzel/friendscript/scripting"
	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/testify/require"
)

var testJsonData = `{
	"hello": {
		"there": true
	}
}`

var testYamlData = `{
---
hello:
	there: true
`

type testRuntime struct{}

func (self *testRuntime) Scope() *scripting.Scope {
	return nil
}

func (self *testRuntime) Run(scriptName string, options *utils.RunOptions) (interface{}, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (self *testRuntime) GetReaderForPath(path string) (io.ReadCloser, error) {
	return nil, nil
}

func (self *testRuntime) GetWriterForPath(path string) (string, io.Writer, error) {
	return ``, nil, nil
}

func (self *testRuntime) Open(fileOrReader interface{}) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(testJsonData)), nil
}

func (self *testRuntime) RegisterPathWriter(handler utils.PathWriterFunc) {}
func (self *testRuntime) RegisterPathReader(handler utils.PathReaderFunc) {}

func TestParseJson(t *testing.T) {
	var cmd = New(new(testRuntime))

	out, err := cmd.Json(bytes.NewBufferString(testJsonData))
	require.NoError(t, err)
	require.EqualValues(t, map[string]interface{}{
		`hello`: map[string]interface{}{
			`there`: true,
		},
	}, out)
}

func TestParseYaml(t *testing.T) {
	var cmd = New(new(testRuntime))

	out, err := cmd.Yaml(bytes.NewBufferString(testYamlData))
	require.NoError(t, err)
	require.EqualValues(t, map[interface{}]interface{}{
		`hello`: map[interface{}]interface{}{
			`there`: true,
		},
	}, out)
}
