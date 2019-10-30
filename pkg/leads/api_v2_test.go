package leads

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeAndDecide(t *testing.T) {
	assert := assert.New(t)

	// phoneTest := HelperRandstring(9)
	// ipTest := HelperRandstring(9)

	tests := []struct {
		Description      string
		TypeRequest      string
		StatusCode       int
		Lead             Lead
		ExpectedResult   bool
		ExpectedResponse ResponseAPI
	}{
		{
			Description: "when HandleFunction receive a POST request with no data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead:        Lead{},
			ExpectedResponse: ResponseAPI{
				Success:     false,
				SmartCenter: false,
			},
		},
		{
			Description: "when HandleFunction receive a POST request without sou_id value",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				IsSmartCenter: false,
			},
			ExpectedResult: false,
			ExpectedResponse: ResponseAPI{
				Success:     false,
				SmartCenter: false,
			},
		},
		// {
		// 	Description: "when HandleFunction receive a POST request without lea_type value",
		// 	TypeRequest: http.MethodPost,
		// 	StatusCode:  http.StatusOK,
		// 	Lead: Lead{
		// 		SouID:         15,
		// 		IsSmartCenter: false,
		// 	},
		// 	ExpectedResponse: ResponseAPI{
		// 		Success:     true,
		// 		SmartCenter: false,
		// 	},
		// },
	}

	for _, test := range tests {

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			testresp, err := json.Marshal(test.Lead)
			assert.NoError(err)

			res.WriteHeader(test.StatusCode)
			res.Write(testresp)
		}))
		defer func() { ts.Close() }()

		body, err := json.Marshal(test.Lead)
		assert.NoError(err)

		req, err := http.NewRequest(test.TypeRequest, ts.URL, bytes.NewBuffer(body))
		assert.NoError(err)

		http := &http.Client{}
		resp, err := http.Do(req)
		assert.NoError(err)

		assert.Equal(resp.StatusCode, test.StatusCode)

		response := ResponseAPI{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(err)

		log.Println(response)

		assert.Equal(response.Success, test.ExpectedResponse.Success)
		assert.Equal(response.SmartCenter, test.ExpectedResponse.SmartCenter)
	}
}
