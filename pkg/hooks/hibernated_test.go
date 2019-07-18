package leads

import (
	"testing"

	model "github.com/bysidecar/leads/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestActiveHibernated(t *testing.T) {
	assert := assert.New(t)

	var hibernated Hibernated

	tests := []struct {
		Description    string
		Lead           model.Lead
		ExpectedResult bool
	}{
		{
			Description: "When an active campaign is used",
			Lead: model.Lead{
				SouID:        5,
				SouIDLeontel: 6,
			},
			ExpectedResult: false,
		},
		{
			Description: "When an hibernated campaign is used",
			Lead: model.Lead{
				SouID:        31,
				SouIDLeontel: 39,
			},
			ExpectedResult: true,
		},
		{
			Description: "When a non existing campaign is used",
			Lead: model.Lead{
				SouID:        999,
				SouIDLeontel: 0,
			},
			ExpectedResult: true,
		},
	}

	for _, test := range tests {
		result := hibernated.Active(test.Lead)

		assert.Equal(result, test.ExpectedResult)
	}
}
