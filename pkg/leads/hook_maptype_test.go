package leads

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestActiveMapType(t *testing.T) {
	assert := assert.New(t)

	var mp MapType

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "when MapType hook is successfully activated",
			Lead: Lead{
				LeatypeID: 2,
				SouID:     74,
			},
			Active: true,
		},
		{
			Description: "when MapType hook is successfully activated",
			Lead: Lead{
				LeatypeID: 8,
				SouID:     75,
			},
			Active: true,
		},
		{
			Description: "when MapType hook is successfully activated",
			Lead: Lead{
				LeatypeID: 24,
				SouID:     76,
			},
			Active: true,
		},
		{
			Description: "when MapType hook is successfully activated",
			Lead: Lead{
				LeatypeID: 8,
				SouID:     77,
			},
			Active: true,
		},
		{
			Description: "when MapType hook is not activated because type is not in list",
			Lead: Lead{
				LeatypeID: 99,
				SouID:     74,
			},
			Active: false,
		},
		{
			Description: "when Gclid hook is not activated because sou_id is not in list",
			Lead: Lead{
				LeatypeID: 99,
				SouID:     75,
			},
			Active: false,
		},
		{
			Description: "when MapType hook is successfully activated",
			Lead: Lead{
				LeatypeID: 3,
				SouID:     64,
			},
			Active: true,
		},
		{
			Description: "when MapType hook is successfully activated",
			Lead: Lead{
				LeatypeID: 8,
				SouID:     79,
			},
			Active: true,
		},
		{
			Description: "when MapType hook is successfully activated",
			Lead: Lead{
				LeatypeID: 8,
				SouID:     78,
			},
			Active: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := mp.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerformMapType(t *testing.T) {
	assert := assert.New(t)

	var mt MapType

	type ExpectedResult struct {
		Result             bool
		SouID              int64
		LeaTypeID          int64
		SouIDLeontel       int64
		SouDescLeontel     string
		LeatypeIDLeontel   int64
		LeatypeDescLeontel string
	}

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_5")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_5")
	defer db.Close()

	database := Database{}
	database.DB = db

	tests := []struct {
		Description    string
		Lead           Lead
		Helper         Lead
		Response       HookResponse
		ExpectedResult ExpectedResult
	}{
		{
			Description: "When a lead is not from active campaign",
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
			ExpectedResult: ExpectedResult{
				Result:    false,
				SouID:     64,
				LeaTypeID: 1,
			},
		},
		{
			Description: "When a lead is from active campaign but type it is not",
			Lead: Lead{
				SouID:     75,
				LeatypeID: 10,
			},
			Helper: Lead{
				SouIDLeontel:       84,
				SouDescLeontel:     "T EMPRESAS E2E",
				LeatypeIDLeontel:   14,
				LeatypeDescLeontel: "PENDIENTE FIRMA",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:    false,
				SouID:     75,
				LeaTypeID: 10,
			},
		},
		{
			Description: "When a lead is from active campaign and type is a valid type",
			Lead: Lead{
				SouID:     74,
				LeatypeID: 2,
			},
			Helper: Lead{
				SouIDLeontel:       83,
				SouDescLeontel:     "R EMPRESAS E2E",
				LeatypeIDLeontel:   13,
				LeatypeDescLeontel: "C2C",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:    true,
				SouID:     74,
				LeaTypeID: 9,
			},
		},
		{
			Description: "When a lead is from active campaign and type is a valid type",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 27,
			},
			Helper: Lead{
				SouIDLeontel:       73,
				SouDescLeontel:     "R CABLE E2E",
				LeatypeIDLeontel:   30,
				LeatypeDescLeontel: "C2C",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:    true,
				SouID:     64,
				LeaTypeID: 9,
			},
		},
		{
			Description: "When a lead is from active campaign and type is a valid type",
			Lead: Lead{
				SouID:     79,
				LeatypeID: 8,
			},
			Helper: Lead{
				SouIDLeontel:       88,
				SouDescLeontel:     "VIRGIN",
				LeatypeIDLeontel:   13,
				LeatypeDescLeontel: "C2C",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:    true,
				SouID:     79,
				LeaTypeID: 9,
			},
		},
		{
			Description: "When a lead is from active campaign and type is a valid type",
			Lead: Lead{
				SouID:     78,
				LeatypeID: 8,
			},
			Helper: Lead{
				SouIDLeontel:       87,
				SouDescLeontel:     "ENDESA",
				LeatypeIDLeontel:   13,
				LeatypeDescLeontel: "C2C",
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: ExpectedResult{
				Result:    true,
				SouID:     78,
				LeaTypeID: 9,
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

			mt.Perform(&cont)

			assert.Equal(cont.Lead.LeatypeID, test.ExpectedResult.LeaTypeID)
			assert.Equal(cont.Lead.SouID, test.ExpectedResult.SouID)
		})
	}
}

func TestGetNewType(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Description    string
		Lead           Lead
		ExpectedResult int64
	}{
		{
			Description: "When a lead is not from active campaign",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 1,
			},
			ExpectedResult: 1,
		},
		{
			Description: "When a lead is from active campaign but lea type it is not",
			Lead: Lead{
				SouID:     77,
				LeatypeID: 1,
			},
			ExpectedResult: 1,
		},
		{
			Description: "When a lead is from active campaign but lea type it is not",
			Lead: Lead{
				SouID:     75,
				LeatypeID: 10,
			},
			ExpectedResult: 10,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     77,
				LeatypeID: 8,
			},
			ExpectedResult: 20,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     74,
				LeatypeID: 24,
			},
			ExpectedResult: 9,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 24,
			},
			ExpectedResult: 9,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 27,
			},
			ExpectedResult: 9,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 30,
			},
			ExpectedResult: 9,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 8,
			},
			ExpectedResult: 9,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     79,
				LeatypeID: 8,
			},
			ExpectedResult: 9,
		},
		{
			Description: "When a lead is from active campaign and lea type is an expected value",
			Lead: Lead{
				SouID:     78,
				LeatypeID: 8,
			},
			ExpectedResult: 9,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			getNewType(&test.Lead)

			assert.Equal(test.Lead.LeatypeID, test.ExpectedResult)
		})
	}
}
