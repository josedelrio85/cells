package apic2c

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerFunction(t *testing.T) {
	assert := assert.New(t)

	phoneTest := "666666666"
	ipTest := "127.0.0.1"

	tests := []struct {
		Description string
		Storer      Storer
		TypeRequest string
		StatusCode  int
		Lead        Lead
	}{
		{
			Description: "when HandleFunction receive a GET request",
			TypeRequest: http.MethodGet,
			StatusCode:  http.StatusMethodNotAllowed,
			Lead:        Lead{},
		},
		{
			Description: "when HandleFunction receive a POST request",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			Lead:        Lead{},
		},
		{
			Description: "when HandleFunction TODO DESCRIPTION",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			// Storer: &FakeDb{
			// 	OpenFunc:        func() error { return nil },
			// 	CloseFunc:       func() error { return nil },
			// 	CreateTableFunc: func(interface{}) error { return nil },
			// 	UpdateFunc:      func(interface{}, string, []string) error { return nil },
			// 	InsertFunc:			 func(interface{}) error {return nil},
			// },
			Lead: Lead{
				SouID:     15,
				LeatypeID: 1,
				LeaPhone:  &phoneTest,
				LeaIP:     &ipTest,
			},
		},
	}

	for _, test := range tests {
		ch := Handler{}
		ts := httptest.NewServer(ch.ch)
		defer ts.Close()

		body, err := json.Marshal(test.Lead)
		if err != nil {
			t.Errorf("error marshaling test json: Err: %v", err)
			return
		}

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
