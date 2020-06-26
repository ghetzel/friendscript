// This package provides default implementations of commands that can be used when creating modules.
package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/timeutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/kyokomi/emoji"
	defaults "github.com/mcuadros/go-defaults"
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

// Outputs a line to the log.
func (self *Commands) Log(message interface{}) error {
	if message == nil {
		return nil
	} else if typeutil.IsScalar(reflect.ValueOf(message)) {
		emoji.Printf("%v\n", message)
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
		return args.Fallback, nil
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

	// Specifies a key in the scope of the evaluate script that will be used as the result value of this command.
	ResultKey string `json:"result_key"` // result

	// Provides a set of initial variables to the script.
	Data map[string]interface{} `json:"data"` // null

	// If true, the scope of the running script will not be able to modify data in the parent scope.
	Isolated bool `json:"isolated" default:"true"`
}

// Evaluates another Friendscript loaded from another file. The filename is the
// absolute path or basename of the file to search for in the FRIENDSCRIPT_PATH
// environment variable to load and evaluate. The FRIENDSCRIPT_PATH variable
// behaves like the the traditional *nix PATH variable, wherein multiple paths
// can be specified as a colon-separated (:) list. The directory of the calling
// script (if available) will always be checked first.
//
// Returns: The value of the variable named by result_key at the end of the
// evaluated script's execution.
//
func (self *Commands) Run(filename string, args *RunArgs) (interface{}, error) {
	if self.env == nil {
		return nil, fmt.Errorf("no environment found")
	}

	if args == nil {
		args = &RunArgs{}
	}

	defaults.SetDefaults(args)

	var basePath string

	if ctx := self.env.Scope().EvalContext(); ctx != nil {
		basePath = filepath.Dir(ctx.Filename)
	}

	return self.env.Run(filename, &utils.RunOptions{
		Isolated:  false,
		ResultKey: args.ResultKey,
		Data:      args.Data,
		BasePath:  basePath,
	})
}

// Pauses execution of the current script for the given duration.
func (self *Commands) Wait(delay interface{}) error {
	var duration time.Duration

	if delayD, ok := delay.(time.Duration); ok {
		duration = delayD
	} else if delayMs, err := stringutil.ConvertToInteger(delay); err == nil {
		duration = time.Duration(delayMs) * time.Millisecond
	} else if delayParsed, err := timeutil.ParseDuration(fmt.Sprintf("%v", delay)); err == nil {
		duration = delayParsed
	} else {
		return fmt.Errorf("invalid duration: %v", err)
	}

	log.Infof("Waiting for %v", duration)
	time.Sleep(duration)
	return nil
}
