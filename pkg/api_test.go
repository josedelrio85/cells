package apic2c

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHandlerFunction(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Description string
		// Storer Storer
		TypeRequest string
		StatusCode  int
	}{
		{
			Description: "when HandleFunction receive a GET request",
			TypeRequest: http.MethodGet,
			StatusCode:  http.StatusMethodNotAllowed,
		},
		{
			Description: "when HandleFunction receive a POST request",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
		},
	}

	for _, test := range tests {
		ch := Handler{}
		ts := httptest.NewServer(ch.HandleFunction())
		defer ts.Close()

		var body []byte
		req, err := http.NewRequest(test.TypeRequest, ts.URL, bytes.NewBuffer(body))
		if err != nil {
			t.Errorf("error createing the test Request: err %v", err)
			return
		}

		http := &http.Client{}
		resp, err := http.Do(req)
		if err != nil {
			t.Errorf("error sending test request: Err %v", err)
			return
		}

		assert.Equal(resp.StatusCode, test.StatusCode)
	}
}
