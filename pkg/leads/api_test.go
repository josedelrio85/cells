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

func TestGetLeontelValues(t *testing.T) {
	assert := assert.New(t)

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_0")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_0")
	defer db.Close()

	tests := []struct {
		Description    string
		Lead           Lead
		ExpectedResult Lead
	}{
		{
			Description: "CREDITEA END TO END	9 => 13 | C2C 1 => 2",
			Lead: Lead{
				SouID:     9,
				LeatypeID: 1,
			},
			ExpectedResult: Lead{
				SouIDLeontel:       13,
				SouDescLeontel:     "CREDITEA END TO END",
				LeatypeIDLeontel:   2,
				LeatypeDescLeontel: "C2C",
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
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			row := fmt.Sprintf("%d,%s,%d", test.Lead.SouID, test.ExpectedResult.SouDescLeontel, test.ExpectedResult.SouIDLeontel)
			rs := mock.NewRows([]string{"sou_id", "sou_description", "sou_idcrm"}).
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

			err := test.Lead.GetLeontelValues(db)
			assert.NoError(err)

			assert.Equal(test.ExpectedResult.SouIDLeontel, test.Lead.SouIDLeontel)
			assert.Equal(test.ExpectedResult.LeatypeIDLeontel, test.Lead.LeatypeIDLeontel)
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
			Description: "CREDITEA END TO END	9 => 13 | C2C 1 => 2",
			StatusCode: http.StatusOK,
			Interaction: Interaction{
				Provider:    "CREDITEA END TO END",
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

func TestLeadToLeontel(t *testing.T) {
	assert := assert.New(t)

	t1 := "A"
	t2 := "B"
	t3 := "C"
	t4 := "D"
	t5 := "E"
	t6 := "F"
	t7 := "G"
	t8 := ""
	observations := "H"

	obsTest := "A -- B"
	obsTest2 := "A -- C -- D -- B"
	obsTest3 := "A -- B -- C -- D -- E -- F"
	obsTest4 := "A -- B -- C -- D -- E -- F -- G -- A -- B -- C -- D -- E"
	obsTest5 := "A -- B -- C"
	obsTest6 := "B -- C -- D -- E -- F -- G"
	obsTest7 := "B -- C -- D -- E -- F"

	obsTest8 := fmt.Sprintf("%s -- %s -- %s -- %s", obsTest6, obsTest6, obsTest7, obsTest)

	tests := []struct {
		Index          int
		Description    string
		Lead           Lead
		ExpectedResult LeadLeontel
	}{
		{
			Index:       1,
			Description: "check data returned for sou_id 9 Creditea EndToEnd",
			Lead: Lead{
				SouID:        9,
				SouIDLeontel: 13,
				LeaDNI:       &t1,
				Creditea: &Creditea{
					RequestedAmount: &t2,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     13,
				Dninie:        &t1,
				Observaciones: &obsTest,
			},
		},
		{
			Index:       2,
			Description: "check data returned for sou_id 11 Creditea Rastreator",
			Lead: Lead{
				SouID:        11,
				SouIDLeontel: 15,
				LeaDNI:       &t1,
				Creditea: &Creditea{
					RequestedAmount: &t2,
					NetIncome:       &t3,
					ContractType:    &t4,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     15,
				Dninie:        &t1,
				Observaciones: &obsTest2,
			},
		},
		{
			Index:       3,
			Description: "check data returned for sou_id 46-49 Microsoft Hazelcambio + Recomendador",
			Lead: Lead{
				SouID:        46,
				SouIDLeontel: 61,
				Observations: &observations,
				Microsoft: &Microsoft{
					ComputerType: &t1,
					Sector:       &t2,
					Usecase:      &t3,
					Budget:       &t4,
					Performance:  &t5,
					Movility:     &t6,
					Office365:    &t7,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:      61,
				Tipoordenador:  &t1,
				Sector:         &t2,
				Tipouso:        &t3,
				Presupuesto:    &t4,
				Rendimiento:    &t5,
				Movilidad:      &t6,
				Office365:      &t7,
				Observaciones2: &observations,
			},
		},
		{
			Index:       4,
			Description: "check data returned for sou_id 48 Microsoft Calculadora",
			Lead: Lead{
				SouID:        48,
				SouIDLeontel: 61,
				Observations: &observations,
				Microsoft: &Microsoft{
					DevicesAverageAge:      &t1,
					DevicesOperatingSystem: &t2,
					DevicesHangFrequency:   &t3,
					DevicesNumber:          &t4,
					DevicesLastYearRepairs: &t5,
					DevicesStartupTime:     &t6,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:      61,
				Observaciones2: &obsTest3,
			},
		},
		{
			Index:       5,
			Description: "check data returned for sou_id 50 Microsoft Ofertas",
			Lead: Lead{
				SouID:         50,
				SouIDLeontel:  61,
				IsSmartCenter: false,
				Microsoft:     &Microsoft{},
			},
			ExpectedResult: LeadLeontel{
				LeaSource: 61,
			},
		},
		{
			Index:       6,
			Description: "check data returned for sou_id 51 Microsoft Ficha Producto",
			Lead: Lead{
				SouID:        51,
				SouIDLeontel: 61,
				Observations: &observations,
				Microsoft: &Microsoft{
					ProductType:        &t1,
					ProductName:        &t2,
					ProductID:          &t3,
					OriginalPrice:      &t4,
					Price:              &t5,
					Brand:              &t6,
					DiscountPercentage: &t7,
					DiscountCode:       &t1,
					ProcessorType:      &t2,
					DiskCapacity:       &t3,
					Graphics:           &t4,
					WirelessInterface:  &t5,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:      61,
				Observaciones2: &obsTest4,
			},
		},
		{
			Index:       7,
			Description: "check data returned for sou_id 54 R Cable Expansion",

			Lead: Lead{
				SouID:        54,
				SouIDLeontel: 63,
				RcableExp: &RcableExp{
					RespValues: &t1,
					Location:   &t2,
					Answer:     &t3,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     63,
				Observaciones: &obsTest5,
			},
		},
		{
			Index:       8,
			Description: "check data returned for sou_id 64 R Cable End To End (Kinkon) C2C",
			Lead: Lead{
				SouID:         64,
				SouIDLeontel:  73,
				LeaPhone:      &t1,
				IsSmartCenter: false,
				Kinkon: &Kinkon{
					CovData:     &CovData{},
					Portability: &Portability{},
					HolderData:  &HolderData{},
					BillingInfo: &BillingInfo{},
				},
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     73,
				Telefono:      &t1,
				Observaciones: &t8,
			},
		},
		{
			Index:       9,
			Description: "check data returned for sou_id 64 R Cable End To End (Kinkon) Coverture check",
			Lead: Lead{
				SouID:            64,
				SouIDLeontel:     73,
				LeatypeIDLeontel: 26,
				Kinkon: &Kinkon{
					Coverture: &t1,
					CovData: &CovData{
						State:    &t2,
						Town:     &t3,
						Street:   &t4,
						Number:   &t5,
						Floor:    &t6,
						CovPhone: &t7,
					},
					Portability: &Portability{},
					HolderData:  &HolderData{},
					BillingInfo: &BillingInfo{},
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     73,
				LeaType:       26,
				Observaciones: &obsTest6,
			},
		},
		{
			Index:       10,
			Description: "check data returned for sou_id 66 Euskaltel End To End (Kinkon) Hiring process",
			Lead: Lead{
				SouID:            66,
				SouIDLeontel:     75,
				LeatypeIDLeontel: 30,
				Kinkon: &Kinkon{
					Coverture: &t1,
					CovData: &CovData{
						State:    &t2,
						Town:     &t3,
						Street:   &t4,
						Number:   &t5,
						Floor:    &t6,
						CovPhone: &t7,
					},
					Portability: &Portability{
						Phone:                &t2,
						PhoneProvider:        &t3,
						MobilePhone:          &t4,
						MobilePhoneProvider:  &t5,
						MobilePhone2:         &t6,
						MobilePhoneProvider2: &t7,
					},
					HolderData: &HolderData{
						Name:         &t2,
						Surname:      &t3,
						Idnumber:     &t4,
						Mail:         &t5,
						ContactPhone: &t6,
					},
					BillingInfo: &BillingInfo{
						AccountHolder: &t1,
						AccountNumber: &t2,
					},
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     75,
				LeaType:       30,
				Observaciones: &obsTest8,
			},
		},
	}

	for _, test := range tests {
		leontel := test.Lead.LeadToLeontel()
		assert.Equal(test.ExpectedResult, leontel)
	}
}

func TestSenLeadToLeontel(t *testing.T) {
	assert := assert.New(t)

	phone := HelperRandstring(9)

	tests := []struct {
		Description      string
		TypeRequest      string
		StatusCode       int
		LeadLeontel      []LeadLeontel
		ExpectedResult   bool
		ExpectedResponse LeontelResp
	}{
		{
			Description: "when HandleFunction receive a POST request with no data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			LeadLeontel: []LeadLeontel{},
			ExpectedResponse: LeontelResp{
				Success: false,
				Error:   "Origen del lead erroneo use getSources()",
			},
		},
		{
			Description: "when HandleFunction receive a POST request with data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			LeadLeontel: []LeadLeontel{
				{
					LeaSource: 23,
					LeaType:   2,
					Telefono:  &phone,
				},
			},
			ExpectedResponse: LeontelResp{
				Success: true,
				ID:      9999,
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

		body, err := json.Marshal(test.LeadLeontel)
		assert.NoError(err)

		req, err := http.NewRequest(test.TypeRequest, ts.URL, bytes.NewBuffer(body))
		assert.NoError(err)

		http := &http.Client{}
		resp, err := http.Do(req)
		assert.NoError(err)

		assert.Equal(resp.StatusCode, test.StatusCode)

		response := LeontelResp{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(err)

		log.Println(response)

		assert.Equal(response.Success, test.ExpectedResponse.Success)
		assert.Equal(response.ID, test.ExpectedResponse.ID)
	}
}
