package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestDecodeAndDecide(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Description      string
		TypeRequest      string
		StatusCode       int
		Lead             Lead
		ExpectedResult   bool
		ExpectedResponse ResponseAPI
	}{
		{
			Description: "when HandleFunction receive a POST request with no data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead:        Lead{},
			ExpectedResponse: ResponseAPI{
				Success:     false,
				SmartCenter: false,
			},
		},
		{
			Description: "when HandleFunction receive a POST request without sou_id value",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				IsSmartCenter: false,
			},
			ExpectedResult: false,
			ExpectedResponse: ResponseAPI{
				Success:     false,
				SmartCenter: false,
			},
		},
		{
			Description: "when HandleFunction receive a POST request without lea_type value",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			Lead: Lead{
				SouID:         15,
				IsSmartCenter: false,
			},
			ExpectedResponse: ResponseAPI{
				Success:     true,
				SmartCenter: false,
			},
		},
	}

	for _, test := range tests {

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			testresp, err := json.Marshal(test.ExpectedResponse)
			assert.NoError(err)

			res.WriteHeader(test.StatusCode)
			res.Write(testresp)
		}))
		defer func() { ts.Close() }()

		body, err := json.Marshal(test.Lead)
		assert.NoError(err)

		req, err := http.NewRequest(test.TypeRequest, ts.URL, bytes.NewBuffer(body))
		assert.NoError(err)

		http := &http.Client{}
		resp, err := http.Do(req)
		assert.NoError(err)

		assert.Equal(resp.StatusCode, test.StatusCode)

		response := ResponseAPI{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(err)

		log.Println(response)

		assert.Equal(response.Success, test.ExpectedResponse.Success)
		assert.Equal(response.SmartCenter, test.ExpectedResponse.SmartCenter)
	}
}

func TestGetSourceValues(t *testing.T) {
	assert := assert.New(t)

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_3")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_3")
	defer db.Close()

	tests := []struct {
		Description    string
		Lead           Lead
		ExpectedResult Lead
	}{
		{
			Description: "RCABLE END TO END	64 => 73 | C2C 1 => 2",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 1,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       73,
				SouDescLeontel:     "RCABLE END TO END",
				LeatypeIDLeontel:   2,
				LeatypeDescLeontel: "C2C",
				SouIDEvolution:     0,
			},
		},
		{
			Description: "EVO BANCO 3 => 4 | INACTIVIDAD 2 => 3",
			Lead: Lead{
				SouID:     3,
				LeatypeID: 2,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       4,
				SouDescLeontel:     "EVO BANCO",
				LeatypeIDLeontel:   3,
				LeatypeDescLeontel: "INACTIVIDAD",
				SouIDEvolution:     0,
			},
		},
		{
			Description: "R CABLE EXPANSION END TO END 54 => 63 | FDH 8 => 12",
			Lead: Lead{
				SouID:     54,
				LeatypeID: 8,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       63,
				SouDescLeontel:     "R CABLE EXPANSION END TO END",
				LeatypeIDLeontel:   12,
				LeatypeDescLeontel: "FDH",
				SouIDEvolution:     0,
			},
		},
		{
			Description: "R CABLE END TO END 64 => 73 | SEM 25 => 27",
			Lead: Lead{
				SouID:     64,
				LeatypeID: 25,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       73,
				SouDescLeontel:     "R CABLE END TO END",
				LeatypeIDLeontel:   27,
				LeatypeDescLeontel: "SEM",
				SouIDEvolution:     0,
			},
		},
		{
			Description: "ADESLAS E2E 77 => 86 | C2C 1 => 1",
			Lead: Lead{
				SouID:     77,
				LeatypeID: 1,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       86,
				SouDescLeontel:     "ADESLAS END TO END",
				LeatypeIDLeontel:   2,
				LeatypeDescLeontel: "C2C",
				SouIDEvolution:     0,
			},
		},
		{
			Description: "R CABLE END TO END 74 => 83 | SEM 25 => 27",
			Lead: Lead{
				SouID:     74,
				LeatypeID: 25,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       83,
				SouDescLeontel:     "R CABLE END TO END",
				LeatypeIDLeontel:   27,
				LeatypeDescLeontel: "SEM",
				SouIDEvolution:     0,
			},
		},
		{
			Description: "ENDESA 78 => 87 | CORREGISTRO 30 => 33",
			Lead: Lead{
				SouID:     78,
				LeatypeID: 30,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       87,
				SouDescLeontel:     "ENDESA",
				LeatypeIDLeontel:   33,
				LeatypeDescLeontel: "CORREGISTRO",
				SouIDEvolution:     0,
			},
		},
		{
			Description: "VIRGIN 79 => 88 | CORREGISTRO 30 => 33",
			Lead: Lead{
				SouID:     79,
				LeatypeID: 30,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       88,
				SouDescLeontel:     "VIRGIN",
				LeatypeIDLeontel:   33,
				LeatypeDescLeontel: "CORREGISTRO",
				SouIDEvolution:     100000006,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			row := fmt.Sprintf("%d,%s,%d,%d", test.Lead.SouID, test.ExpectedResult.SouDescLeontel, test.ExpectedResult.SouIDLeontel, test.ExpectedResult.SouIDEvolution)
			rs := mock.NewRows([]string{"sou_id", "sou_description", "sou_idcrm", "sou_id_evolution"}).
				FromCSVString(row)

			mock.ExpectQuery("SELECT (.+)").
				WithArgs(test.Lead.SouID).
				WillReturnRows(rs)

			row2 := fmt.Sprintf("%d,%s,%d", test.Lead.LeatypeID, test.ExpectedResult.LeatypeDescLeontel, test.ExpectedResult.LeatypeIDLeontel)
			rs2 := mock.NewRows([]string{"leatype_id", "leatype_description", "leatype_idcrm"}).
				FromCSVString(row2)

			mock.ExpectQuery("SELECT (.+)").
				WithArgs(test.Lead.LeatypeID).
				WillReturnRows(rs2)

			err := test.Lead.GetSourceValues(db)
			assert.NoError(err)

			assert.Equal(test.ExpectedResult.SouIDLeontel, test.Lead.SouIDLeontel)
			assert.Equal(test.ExpectedResult.LeatypeIDLeontel, test.Lead.LeatypeIDLeontel)
			assert.Equal(test.ExpectedResult.SouIDEvolution, test.Lead.SouIDEvolution)
		})
	}
}

func TestGetPassport(t *testing.T) {
	assert := assert.New(t)

	ip := HelperRandstring(14)

	tests := []struct {
		Description    string
		StatusCode     int
		Interaction    Interaction
		ExpectedResult Lead
	}{
		{
			Description: "RCABLE END TO END	64 => 73 | C2C 1 => 2",
			StatusCode: http.StatusOK,
			Interaction: Interaction{
				Provider:    "RCABLE END TO END",
				Application: "C2C",
				IP:          ip,
			},
			ExpectedResult: Lead{
				PassportID:    HelperRandstring(12),
				PassportIDGrp: HelperRandstring(12),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				testresp, err := json.Marshal(test.ExpectedResult)
				assert.NoError(err)

				res.WriteHeader(test.StatusCode)
				res.Write(testresp)
			}))
			defer func() { ts.Close() }()

			bytevalues, err := json.Marshal(test.Interaction)
			assert.NoError(err)

			req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(bytevalues))
			assert.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(err)

			rawdata, _ := ioutil.ReadAll(resp.Body)

			passport := Passport{}
			err = json.Unmarshal(rawdata, &passport)
			assert.NoError(err)

			assert.Equal(test.ExpectedResult.PassportID, passport.PassportID)
			assert.Equal(test.ExpectedResult.PassportIDGrp, passport.PassportIDGrp)
		})
	}
}
