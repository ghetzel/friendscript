package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestPost(t *testing.T) {
	assert := require.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case `/json`:
			var x map[string]interface{}

			json.NewDecoder(req.Body).Decode(&x)
			assert.EqualValues(`this is a test`, x[`k1`])
			assert.EqualValues(`1`, req.Header.Get(`X-Friendscript-Testing`))

		case `/cookies`:
			cookie, err := req.Cookie(`OneTestCookie`)
			assert.NoError(err)
			assert.NotNil(cookie)
			assert.EqualValues(`Greetings!`, cookie.Value)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	var client = New(nil)

	client.Post(fmt.Sprintf("%v/json", server.URL), &RequestArgs{
		Headers: map[string]interface{}{
			`X-Friendscript-Testing`: 1,
		},
		RequestType: `json`,
		Body: map[string]interface{}{
			`k1`: `this is a test`,
		},
	})

	client = New(nil)

	client.Get(fmt.Sprintf("%v/cookies", server.URL), &RequestArgs{
		Cookies: map[string]interface{}{
			`OneTestCookie`: `Greetings!`,
		},
	})
}

func TestIsErrorStatus(t *testing.T) {
	assert := require.New(t)

	for i := 0; i < 999; i++ {
		v := isErrorStatus(i, `402-999,0-400`)

		if i == 401 {
			assert.True(v)
		} else {
			assert.False(v)
		}
	}

	assert.True(isErrorStatus(199, `200`))
	assert.False(isErrorStatus(200, `200`))
	assert.True(isErrorStatus(201, `200`))
	assert.True(isErrorStatus(199, `200-299`))
	assert.False(isErrorStatus(200, `200-299`))
	assert.False(isErrorStatus(299, `200-299`))
	assert.True(isErrorStatus(300, `200-299`))

	assert.True(isErrorStatus(199, `200,204,404`))
	assert.False(isErrorStatus(200, `200,204,404`))
	assert.True(isErrorStatus(201, `200,204,404`))

	assert.True(isErrorStatus(203, `200,204,404`))
	assert.False(isErrorStatus(204, `200,204,404`))
	assert.True(isErrorStatus(205, `200,204,404`))

	assert.True(isErrorStatus(403, `200,204,404`))
	assert.False(isErrorStatus(404, `200,204,404`))
	assert.True(isErrorStatus(405, `200,204,404`))

}
