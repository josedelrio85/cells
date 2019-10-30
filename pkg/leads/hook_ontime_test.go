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
