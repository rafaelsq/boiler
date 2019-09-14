package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JSON(w, r, map[string]interface{}{
			"err": make(chan int),
		})
	}))
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s?debug", ts.URL))
	assert.Nil(t, err)

	b, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	res.Body.Close()

	assert.Equal(t, "could not encode response", string(b))
}
