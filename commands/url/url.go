// Commands for processing and working with URLs.
package url

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-stockutil/executil"
	"github.com/ghetzel/go-stockutil/log"
	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

// if unspecified, default to HTTPS.  This is probably in violation of some RFC,
// but in 2020 I'd rather default to encryption over none.  If the scheme is important
// to the user, they should specify it.  This is also exposed as a package-level variable
// that implementers may change directly or via environment variable
var DefaultUrlScheme = executil.Env(`FRIENDSCRIPT_URL_DEFAULT_SCHEME`, `https`)

type Commands struct {
	utils.Module
	env utils.Runtime
}

type URL struct {
	Scheme   string                 `json:"scheme"`
	Host     string                 `json:"host"`
	Domain   string                 `json:"domain"`
	Port     int                    `json:"port"`
	Path     string                 `json:"path"`
	RawQuery string                 `json:"rawquery"`
	Fragment string                 `json:"fragment"`
	Query    map[string]interface{} `json:"query"`
	Full     string                 `json:"full"`
	url      *url.URL
}

func (self *URL) String() string {
	var u = new(url.URL)

	u.Scheme = self.Scheme
	u.Fragment = self.Fragment
	u.Path = self.Path

	if self.Port == 80 && self.Scheme == `http` {
		u.Host = self.Domain
	} else if self.Port == 443 && self.Scheme == `https` {
		u.Host = self.Domain
	} else {
		u.Host = self.Host
	}

	if len(self.Query) > 0 {
		var vals = make(url.Values)

		for k, v := range self.Query {
			vals.Set(k, typeutil.String(v))
		}

		u.RawQuery = vals.Encode()
	}

	return u.String()
}

func (self *URL) sync() error {
	if self.url != nil {
		// net/url treats bare domains as URL paths, so we want to detect that case and treat them as hosts
		if self.url.Host == `` && !strings.HasPrefix(self.url.Path, `/`) {
			self.url.Host = self.url.Path
			self.url.Path = `/`
		}

		self.Scheme = self.url.Scheme
		self.Host = self.url.Host
		self.Path = self.url.Path
		self.RawQuery = self.url.RawQuery
		self.Fragment = self.url.Fragment
	}

	// set a default scheme
	if self.Scheme == `` {
		self.Scheme = DefaultUrlScheme
	}

	// normalize scheme
	self.Scheme = strings.ToLower(self.Scheme)

	// work out the domain and (if applicable) port number
	if self.Domain == `` {
		if self.Host != `` {
			if h, p, err := net.SplitHostPort(self.Host); err == nil {
				self.Domain = h

				if self.Port == 0 {
					self.Port = int(typeutil.Int(p))
				}
			} else if log.ErrContains(err, `missing port in address`) {
				self.Domain = self.Host
			} else {
				return err
			}
		}
	}

	// if a port is not explicitly specified, try to work it out from the scheme
	if self.Port == 0 {
		switch self.Scheme {
		case `http`:
			self.Port = 80
		case `https`:
			self.Port = 443
		}
	}

	// set a default path
	if self.Path == `` {
		self.Path = `/`
	}

	// parse and populate querystring values into a map
	if qmap, err := urlValuesToMap(self.RawQuery); err == nil {
		self.Query = qmap
	} else {
		return err
	}

	self.Full = self.String()
	return nil
}

func New(env utils.Runtime) *Commands {
	var cmd = &Commands{
		env: env,
	}

	cmd.Module = utils.NewDefaultExecutor(cmd)
	return cmd
}

// Parse the given URL string or structure, and return a structured representation of the various parts of a URL.
func (self *Commands) Parse(u interface{}) (*URL, error) {
	var parsed = &URL{
		Query: make(map[string]interface{}),
	}

	if v, ok := u.(*URL); ok {
		if v.url != nil {
			parsed = v
		} else {
			return nil, fmt.Errorf("invalid URL struct")
		}
	} else if v, ok := u.(*url.URL); ok {
		parsed.url = v
	} else if s, ok := u.(string); ok && s != `` {
		if pu, err := url.Parse(s); err == nil {
			parsed.url = pu
		} else {
			return nil, err
		}
	} else if typeutil.IsMap(u) {
		var umap = maputil.M(u)

		parsed.Scheme = umap.String(`scheme`)
		parsed.Host = umap.String(`host`)
		parsed.Path = umap.String(`path`)
		parsed.RawQuery = umap.String(`rawquery`)
		parsed.Fragment = umap.String(`fragment`)
	} else {
		return nil, fmt.Errorf("invalid URL")
	}

	// make sure all struct fields are populated
	if err := parsed.sync(); err == nil {
		return parsed, nil
	} else {
		return nil, err
	}
}

// Take a map or previous URL response structure and encode the values into a string
// that can be used in another URL or form post data.  This command does not automaticlly
// prepend a "?" character to the output.
func (self *Commands) EncodeQuery(querymap interface{}) (string, error) {
	var qmap map[string]interface{}
	var qvals = make(url.Values)

	if u, ok := querymap.(*URL); ok {
		qmap = u.Query
	} else if typeutil.IsMap(querymap) {
		qmap = typeutil.MapNative(querymap)
	} else if s, ok := querymap.(string); ok {
		if u, err := self.Parse(s); err == nil {
			qmap = u.Query
		} else {
			return ``, err
		}
	} else {
		return ``, fmt.Errorf("expected URL string, structure or map of key-value pairs")
	}

	// put the values into the url.Values map
	for k, v := range qmap {
		qvals[k] = []string{typeutil.String(v)}
	}

	// return the encoded values
	return qvals.Encode(), nil
}

// Take a URL or map of query string key=value pairs and return a map of values.
func (self *Commands) ParseQuery(urlOrQueryString interface{}) (map[string]interface{}, error) {

	if u, ok := urlOrQueryString.(*URL); ok {
		return u.Query, nil
	} else if typeutil.IsMap(urlOrQueryString) {
		return typeutil.MapNative(urlOrQueryString), nil
	} else if s, ok := urlOrQueryString.(string); ok {
		if strings.HasPrefix(s, `?`) {
			s = strings.TrimPrefix(s, `?`)

			return urlValuesToMap(s)

		} else if u, err := self.Parse(s); err == nil {
			return u.Query, nil
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("expected URL string, structure or map of key-value pairs")
	}
}

// Escapes the string so it can be safely placed inside a URL path segment, replacing special characters (including /) with %XX sequences as needed.
func (self *Commands) Escape(stringOrStruct interface{}) (string, error) {
	if u, err := self.Parse(stringOrStruct); err == nil {
		return url.PathEscape(u.Path), nil
	} else {
		return ``, err
	}
}

func (self *Commands) Unescape(stringOrStruct interface{}) (string, error) {
	if u, err := self.Parse(stringOrStruct); err == nil {
		return url.PathUnescape(u.Path)
	} else {
		return ``, err
	}
}

func urlValuesToMap(rawquery string) (map[string]interface{}, error) {
	if vals, err := url.ParseQuery(rawquery); err == nil {
		var out = make(map[string]interface{})

		for k, vv := range vals {
			switch len(vv) {
			case 0:
				out[k] = nil
			case 1:
				out[k] = vv[0]
			default:
				out[k] = vv
			}
		}

		return out, nil
	} else {
		return nil, err
	}

}
