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

func TestActiveOntime(t *testing.T) {
	assert := assert.New(t)

	var ontime Ontime

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "When LeatypeID is a type susceptible of being checked",
			Lead: Lead{
				LeatypeID: 1,
			},
			Active: true,
		},
		{
			Description: "When LeatypeID is a type susceptible of being checked",
			Lead: Lead{
				LeatypeID: 3,
			},
			Active: true,
		},
		{
			Description: "When LeatypeID is a type susceptible of being checked",
			Lead: Lead{
				LeatypeID: 4,
			},
			Active: true,
		},
		{
			Description: "When LeatypeID is a type susceptible of being checked",
			Lead: Lead{
				LeatypeID: 9,
			},
			Active: true,
		},
		{
			Description: "When LeatypeID is not susceptible of being checked",
			Lead: Lead{
				LeatypeID: 99,
			},
			Active: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			active := ontime.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestOnholiday(t *testing.T) {
	assert := assert.New(t)

	result := getExpectedResultHoliday()

	tests := []struct {
		Description    string
		Input          InputDataOntime
		ExpectedResult bool
		ExpectedStatus int
	}{
		{
			Description: "if today is holiday the result will be true, false other case",
			Input: InputDataOntime{
				SouID: 6,
			},
			ExpectedResult: result,
			ExpectedStatus: http.StatusOK,
		},
		{
			Description: "when we use a campaign without a registered timetable",
			Input: InputDataOntime{
				SouID: 999,
			},
			ExpectedResult: false,
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			// generate a test server so we can capture and inspect the request
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				testresp, err := json.Marshal(test.ExpectedResult)
				assert.NoError(err)

				// Writes the status response that we expect for the actual test
				res.WriteHeader(test.ExpectedStatus)

				// Writes the response that we expect for the actual test
				res.Write(testresp)
			}))
			defer func() { testServer.Close() }()

			bytevalues, err := json.Marshal(test.Input)
			assert.NoError(err)

			req, err := http.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(bytevalues))
			assert.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(err)

			rawdata, _ := ioutil.ReadAll(resp.Body)
			var structdata bool

			err = json.Unmarshal(rawdata, &structdata)

			assert.NoError(err)
			assert.Equal(test.ExpectedResult, structdata)
		})
	}
}

func TestOntime(t *testing.T) {
	assert := assert.New(t)

	result := getExpectedResultOnTime()

	tests := []struct {
		Description    string
		Input          InputDataOntime
		ExpectedResult bool
		ExpectedStatus int
	}{
		{
			Description: "When we are on time",
			Input: InputDataOntime{
				SouID: 6,
			},
			ExpectedResult: result,
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			// generate a test server so we can capture and inspect the request
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				testresp, err := json.Marshal(test.ExpectedResult)
				assert.NoError(err)

				// Writes the status response that we expect for the actual test
				res.WriteHeader(test.ExpectedStatus)

				// Writes the response that we expect for the actual test
				res.Write(testresp)
			}))
			defer func() { testServer.Close() }()

			bytevalues, err := json.Marshal(test.Input)
			assert.NoError(err)

			req, err := http.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(bytevalues))
			assert.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(err)

			rawdata, _ := ioutil.ReadAll(resp.Body)
			var structdata bool

			err = json.Unmarshal(rawdata, &structdata)

			assert.NoError(err)
			assert.Equal(test.ExpectedResult, structdata)
		})
	}
}

type ExpectedResult struct {
	ResultHoliday bool
	ResultOntime  bool
}

func TestPerformOntime(t *testing.T) {
	assert := assert.New(t)
	var ontime Ontime

	rhol := getExpectedResultHoliday()
	rot := getExpectedResultOnTime()

	tests := []struct {
		Description    string
		Input          InputDataOntime
		Input2         InputDataOntime
		ExpectedResult ExpectedResult
	}{
		{
			Description: "When is not holiday and we are on time",
			Input: InputDataOntime{
				SouID: 6,
			},
			Input2: InputDataOntime{
				SouID: 6,
			},
			ExpectedResult: ExpectedResult{
				ResultHoliday: rhol,
				ResultOntime:  rot,
			},
		},
	}

	for _, test := range tests {
		err := ontime.checkHoliday(test.Input)

		assert.NoError(err)

		err = ontime.checkOntime(test.Input2)

		assert.NoError(err)

		assert.Equal(ontime.ResultHoliday, test.ExpectedResult.ResultHoliday)
		assert.Equal(ontime.ResultOntime, test.ExpectedResult.ResultOntime)
	}
}

func getExpectedResultHoliday() bool {
	today := time.Now().Format("2006-01-02")
	holdays := map[string]bool{
		"2019-01-01": true,
		"2019-01-06": true,
		"2019-04-10": true,
		"2019-05-01": true,
		"2019-08-15": true,
		"2019-10-12": true,
		"2019-12-08": true,
		"2019-12-25": true,
	}
	if holdays[today] {
		return true
	}
	return false
}

func getExpectedResultOnTime() bool {
	loc, _ := time.LoadLocation("UTC")
	intday := time.Now().In(loc).Weekday()
	inthour := time.Now().In(loc).Hour()

	if intday > 0 && intday < 6 {
		if inthour > 9 && inthour < 21 {
			return true
		}
	}
	return false
}
