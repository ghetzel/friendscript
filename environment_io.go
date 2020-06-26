package friendscript

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-stockutil/typeutil"
)

// Registers a new handler function that will be used for turning paths into writable streams.
func (self *Environment) RegisterPathWriter(handler utils.PathWriterFunc) {
	self.pathWriters = append([]utils.PathWriterFunc{
		handler,
	}, self.pathWriters...)
}

// Registers a new handler function that will be used for turning paths into readable streams.
func (self *Environment) RegisterPathReader(handler utils.PathReaderFunc) {
	self.pathReaders = append([]utils.PathReaderFunc{
		handler,
	}, self.pathReaders...)
}

// Takes a path string and consults all registered PathWriterFuncs.  The first one to claim it can handle
// the path will be responsible for returning a possibly-rewritten path string and an io.Writer that will
// accept the data being written.
func (self *Environment) GetWriterForPath(path string) (string, io.Writer, error) {
	for _, handler := range self.pathWriters {
		if p, w, err := handler(path); err == nil {
			// non-nil io.Writer + nil error = a handled request
			if w != nil {
				return p, w, nil
			}
		} else {
			return ``, nil, err
		}
	}

	return ``, nil, nil
}

// Takes a path string and consults all registered PathReaderFuncs.  The first one to claim it can handle
// the path will be responsible for returning an io.ReadCloser that represents the stream of data being
// sought.
func (self *Environment) GetReaderForPath(path string) (io.ReadCloser, error) {
	for _, handler := range self.pathReaders {
		if r, err := handler(path); err == nil {
			// non-nil RC + nil error = a handled request
			if r != nil {
				return r, nil
			}
		} else {
			return nil, err
		}
	}

	return os.Open(path)
}

// Open a readable destination file for reading.  If fileOrReader is a string, it will be treated
// as a path and will be sent to GetReaderForPath().  If it is an io.Reader, it will be returned
// without reading from it.
func (self *Environment) Open(fileOrReader interface{}) (io.ReadCloser, error) {
	var rc io.ReadCloser

	if b, ok := fileOrReader.([]byte); ok {
		rc = ioutil.NopCloser(bytes.NewBuffer(b))
	} else if i, ok := fileOrReader.([]interface{}); ok {
		rc = ioutil.NopCloser(bytes.NewBuffer(
			typeutil.Bytes(i),
		))
	} else if filename, ok := fileOrReader.(string); ok {
		if file, err := self.GetReaderForPath(filename); err == nil {
			rc = file
		} else {
			return nil, err
		}
	} else if r, ok := fileOrReader.(io.Reader); ok {
		rc = ioutil.NopCloser(r)
	} else {
		return nil, fmt.Errorf("argument must be a string or stream, got: %T", fileOrReader)
	}

	if rc == nil {
		return nil, fmt.Errorf("no readable data available")
	}

	return rc, nil
}
