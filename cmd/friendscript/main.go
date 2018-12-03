package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/friendscript"
	"github.com/ghetzel/go-stockutil/log"
)

func main() {
	app := cli.NewApp()
	app.Name = `friendscript`
	app.Usage = `Friendscript eval`
	app.Version = friendscript.Version
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `info`,
			EnvVar: `LOGLEVEL`,
		},
		cli.BoolFlag{
			Name:  `print-vars, P`,
			Usage: `Print the final state of all variables upon script completion.`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))
		return nil
	}

	app.Action = func(c *cli.Context) {
		// evaluate Friendscript
		script := friendscript.NewEnvironment()

		var input io.Reader

		if c.NArg() > 0 {
			filename := c.Args().First()

			switch filename {
			case `-`:
				input = os.Stdin
			default:
				if file, err := os.Open(c.Args().First()); err == nil {
					log.Debugf("Friendscript being read from file %s", file.Name())
					input = file
				} else {
					log.Fatalf("file error: %v", err)
					return
				}
			}
		} else {
			log.Fatalf("Must specify a Friendscript filename to execute.")
			return
		}

		if scope, err := script.EvaluateReader(input); err == nil {
			if c.Bool(`print-vars`) {
				fmt.Println(scope)
			}
		} else {
			log.Fatalf("runtime error: %v", err)
		}
	}

	app.Run(os.Args)
}
