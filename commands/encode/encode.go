// Commands used to encode data into various formats.
package encode

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-defaults"
	"github.com/ghetzel/go-stockutil/typeutil"
	"gopkg.in/yaml.v2"
)

type Commands struct {
	utils.Module
	env utils.Runtime
}

func New(env utils.Runtime) *Commands {
	var cmd = &Commands{
		env: env,
	}

	cmd.Module = utils.NewDefaultExecutor(cmd)
	return cmd
}

type JsonArgs struct {
	// Indent output with the given number of spaces.
	Indent int `json:"indent" default:"2"`
}

// Encode the given data into a JSON document.
func (self *Commands) Json(value any, args *JsonArgs) (out string, merr error) {
	if args == nil {
		args = new(JsonArgs)
	}

	defaults.SetDefaults(args)
	var outb []byte

	if args.Indent == 0 {
		outb, merr = json.Marshal(value)
	} else {
		outb, merr = json.MarshalIndent(value, ``, strings.Repeat(` `, args.Indent))
	}

	out = string(outb)
	return
}

// Encode the given data into a YAML document.
func (self *Commands) Yaml(value any) (string, error) {
	var outb, err = yaml.Marshal(value)
	return string(outb), err
}

// Encode the given data into base64
func (self *Commands) Base64(value any) (string, error) {
	return base64.StdEncoding.EncodeToString(
		[]byte(typeutil.String(value)),
	), nil
}
