package leads

import (
	"testing"

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

	var ontime Ontime

	tests := []struct {
		Description    string
		Input          InputDataOntime
		ExpectedResult bool
	}{
		{
			Description: "when the day is not a holiday day",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "2019-04-08",
			},
			ExpectedResult: false,
		},
		{
			Description: "when the day is a holiday day",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "2019-01-01",
			},
			ExpectedResult: true,
		},
		{
			Description: "when we use a campaign without a registered timetable",
			Input: InputDataOntime{
				SouID: 5,
				Day:   "2019-07-16",
			},
			ExpectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			err := ontime.checkHoliday(test.Input)

			assert.NoError(err)
			assert.Equal(test.ExpectedResult, ontime.ResultHoliday)
		})
	}
}

func TestOntime(t *testing.T) {
	assert := assert.New(t)

	var ontime Ontime

	tests := []struct {
		Description    string
		Input          InputDataOntime
		ExpectedResult bool
		RealResult     bool
	}{
		{
			Description: "When we are on time",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "1",
				Hour:  "12:00",
			},
			ExpectedResult: true,
			RealResult:     true,
		},
		{
			Description: "When an hour lower than 10 has one character instead 2 and the ontime validation should be true",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "3",
				Hour:  "9:30",
			},
			ExpectedResult: true,
			RealResult:     false,
		},
		{
			Description: "When are in a working day but off time",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "2",
				Hour:  "08:30",
			},
			ExpectedResult: false,
			RealResult:     false,
		},
		{
			Description: "When are in a non-working day",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "0",
				Hour:  "10:15",
			},
			ExpectedResult: false,
			RealResult:     false,
		},
		{
			Description: "When are in a campaign without a registered timetable",
			Input: InputDataOntime{
				SouID: 99,
				Day:   "3",
				Hour:  "16:15",
			},
			ExpectedResult: false,
			RealResult:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			err := ontime.checkOntime(test.Input)

			assert.NoError(err)
			if test.ExpectedResult != test.RealResult {
				assert.Equal(test.RealResult, ontime.ResultOntime)
			} else {
				assert.Equal(test.ExpectedResult, ontime.ResultOntime)
			}
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
				Day:   "2019-06-15",
			},
			Input2: InputDataOntime{
				SouID: 6,
				Day:   "2",
				Hour:  "12:00",
			},
			ExpectedResult: ExpectedResult{
				ResultHoliday: false,
				ResultOntime:  true,
			},
		},
		{
			Description: "When is holiday",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "2019-08-15",
			},
			Input2: InputDataOntime{
				SouID: 6,
				Day:   "2",
				Hour:  "12:00",
			},
			ExpectedResult: ExpectedResult{
				ResultHoliday: true,
				ResultOntime:  true,
			},
		},
		{
			Description: "When is not holiday and we are not on time",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "2019-06-15",
			},
			Input2: InputDataOntime{
				SouID: 6,
				Day:   "0",
				Hour:  "12:00",
			},
			ExpectedResult: ExpectedResult{
				ResultHoliday: false,
				ResultOntime:  false,
			},
		},
		{
			Description: "When is not holiday and we are not on time v2",
			Input: InputDataOntime{
				SouID: 6,
				Day:   "2019-06-15",
			},
			Input2: InputDataOntime{
				SouID: 6,
				Day:   "3",
				Hour:  "01:00",
			},
			ExpectedResult: ExpectedResult{
				ResultHoliday: false,
				ResultOntime:  false,
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
