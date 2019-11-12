package leads

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
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
		Result             bool
		SouID              int64
		SouIDLeontel       int64
		SouDescLeontel     string
		LeatypeIDLeontel   int64
		LeatypeDescLeontel string
	}

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_4")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_4")
	defer db.Close()

	database := Database{}
	database.DB = db

	gclidvalue := HelperRandstring(19)

	tests := []struct {
		Description    string
		Lead           Lead
		Helper         Lead
		Response       HookResponse
		ExpectedResult ExpectedResult
	}{
		{
			Description: "When a lead has not gclid",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 1,
			},
			Helper: Lead{
				SouIDLeontel:       73,
				SouDescLeontel:     "R CABLE END TO END",
				LeatypeIDLeontel:   2,
				LeatypeDescLeontel: "C2C",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			// When there is not gclid value, the sou_id does not change
			// and consecuently sou_id Leontel will not change either (in this hook)
			ExpectedResult: ExpectedResult{
				Result:       false,
				SouID:        64,
				SouIDLeontel: 0,
			},
		},
		{
			Description: "When a lead has gclid value",
			Lead: Lead{
				SouID:     15,
				Gclid:     &gclidvalue,
				LeatypeID: 1,
			},
			Helper: Lead{
				SouIDLeontel:       23,
				SouDescLeontel:     "ALTERNA",
				LeatypeIDLeontel:   2,
				LeatypeDescLeontel: "C2C",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:       true,
				SouID:        15,
				SouIDLeontel: 23,
			},
		},
		{
			Description: "When a lead has gclid value",
			Lead: Lead{
				SouID:     64,
				Gclid:     &gclidvalue,
				LeatypeID: 1,
			},
			Helper: Lead{
				SouIDLeontel:       73,
				SouDescLeontel:     "R CABLE END TO END",
				LeatypeIDLeontel:   2,
				LeatypeDescLeontel: "C2C",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:       true,
				SouID:        64,
				SouIDLeontel: 73,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			if test.Lead.Gclid != nil {
				row := fmt.Sprintf("%d,%s,%d", test.Lead.SouID, test.Helper.SouDescLeontel, test.Helper.SouIDLeontel)
				rs := mock.NewRows([]string{"sou_id", "sou_description", "sou_idcrm"}).
					FromCSVString(row)

				mock.ExpectQuery("SELECT (.+)").
					WithArgs(test.Lead.SouID).
					WillReturnRows(rs)

				row2 := fmt.Sprintf("%d,%s,%d", test.Lead.LeatypeID, test.Helper.LeatypeDescLeontel, test.Helper.LeatypeIDLeontel)
				rs2 := mock.NewRows([]string{"leatype_id", "leatype_description", "leatype_idcrm"}).
					FromCSVString(row2)

				mock.ExpectQuery("SELECT (.+)").
					WithArgs(test.Lead.LeatypeID).
					WillReturnRows(rs2)
			}

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
