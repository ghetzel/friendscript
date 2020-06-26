// Commands used to load and parse various formats of serialized data (JSON, YAML, etc.)
//
package parse

import (
	"encoding/json"

	"github.com/ghetzel/friendscript/utils"
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

// Parses the given file as a JSON document and returns the resulting value.
func (self *Commands) Json(fileOrReader interface{}) (interface{}, error) {
	if rc, err := self.env.Open(fileOrReader); err == nil {
		defer rc.Close()
		var out interface{}

		if err := json.NewDecoder(rc).Decode(&out); err == nil {
			return out, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// Parses the given file as a YAML document and returns the resulting value.
func (self *Commands) Yaml(fileOrReader interface{}) (interface{}, error) {
	if rc, err := self.env.Open(fileOrReader); err == nil {
		defer rc.Close()
		var out interface{}

		if err := yaml.NewDecoder(rc).Decode(&out); err == nil {
			return out, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
