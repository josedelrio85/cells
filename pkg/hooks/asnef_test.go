package leads

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/bysidecar/leads/pkg/model"
)

func TestActive(t *testing.T) {
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
				SouID:         9,
				IsSmartCenter: true,
			},
			Active: true,
		},
		{
			Description: "when Asnef successfully is activated",
			Lead: model.Lead{
				SouID:         58,
				IsSmartCenter: true,
			},
			Active: true,
		},
		{
			Description: "when Asnef is not activated",
			Lead: model.Lead{
				SouID:         99,
				IsSmartCenter: true,
			},
			Active: false,
		},
		{
			Description: "when Asnef is not activated",
			Lead: model.Lead{
				SouID:         0,
				IsSmartCenter: true,
			},
			Active: false,
		},
		{
			Description: "when IsSmartCenter is false",
			Lead: model.Lead{
				SouID:         9,
				IsSmartCenter: false,
			},
			Active: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := asnef.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerform(t *testing.T) {
	assert := assert.New(t)

	var asnef Asnef
	phone := "665932355"
	dni := "79317432t"
	motivo := "Check Asnef marcado."
	cantidad := "1000"

	lead := model.Lead{
		SouID: 9,
	}
	candidates := GetCandidates(lead)

	tests := []struct {
		Description string
		Lead        model.Lead
		Response    HookResponse
	}{
		{
			Description: "when Asnef successfully checks and gives ok to the lead",
			Lead: model.Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: model.Creditea{
					Motivo: nil,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
				Result:     false,
			},
		},
		{
			Description: "when Asnef/Already client checks were clicked. Client activates the limitation",
			Lead: model.Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: model.Creditea{
					Cantidadsolicitada: &cantidad,
					Motivo:             &motivo,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
				Result:     true,
			},
		},
		{
			Description: "when Asnef validations is not passed",
			Lead: model.Lead{
				SouID:         9,
				LeaPhone:      &candidates[0].Telefono,
				LeaDNI:        &candidates[0].DNI,
				IsSmartCenter: false,
				Creditea: model.Creditea{
					Motivo: nil,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
				Result:     true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			response := asnef.Perform(&test.Lead)

			assert.Equal(test.Response, response)
			assert.Equal(test.Response.Result, response.Result)
			assert.Equal(test.Response.StatusCode, response.StatusCode)
			assert.Equal(test.Response.Err, response.Err)
		})
	}
}
