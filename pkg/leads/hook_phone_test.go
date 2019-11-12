package leads

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestActivePhone(t *testing.T) {
	assert := assert.New(t)

	var ph Phone

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "when Phone hook is successfully activated",
			Lead: Lead{
				SouID: 64,
			},
			Active: true,
		},
		{
			Description: "when Phone hook is successfully activated",
			Lead: Lead{
				SouID: 15,
			},
			Active: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := ph.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerformPhone(t *testing.T) {
	assert := assert.New(t)

	var ph Phone

	p1 := HelperRandstring(9)
	p2 := "6543217890"
	p3 := "01234567987987987987"
	p4 := "666666666"
	p5 := "648921456"

	tests := []struct {
		Description    string
		Lead           Lead
		Response       HookResponse
		ExpectedResult bool
	}{
		{
			Description: "when phone parameter is not present",
			Lead:        Lead{},
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New(""),
			},
			ExpectedResult: false,
		},
		{
			Description: "when a not valid phone parameter is received",
			Lead: Lead{
				LeaPhone: &p1,
			},
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New(""),
			},
			ExpectedResult: false,
		},
		{
			Description: "when a not valid phone parameter is received",
			Lead: Lead{
				LeaPhone: &p2,
			},
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New(""),
			},
			ExpectedResult: false,
		},
		{
			Description: "when a not valid phone parameter is received",
			Lead: Lead{
				LeaPhone: &p3,
			},
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New(""),
			},
			ExpectedResult: false,
		},
		{
			Description: "when a in quarantine phone parameter is received",
			Lead: Lead{
				LeaPhone: &p4,
			},
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New(""),
			},
			ExpectedResult: false,
		},
		{
			Description: "when a ot valid phone parameter is received",
			Lead: Lead{
				LeaPhone: &p5,
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			cont := Handler{
				Lead: test.Lead,
			}
			response := ph.Perform(&cont)
			assert.Equal(test.Response.StatusCode, response.StatusCode)
			if response.StatusCode == 422 {
				assert.False(test.ExpectedResult)
			} else {
				assert.True(test.ExpectedResult)
			}
		})
	}
}
