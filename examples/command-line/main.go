package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghetzel/friendscript"
	"github.com/ghetzel/friendscript/commands/core"
	"github.com/ghetzel/friendscript/utils"
)

type CoreCommands struct {
	*core.Commands
}

func NewCoreCommands(env utils.Scopeable) *CoreCommands {
	cmd := &CoreCommands{
		Commands: core.New(env),
	}

	cmd.SetInstance(cmd)

	return cmd
}

// Exit the program with the given status code.
//
// This will be available within Friendscript as the "exit" command.
//
func (self *CoreCommands) Exit(status int) error {
	os.Exit(status)
	return nil
}

// Return the list of files and subdirectories in the given directory path.
//
// This will be available within Friendscript as the "ls" command.
//
func (self *CoreCommands) Ls(path string) ([]string, error) {
	if entries, err := ioutil.ReadDir(path); err == nil {
		paths := make([]string, 0)

		for _, entry := range entries {
			paths = append(paths, entry.Name())
		}

		return paths, nil
	} else {
		return nil, err
	}
}

func main() {
	// create a new Friendscript scripting environment
	environment := friendscript.NewEnvironment()

	// add in our commands module, which extends the default "core" commands module
	// by adding two new commands: "exit" and "ls".
	environment.RegisterModule(``, NewCoreCommands(environment))

	if len(os.Args) > 1 {
		for _, scriptPath := range os.Args[1:] {
			if _, err := environment.EvaluateFile(scriptPath); err == nil {
				os.Exit(0)
			} else {
				fmt.Printf("script error: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Printf("usage: %v SCRIPT [SCRIPT ..]\n", os.Args[0])
		os.Exit(127)
	}
}
