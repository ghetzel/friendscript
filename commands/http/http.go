// Commands for interacting with HTTP resources
package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-stockutil/httputil"
	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	defaults "github.com/mcuadros/go-defaults"
)

type Commands struct {
	utils.Module
	scopeable utils.Scopeable
	defaults  RequestArgs
}

type RequestArgs struct {
	// The headers to send with the request.
	Headers map[string]interface{} `json:"headers"`

	// Query string parameters to add to the request.
	Params map[string]interface{} `json:"params"`

	// The amount of time to wait for the request to complete.
	Timeout time.Duration `json:"timeout" default:"30s"`

	// The body of the request. This is processed according to what is specified in RequestType.
	Body interface{} `json:"body"`

	// The type of data in Body, specifying how it should be encoded.  Valid values are "raw", "form", and "json"
	RequestType string `json:"type" default:"raw"`

	// Specify how the response body should be decoded.  Can be "raw", or a MIME type that overrides the Content-Type response header.
	ResponseType string `json:"response_type"`

	// Whether to disable TLS peer verification.
	DisableVerifySSL bool `json:"disable_verify_ssl"`
}

func (self *RequestArgs) Merge(other *RequestArgs) RequestArgs {
	out := *self

	if other != nil {
		out.Headers, _ = maputil.Merge(out, other.Headers)
		out.Params, _ = maputil.Merge(out, other.Params)

		if v := other.RequestType; v != `` {
			out.RequestType = v
		}

		if v := other.ResponseType; v != `` {
			out.ResponseType = v
		}

		if other.DisableVerifySSL {
			out.DisableVerifySSL = true
		}

		if v := other.Timeout; v > 0 {
			out.Timeout = v
		}

		if v := other.Body; v != nil {
			out.Body = v
		}
	}

	return out
}

type HttpResponse struct {
	// The numeric HTTP status code of the response.
	Status int `json:"status"`

	// The time (in millisecond) that the request took to complete.
	Took int64 `json:"took"`

	// Response headers sent back from the server.
	Headers map[string]interface{} `json:"headers"`

	// The MIME type of the response body (if any).
	ContentType string `json:"type"`

	// The length of the response body in bytes.
	Length int64 `json:"length"`

	// The decoded response body (if any).
	Body interface{} `json:"body"`
}

func New(scopeable utils.Scopeable) *Commands {
	reqargs := &RequestArgs{}
	defaults.SetDefaults(reqargs)

	cmd := &Commands{
		scopeable: scopeable,
		defaults:  *reqargs,
	}

	cmd.Module = utils.NewDefaultExecutor(cmd)
	return cmd
}

// Set default options that apply to all subsequent HTTP requests.
func (self *Commands) Defaults(args *RequestArgs) error {
	if args == nil {
		args = &RequestArgs{}
	}

	defaults.SetDefaults(args)
	self.defaults = *args
	return nil
}

// Perform an HTTP GET request.
func (self *Commands) Get(url string, args *RequestArgs) (*HttpResponse, error) {
	return self.request(`GET`, url, args)
}

// Perform an HTTP POST request.
func (self *Commands) Post(url string, args *RequestArgs) (*HttpResponse, error) {
	return self.request(`POST`, url, args)
}

// Perform an HTTP PUT request.
func (self *Commands) Put(url string, args *RequestArgs) (*HttpResponse, error) {
	return self.request(`PUT`, url, args)
}

// Perform an HTTP DELETE request.
func (self *Commands) Delete(url string, args *RequestArgs) (*HttpResponse, error) {
	return self.request(`DELETE`, url, args)
}

// Perform an HTTP OPTIONS request.
func (self *Commands) Options(url string, args *RequestArgs) (*HttpResponse, error) {
	return self.request(`OPTIONS`, url, args)
}

// Perform an HTTP HEAD request.
func (self *Commands) Head(url string, args *RequestArgs) (*HttpResponse, error) {
	return self.request(`HEAD`, url, args)
}

func (self *Commands) request(method string, url string, args *RequestArgs) (*HttpResponse, error) {
	reqargs := self.defaults.Merge(args)

	client := &http.Client{
		Timeout: reqargs.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: reqargs.DisableVerifySSL,
			},
		},
	}

	// encode the body (if any) in preparation for sending in the request
	if body, contentType, err := encodeBody(reqargs.RequestType, reqargs.Body); err == nil {
		// get a new request
		if req, err := http.NewRequest(method, url, body); err == nil {
			// set query string parameters
			if len(reqargs.Params) > 0 {
				for k, v := range reqargs.Params {
					httputil.SetQ(req.URL, k, v)
				}
			}

			// set the body
			if body != nil {
				req.Body = ioutil.NopCloser(body)
			}

			// get headers in place
			for k, v := range self.defaults.Headers {
				req.Header.Set(k, typeutil.String(v))
			}

			// set content type detected during encoding
			if contentType != `` {
				req.Header.Set(`Content-Type`, contentType)
			}

			// set header overrides
			if args != nil {
				for k, v := range args.Headers {
					req.Header.Set(k, typeutil.String(v))
				}
			}

			start := time.Now()

			// perform the request
			if response, err := client.Do(req); err == nil {
				// build the response
				res := &HttpResponse{
					Status:  response.StatusCode,
					Headers: make(map[string]interface{}),
					Took:    int64(time.Since(start).Nanoseconds() / 1e6),
				}

				// add (autotyped) headers
				for k, vs := range response.Header {
					if len(vs) == 1 {
						res.Headers[k] = typeutil.Auto(vs[0])
					} else {
						res.Headers[k] = sliceutil.Autotype(vs)
					}
				}

				// decode (i.e.: decompress) response
				if response.ContentLength > 0 {
					if response.Body != nil {
						defer response.Body.Close()
					}

					if decoded, err := httputil.DecodeResponse(response); err == nil {
						// read and parse response body
						if data, err := ioutil.ReadAll(decoded); err == nil {
							res.Length = int64(len(data))
							res.Body = data

							switch reqargs.ResponseType {
							case `raw`:
								break
							default:
								// automatically decode response
								rt := response.Header.Get(`Content-Type`)

								if reqargs.ResponseType != `` {
									rt = reqargs.ResponseType
								}

								switch rt {
								case `application/json`:
									if err := json.Unmarshal(data, &res.Body); err != nil {
										return nil, err
									}
								}
							}
						} else {
							return nil, err
						}

					} else {
						return nil, err
					}
				}

				return res, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func encodeBody(enctype string, body interface{}) (io.Reader, string, error) {
	var reader io.Reader
	var contentType string = `application/octet-stream`

	if body != nil {
		switch enctype {
		case ``, `raw`:
			if r, ok := body.(io.Reader); ok {
				reader = r
			} else if b, ok := body.([]byte); ok {
				reader = bytes.NewBuffer(b)
			} else {
				contentType = `text/plain`
				reader = bytes.NewBufferString(typeutil.String(body))
			}

		case `form`:
			contentType = `application/x-www-form-urlencoded`

			if typeutil.IsMap(body) {
				values := make(url.Values)

				for key, value := range maputil.M(body).MapNative() {
					values.Set(key, typeutil.String(value))
				}

				reader = bytes.NewBufferString(values.Encode())

			} else if r, ok := body.(io.Reader); ok {
				contentType = `multipart/form-data`
				reader = r

			} else if b, ok := body.([]byte); ok {
				contentType = `multipart/form-data`
				reader = bytes.NewBuffer(b)

			} else {
				reader = bytes.NewBufferString(typeutil.String(body))
			}

		case `json`:
			contentType = `application/json`

			if typeutil.IsMap(body) {
				if data, err := json.Marshal(maputil.M(body).MapNative()); err == nil {
					reader = bytes.NewBuffer(data)
				} else {
					return nil, ``, err
				}
			} else {
				reader = bytes.NewBufferString(typeutil.String(body))
			}

		default:
			return nil, ``, fmt.Errorf("Unknown encoding type %q", enctype)
		}
	}

	return reader, contentType, nil
}
