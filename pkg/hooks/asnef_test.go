package leads

import (
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/bysidecar/leads/pkg/model"
)

func ActiveTest(t *testing.T) {
	assert := assert.New(t)

	var asnef Asnef

	tests := []struct {
		Description string
		Lead        model.Lead
		Active      bool
	}{
		{
			Description: "when Asnef successfully is activated",
			Lead: model.Lead{
				SouID: 1,
			},
			Active: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := asnef.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func PerformTest(t *testing.T) {
	assert := assert.New(t)

	var asnef Asnef

	tests := []struct {
		Description string
		Lead        model.Lead
		Response    HookResponse
	}{
		{
			Description: "when Asnef successfully checks and gives ok to the lead",
			Lead:        model.Lead{},
			Response:    HookResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			response := asnef.Perform(test.Lead)

			assert.Equal(test.Response, response)
		})
	}
}
