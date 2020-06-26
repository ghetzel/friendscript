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
	env utils.Runtime
}

func New(env utils.Runtime) *Commands {
	cmd := &Commands{
		env: env,
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

type ReadArgs struct {
	// Whether to attempt to close the source (if possible) after reading.
	Autoclose bool `json:"autoclose" default:"true"`

	// The amount of data (in bytes) to read from the readable stream.
	Length int64 `json:"length" default:"-1"`
}

type ReadResponse struct {
	// The readable data.
	Data io.Reader `json:"data"`

	// The length of the data (in bytes).
	Length int64 `json:"length,omitempty"`
}

func (self *Commands) Read(fileOrReader interface{}, args *ReadArgs) (*ReadResponse, error) {
	if args == nil {
		args = new(ReadArgs)
	}

	defaults.SetDefaults(args)

	// get a readable stream representing the source path we were given
	if stream, err := self.env.Open(fileOrReader); err == nil {
		var buf = bytes.NewBuffer(nil)
		var response = new(ReadResponse)

		// if we're supposed to close the source, defer that now
		if args.Autoclose {
			defer stream.Close()
		}

		if args.Length >= 0 {
			// read the first N bytes
			if n, err := io.CopyN(buf, stream, args.Length); err == nil {
				response.Length = n
			} else {
				return nil, err
			}
		} else {
			// read ALL THE BYTES
			if n, err := io.Copy(buf, stream); err == nil {
				response.Length = n
			} else {
				return nil, err
			}
		}

		// whatever is in the buffer, that's what you get.
		response.Data = buf

		return response, nil
	} else {
		return nil, err
	}

}

type WriteArgs struct {
	// The data to write to the destination.
	Data io.Reader `json:"data"`

	// The data to write as a discrete value.
	Value interface{} `json:"value"`

	// Whether to attempt to close the destination (if possible) after reading/writing.
	Autoclose bool `json:"autoclose" default:"true"`
}

type WriteResponse struct {
	// The filesystem path that the data was written to.
	Path string `json:"path,omitempty"`

	// The size of the data (in bytes).
	Size int64 `json:"size,omitempty"`
}

// Write a value or a stream of data to a file at the given path.  The destination path can be a local
// filesystem path, a URI that uses a custom scheme registered outside of the application, or the string
// "temporary", which will write to a temporary file whose path will be returned in the response.
func (self *Commands) Write(destination interface{}, args *WriteArgs) (*WriteResponse, error) {
	if args == nil {
		args = new(WriteArgs)
	}

	defaults.SetDefaults(args)

	var response = new(WriteResponse)
	var writer io.Writer

	if destination != nil {
		if filename, ok := destination.(string); ok {
			if newPath, w, err := self.env.GetWriterForPath(filename); err == nil {
				writer = w
				response.Path = newPath
			} else {
				return nil, err
			}

			if writer == nil {
				if filename == `temporary` {
					if temp, err := ioutil.TempFile(``, ``); err == nil {
						writer = temp
						response.Path = temp.Name()
					} else {
						return nil, err
					}
				} else if file, err := os.Create(filename); err == nil {
					writer = file
					response.Path = filename
				} else {
					return nil, err
				}
			}
		} else if w, ok := destination.(io.Writer); ok {
			writer = w
		} else {
			return nil, fmt.Errorf("Unsupported destination %T; expected string or io.Writer", destination)
		}
	}

	if writer == nil {
		return response, fmt.Errorf("A destination must be specified")
	}

	if writer != nil {
		var err error

		if args.Data != nil {
			response.Size, err = io.Copy(writer, args.Data)
		} else if args.Value != nil {
			source := bytes.NewBufferString(fmt.Sprintf("%v", args.Value))
			response.Size, err = io.Copy(writer, source)
		} else {
			err = fmt.Errorf("Must specify source data or a discrete value to write")
		}

		// if whatever write (or write attempt) we just did succeeded...
		if err == nil {
			// if we're supposed to autoclose the destination, give that a shot now
			if args.Autoclose {
				if closer, ok := writer.(io.Closer); ok {
					return response, closer.Close()
				}
			}
		} else {
			return response, err
		}
	} else {
		return response, fmt.Errorf("Unable to write to destination")
	}

	return response, nil
}
