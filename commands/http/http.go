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
	"strings"
	"time"

	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-stockutil/httputil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
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

	// A map of cookie key=value pairs to include in the request.
	Cookies map[string]interface{} `json:"cookies"`

	// The amount of time to wait for the request to complete.
	Timeout time.Duration `json:"timeout" default:"30s"`

	// The body of the request. This is processed according to what is specified in RequestType.
	Body interface{} `json:"body"`

	// The type of data in Body, specifying how it should be encoded.  Valid values are "raw", "form", and "json"
	RequestType string `json:"request_type,omitempty" default:"json"`

	// Specify how the response body should be decoded.  Can be "raw", or a MIME type that overrides the Content-Type response header.
	ResponseType string `json:"response_type,omitempty"`

	// Whether to disable TLS peer verification.
	DisableVerifySSL bool `json:"disable_verify_ssl"`

	// The path to the root TLS CA bundle to use for verifying peer certificates.
	CertificateBundle string `json:"ca_bundle"`

	// A comma-separated list of numbers (e.g.: 200) or inclusive number ranges (e.g. 200-399) specifying HTTP statuses that are
	// expected and non-erroneous.
	Statuses string `json:"statuses" default:"200-299"`

	// Whether to continue execution if an error status is encountered.
	ContinueOnError bool `json:"continue_on_error"`
}

func (self *RequestArgs) Merge(other *RequestArgs) *RequestArgs {
	out := &RequestArgs{
		Headers:           self.Headers,
		Params:            self.Params,
		Timeout:           self.Timeout,
		Body:              self.Body,
		RequestType:       self.RequestType,
		ResponseType:      self.ResponseType,
		DisableVerifySSL:  self.DisableVerifySSL,
		Statuses:          self.Statuses,
		ContinueOnError:   self.ContinueOnError,
		CertificateBundle: self.CertificateBundle,
	}

	if other != nil {
		out.Headers, _ = maputil.Merge(out.Headers, other.Headers)
		out.Params, _ = maputil.Merge(out.Params, other.Params)
		out.Cookies, _ = maputil.Merge(out.Cookies, other.Cookies)

		if v := other.CertificateBundle; v != `` {
			out.CertificateBundle = v
		}

		if v := other.RequestType; v != `` {
			out.RequestType = v
		}

		if v := other.ResponseType; v != `` {
			out.ResponseType = v
		}

		if v := other.Statuses; v != `` {
			out.Statuses = v
		}

		if other.DisableVerifySSL {
			out.DisableVerifySSL = true
		}

		if other.ContinueOnError {
			out.ContinueOnError = true
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

	// A textual description of the HTTP response code.
	StatusText string `json:"status_text"`

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

	// If the response status is considered an error, and errors aren't fatal, this will be true.
	Error bool `json:"error"`
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
	// this is the bit that takes any defaults set via http::defaults and overlays the per-request values
	var reqargs = self.defaults.Merge(args)

	client := &http.Client{
		Timeout: reqargs.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: reqargs.DisableVerifySSL,
			},
		},
	}

	// specify CA bundle (if provided)
	if ca := reqargs.CertificateBundle; ca != `` {
		log.Debugf("friendscript/http: Using override CA bundle at %v", ca)

		if err := httputil.SetRootCABundle(client, ca); err != nil {
			return nil, err
		}
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

			log.Debugf("friendscript/http: -> %v %v", req.Method, req.URL)

			if body != nil {
				log.Debugf("friendscript/http: -> encoded body as %v (%v)", reqargs.RequestType, contentType)
			}

			for k, vs := range req.Header {
				log.Debugf("friendscript/http: -> [H] %v: %v", k, strings.Join(vs, `,`))
			}

			// populate cookies
			for k, v := range args.Cookies {
				req.AddCookie(&http.Cookie{
					Name:  k,
					Value: typeutil.String(v),
				})

				log.Debugf("friendscript/http: -> [C] %v: %v", k, v)
			}

			// perform the request
			if response, err := client.Do(req); err == nil {
				// build the response
				res := &HttpResponse{
					Status:     response.StatusCode,
					StatusText: response.Status,
					Headers:    make(map[string]interface{}),
					Took:       int64(time.Since(start).Nanoseconds() / 1e6),
				}

				log.Debugf("friendscript/http: <- HTTP %v (took %vms)", response.Status, res.Took)

				// add (autotyped) headers
				for k, vs := range response.Header {
					log.Debugf("friendscript/http: <- [H] %v: %v", k, strings.Join(vs, `,`))

					if len(vs) == 1 {
						res.Headers[k] = typeutil.Auto(vs[0])
					} else {
						res.Headers[k] = sliceutil.Autotype(vs)
					}
				}

				if response.Body != nil {
					defer response.Body.Close()
				}

				if isErrorStatus(response.StatusCode, reqargs.Statuses) {
					if reqargs.ContinueOnError {
						res.Error = true
					} else {
						log.Debugf("friendscript/http: <- Request error: %v", response.Status)
						return nil, fmt.Errorf("HTTP %v", response.Status)
					}
				}

				// decode (i.e.: decompress) response
				if response.ContentLength < 0 || response.ContentLength > 0 {
					if decoded, err := httputil.DecodeResponse(response); err == nil {
						// read and parse response body
						if data, err := ioutil.ReadAll(decoded); err == nil {
							res.Length = int64(len(data))
							res.Body = string(data)

							if res.Length > 0 {
								log.Debugf("friendscript/http: <- decoding body as %v", reqargs.ResponseType)

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
							}
						} else {
							log.Debugf("friendscript/http: <- Read response failed: %v", err)
							return nil, err
						}
					} else {
						log.Debugf("friendscript/http: <- Decode response failed: %v", err)
						return nil, err
					}
				}

				return res, nil
			} else {
				log.Debugf("friendscript/http: <- Request failed: %v", err)
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

func isErrorStatus(status int, allowed string) bool {
	ranges := strings.Split(allowed, `,`)

	for _, rng := range ranges {
		lowS, highS := stringutil.SplitPair(strings.TrimSpace(rng), `-`)
		low := int(typeutil.Int(lowS))
		high := int(typeutil.Int(highS))

		if high == 0 {
			if status == low {
				return false
			}
		} else {
			if status >= low && status <= high {
				return false
			}
		}
	}

	return true
}
