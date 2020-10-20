package leads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActiveRejectSC(t *testing.T) {
	assert := assert.New(t)

	var rsc RejectSC

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "when RejectSC hook is successfully activated",
			Lead: Lead{
				SouID: 5,
			},
			Active: false,
		},
		{
			Description: "when RejectSC hook is not activated because source is not in list",
			Lead: Lead{
				SouID: 71,
			},
			Active: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := rsc.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerformRejectSC(t *testing.T) {
	assert := assert.New(t)

	var rsc RejectSC

	// evaluate only source values that are active
	// when rejectsc is active, smart center bool value should be false
	// in other cases, smart center shouldn't be changed
	tests := []struct {
		Description    string
		Lead           Lead
		Response       HookResponse
		ExpectedResult bool
	}{
		{
			Description: "when source should be rejected",
			Lead: Lead{
				SouID:         71,
				IsSmartCenter: true,
			},
			ExpectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			cont := Handler{
				Lead: test.Lead,
			}
			rsc.Perform(&cont)
			assert.Equal(cont.Lead.IsSmartCenter, test.ExpectedResult)
		})
	}
}
