// Commands for reading and writing files.
package file

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghetzel/friendscript/utils"
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

type TempArgs struct {
	// A string to prefix temporary filenames with
	Prefix string `json:"prefix" default:"friendscript-"`
}

func (self *Commands) Temp(args *TempArgs) (*os.File, error) {
	if args == nil {
		args = &TempArgs{}
	}

	defaults.SetDefaults(args)

	return ioutil.TempFile(``, args.Prefix)
}

func (self *Commands) Open(filename string) (*os.File, error) {
	return os.Open(filename)
}

func (self *Commands) Create(filename string) (*os.File, error) {
	return os.Create(filename)
}

func (self *Commands) Close(file *os.File) error {
	return file.Close()
}

type WriteArgs struct {
	// The data to write as a stream.
	Data io.Reader `json:"data"`

	// The data to write as a discrete value.
	Value interface{} `json:"value"`

	// Whether to attempt to close the destination (if possible) after reading/writing.
	Autoclose bool `json:"autoclose" default:"true"`
}

func (self *Commands) Write(destination interface{}, args *WriteArgs) error {
	var writer io.Writer

	if args == nil {
		args = &WriteArgs{}
	}

	defaults.SetDefaults(args)

	if f, ok := destination.(io.Writer); ok {
		writer = f
	} else if filename, ok := destination.(string); ok {
		if file, err := os.Open(filename); err == nil {
			writer = file
			destination = file
		} else {
			return err
		}
	} else {
		return fmt.Errorf("Must specify a filename string or writable stream as a destination")
	}

	if writer != nil {
		var err error

		if args.Data != nil {
			_, err = io.Copy(writer, args.Data)
		} else if args.Value != nil {
			source := bytes.NewBufferString(fmt.Sprintf("%v", args.Value))
			_, err = io.Copy(writer, source)
		} else {
			err = fmt.Errorf("Must specify source data or a discrete value to write")
		}

		// if whatever write (or write attempt) we just did succeeded...
		if err == nil {
			// if we're supposed to autoclose the destination, give that a shot now
			if args.Autoclose {
				if closer, ok := destination.(io.Closer); ok {
					return closer.Close()
				}
			}
		} else {
			return err
		}
	} else {
		return fmt.Errorf("Unable to write to destination")
	}

	return nil
}

func (self *Commands) Read(source interface{}) (string, error) {
	var reader io.Reader

	if f, ok := source.(io.Reader); ok {
		reader = f

	} else if filename, ok := source.(string); ok {
		if file, err := os.Open(filename); err == nil {
			reader = file
			source = file
		} else {
			return ``, err
		}
	} else {
		return ``, fmt.Errorf("Must specify a filename string or readable stream as a source")
	}

	// autoclose the source if it's closable
	if closer, ok := source.(io.Closer); ok {
		defer closer.Close()
	}

	if reader != nil {
		data, err := ioutil.ReadAll(reader)

		if err == nil {
			return string(data), nil
		} else {
			return ``, err
		}
	} else {
		return ``, fmt.Errorf("Unable to read from source")
	}
}
