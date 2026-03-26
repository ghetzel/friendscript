package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/friendscript"
	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

var closables []io.Closer

func main() {
	var app = cli.NewApp()
	app.Name = `friendscript`
	app.Version = friendscript.Version
	app.EnableBashCompletion = true
	app.ArgsUsage = `[COMMAND_STRING | FILE]`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `info`,
			EnvVar: `LOGLEVEL`,
		},
		cli.BoolFlag{
			Name:  `execute, c`,
			Usage: `Execute commands provided as command line arguments.`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))
		return nil
	}

	app.Action = func(c *cli.Context) {
		// go handleSignals(func() {
		// 	for _, closer := range closables {
		// 		closer.Close()
		// 	}
		// })

		// evaluate Friendscript / run the REPL
		var script = friendscript.NewEnvironment(nil)

		// pre-populate initial variables
		for _, pair := range c.StringSlice(`var`) {
			var k, v = stringutil.SplitPair(pair, `=`)
			script.Set(k, typeutil.Auto(v))
		}

		var input io.ReadCloser = os.Stdin

		if c.Bool(`execute`) {
			input = io.NopCloser(
				bytes.NewBufferString(strings.Join(c.Args(), ` `)),
			)
		} else if c.NArg() == 0 {
			if state, err := term.GetState(int(os.Stdin.Fd())); err == nil {
				defer term.Restore(int(os.Stdin.Fd()), state)

				var scope, err = script.REPL()

				if err == nil {
					fmt.Println(scope)
				} else {
					log.FatalIf(err)
				}
			} else {
				log.FatalIf(err)
			}

			return
		} else if scriptpath := c.Args().First(); fileutil.Exists(scriptpath) {
			if file, err := os.Open(scriptpath); err == nil {
				log.Debugf("Friendscript being read from file %s", file.Name())
				input = file
				closables = append(closables, file)
			} else {
				log.Fatal(errors.Wrap(err, "file error"))
			}
		}

		if scope, err := script.EvaluateReader(input); err == nil {
			if prints := c.StringSlice(`print-var`); len(prints) > 0 {
				var out = make(map[string]any)

				for _, pv := range prints {
					if vv := scope.Get(pv); vv != nil {
						out[pv] = vv
					}
				}

				if data, err := json.MarshalIndent(out, ``, `  `); err == nil {
					fmt.Println(string(data))
				} else {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}

func handleSignals(handler func()) {
	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for range signalChan {
		handler()
		break
	}

	os.Exit(0)
}
