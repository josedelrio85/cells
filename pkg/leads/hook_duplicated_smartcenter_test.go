package leads

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActiveDuplicatedSC(t *testing.T) {
	assert := assert.New(t)

	var duplicated DuplicatedSmartCenter

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "when DuplicatedSmartCenter hook is successfully activated",
			Lead: Lead{
				SouID: 64,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedSmartCenter hook is successfully activated",
			Lead: Lead{
				SouID: 15,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedSmartCenter hook is successfully activated",
			Lead: Lead{
				SouID: 63,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedSmartCenter hook is successfully activated",
			Lead: Lead{
				SouID: 70,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedSmartCenter hook is successfully activated",
			Lead: Lead{
				SouID: 74,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedSmartCenter hook is not activated",
			Lead: Lead{
				SouID: 1,
			},
			Active: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := duplicated.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerformDuplicatedSC(t *testing.T) {
	assert := assert.New(t)

	phone1 := HelperRandstring(9)
	phone2 := HelperRandstring(9)
	phone3 := HelperRandstring(9)

	tests := []struct {
		Description      string
		Lead             Lead
		ExpectedStatus   int
		ExpectedResponse RespSC
	}{
		{
			Description: "When a lead is not duplicated in smart center",
			Lead: Lead{
				LeaPhone:         &phone1,
				SouIDLeontel:     73,
				LeatypeIDLeontel: 2,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: RespSC{
				Success: false,
				Data:    nil,
				Error:   nil,
			},
		},
		{
			Description: "When another lead is not duplicated in smartcenter",
			Lead: Lead{
				LeaPhone:  &phone2,
				SouID:     23,
				LeatypeID: 26,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: RespSC{
				Success: false,
				Data:    nil,
				Error:   nil,
			},
		},
		{
			Description: "When a lead is duplicated because is open in Leontel environment",
			Lead: Lead{
				LeaPhone:  &phone3,
				SouID:     73,
				LeatypeID: 2,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: RespSC{
				Success: true,
				Data: []DataLeontelResp{
					DataLeontelResp{
						Closed:      "0",
						LeaID:       "999999",
						LeaTs:       time.Now().Format("2006-01-02 15:04:05"),
						Description: "Test",
						SubID:       "99",
					},
				},
				Error: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			// generate a test server so we can capture and inspect the request
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				testresp, err := json.Marshal(test.ExpectedResponse)
				assert.NoError(err)

				// Writes the status response that we expect for the actual test
				res.WriteHeader(test.ExpectedStatus)

				// Writes the response that we expect for the actual test
				res.Write(testresp)
			}))
			defer func() { testServer.Close() }()

			// Pass the leontel data for souid and type
			data := struct {
				Phone   string `json:"TELEFONO"`
				SouID   int64  `json:"lea_source"`
				LeaType int64  `json:"lea_type"`
			}{
				Phone:   *test.Lead.LeaPhone,
				SouID:   test.Lead.SouIDLeontel,
				LeaType: test.Lead.LeatypeIDLeontel,
			}

			bytevalues, err := json.Marshal(data)
			assert.NoError(err)

			req, err := http.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(bytevalues))
			assert.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(err)

			rawdata, _ := ioutil.ReadAll(resp.Body)
			structdata := RespSC{}

			err = json.Unmarshal(rawdata, &structdata)
			assert.NoError(err)

			assert.Equal(test.ExpectedResponse.Success, structdata.Success)
			assert.Equal(test.ExpectedResponse.Data, structdata.Data)
			assert.Equal(test.ExpectedResponse.Error, structdata.Error)
		})
	}
}
