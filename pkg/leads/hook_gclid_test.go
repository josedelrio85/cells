package leads

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActiveGclid(t *testing.T) {
	assert := assert.New(t)

	var gc Gclid

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "when Gclid hook is successfully activated",
			Lead: Lead{
				SouID: 9,
			},
			Active: true,
		},
		{
			Description: "when Gclid hook is successfully activated",
			Lead: Lead{
				SouID: 15,
			},
			Active: true,
		},
		{
			Description: "when Gclid hook is not activated",
			Lead: Lead{
				SouID: 1,
			},
			Active: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := gc.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerformGclid(t *testing.T) {
	assert := assert.New(t)

	var gc Gclid

	type ExpectedResult struct {
		Result       bool
		SouID        int64
		SouIDLeontel int64
	}

	database := helperDb()
	database.Open()
	defer database.Close()

	gclidvalue := HelperRandstring(19)

	tests := []struct {
		Description    string
		Lead           Lead
		Response       HookResponse
		ExpectedResult ExpectedResult
	}{
		{
			Description: "When a lead has not gclid",
			Lead: Lead{
				SouID:     9,
				LeatypeID: 1,
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			// When there is not gclid value, the sou_id does not change
			// and consecuently sou_id Leontel will not change either (in this hook)
			ExpectedResult: ExpectedResult{
				Result:       false,
				SouID:        9,
				SouIDLeontel: 0,
			},
		},
		{
			Description: "When a lead has gclid value",
			Lead: Lead{
				SouID:     9,
				Gclid:     &gclidvalue,
				LeatypeID: 1,
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:       true,
				SouID:        9,
				SouIDLeontel: 13,
			}},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			cont := Handler{
				Lead:   test.Lead,
				Storer: &database,
			}

			response := gc.Perform(&cont)

			assert.Equal(test.Response.StatusCode, response.StatusCode)
			assert.Equal(cont.Lead.SouID, test.ExpectedResult.SouID)
			assert.Equal(cont.Lead.SouIDLeontel, test.ExpectedResult.SouIDLeontel)
		})
	}
}
