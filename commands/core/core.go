// This package provides default implementations of commands that can be used when creating modules.
package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	defaults "github.com/mcuadros/go-defaults"
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

// Outputs a line to the log.
func (self *Commands) Log(message interface{}) error {
	if typeutil.IsScalar(reflect.ValueOf(message)) {
		fmt.Printf("%v\n", message)
	} else if data, err := json.MarshalIndent(message, ``, `  `); err == nil {
		fmt.Printf(string(data) + "\n")
	} else {
		log.Errorf("Failed to log message: %v", err)
		return err
	}

	return nil
}

// Store a value in the current scope. Strings will be automatically converted
// into the appropriate data types (float, int, bool) if possible.
func (self *Commands) Put(value interface{}) (interface{}, error) {
	return value, nil
}

type EnvArgs struct {
	// The value to return if the environment variable does not exist, or
	// (optionally) is empty.
	Fallback interface{} `json:"fallback"`

	// Whether empty values should be ignored or not.
	Required bool `json:"required" default:"false"`

	// Whether automatic type detection should be performed or not.
	DetectType bool `json:"detect_type" default:"true"`

	// If specified, this string will be used to split matching values into a
	// list of values. This is useful for environment variables that contain
	// multiple values joined by a separator (e.g: the PATH variable.)
	Joiner string `json:"joiner"`
}

// Retrieves a system environment variable and returns the value of it, or a
// fallback value if the variable does not exist or (optionally) is empty.
//
// #### Examples
//
// ##### Get the value of the `USER` environment variable and store it
// ```
// env 'USER' -> $user
// ```
//
// ##### Require the `LANG`, `USER`, and `CI` environment variables; and fail they are not set.
// ```
// env 'LANG' { required: true }
// env 'USER' { required: true }
// env 'CI'   { required: true }
// ```
//
func (self *Commands) Env(name string, args *EnvArgs) (interface{}, error) {
	if args == nil {
		args = &EnvArgs{}
	}

	defaults.SetDefaults(args)

	if ev := os.Getenv(name); ev != `` {
		var rv interface{}

		if args.Joiner != `` {
			rv = strings.Split(ev, args.Joiner)
		}

		// perform type detection
		if args.DetectType {
			// for arrays, autotype each element
			if typeutil.IsArray(rv) {
				rv = sliceutil.Autotype(rv)
			} else {
				rv = stringutil.Autotype(ev)
			}
		} else {
			rv = ev
		}

		return rv, nil
	} else if args.Required {
		return nil, fmt.Errorf("Environment variable %q was not specified", name)
	} else {
		return nil, nil
	}
}

// Immediately exit the script in an error-like fashion with a specific message.
func (self *Commands) Fail(message string) error {
	if message == `` {
		message = `Unspecified error`
	}

	return errors.New(message)
}

type RunArgs struct {
	Data          interface{} `json:"data"`           // null
	Isolated      bool        `json:"isolated"`       // true
	PreserveState bool        `json:"preserve_state"` // true
	MergeScopes   bool        `json:"merge_scopes"`   // false
	ResultKey     string      `json:"result_key"`     // result
}

// [SKIP]
// Evaluates another Friendscript loaded from another file. The filename is the
// absolute path or basename of the file to search for in the FRIENDSCRIPT_PATH
// environment variable to load and evaluate. The FRIENDSCRIPT_PATH variable
// behaves like the the traditional *nix PATH variable, wherein multiple paths
// can be specified as a colon-separated (:) list. The current working directory
// will always be checked first.
//
// Returns: The value of the variable named by result_key at the end of the
// evaluated script's execution.
//
func (self *Commands) Run(filename string, args *RunArgs) (interface{}, error) {
	return nil, fmt.Errorf(`Not Implemented Yet`)
}
