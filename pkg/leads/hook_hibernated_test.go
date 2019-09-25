package leads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActiveHibernated(t *testing.T) {
	assert := assert.New(t)

	var hibernated Hibernated

	tests := []struct {
		Description    string
		Lead           Lead
		ExpectedResult bool
	}{
		{
			Description: "When an active campaign is used",
			Lead: Lead{
				SouID:        5,
				SouIDLeontel: 6,
			},
			ExpectedResult: false,
		},
		{
			Description: "When an hibernated campaign is used",
			Lead: Lead{
				SouID:        31,
				SouIDLeontel: 39,
			},
			ExpectedResult: true,
		},
		{
			Description: "When a non existing campaign is used",
			Lead: Lead{
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
