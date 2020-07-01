package url

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestParseString(t *testing.T) {
	var cases = map[string]*URL{
		`google.com`: {
			Scheme: `https`,
			Host:   `google.com`,
			Domain: `google.com`,
			Port:   443,
			Path:   `/`,
			Full:   `https://google.com/`,
		},
		`http://google.com`: {
			Scheme: `http`,
			Host:   `google.com`,
			Domain: `google.com`,
			Port:   80,
			Path:   `/`,
			Full:   `http://google.com/`,
		},
		`https://google.com`: {
			Scheme: `https`,
			Host:   `google.com`,
			Domain: `google.com`,
			Port:   443,
			Path:   `/`,
			Full:   `https://google.com/`,
		},
		`https://google.com:8443`: {
			Scheme: `https`,
			Host:   `google.com:8443`,
			Domain: `google.com`,
			Port:   8443,
			Path:   `/`,
			Full:   `https://google.com:8443/`,
		},
		`http://google.com:8080`: {
			Scheme: `http`,
			Host:   `google.com:8080`,
			Domain: `google.com`,
			Port:   8080,
			Path:   `/`,
			Full:   `http://google.com:8080/`,
		},
		`http://google.com:443`: {
			Scheme: `http`,
			Host:   `google.com:443`,
			Domain: `google.com`,
			Port:   443,
			Path:   `/`,
			Full:   `http://google.com:443/`,
		},
		`https://google.com:80`: {
			Scheme: `https`,
			Host:   `google.com:80`,
			Domain: `google.com`,
			Port:   80,
			Path:   `/`,
			Full:   `https://google.com:80/`,
		},
	}

	for in, want := range cases {
		have, err := New(nil).Parse(in)
		require.NoError(t, err, in)
		require.NotNil(t, have, in)
		require.Equal(t, want.Full, have.Full, in)
		require.Equal(t, want.Scheme, have.Scheme, in)
		require.Equal(t, want.Host, have.Host, in)
		require.Equal(t, want.Domain, have.Domain, in)
		require.Equal(t, want.Port, have.Port, in)
		require.Equal(t, want.Path, have.Path, in)
		require.Empty(t, want.RawQuery, have.RawQuery, in)
		require.Empty(t, want.Query, have.Query, in)
		require.Empty(t, want.Fragment, have.Fragment, in)
	}
}

func TestEncodeQuery(t *testing.T) {
	var cases = map[interface{}]interface{}{
		`?x=1&y=true&z=three`:                    `x=1&y=true&z=three`,
		`https://example.com?x=1&y=true&z=three`: `x=1&y=true&z=three`,
		&URL{
			Query: map[string]interface{}{
				`z`: `three`,
				`y`: true,
				`x`: 1,
			},
		}: `x=1&y=true&z=three`,
	}

	for in, want := range cases {
		have, err := New(nil).EncodeQuery(in)
		require.NoError(t, err)
		require.Equal(t, want, have)
	}
}

func TestEncodeQueryMap(t *testing.T) {
	have, err := New(nil).EncodeQuery(map[string]interface{}{
		`y`: true,
		`x`: 1,
		`z`: `three`,
	})
	require.NoError(t, err)
	require.Equal(t, `x=1&y=true&z=three`, have)
}
