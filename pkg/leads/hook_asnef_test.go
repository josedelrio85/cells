package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestActive(t *testing.T) {
	assert := assert.New(t)

	var asnef Asnef

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "when ASNEF successfully is activated",
			Lead: Lead{
				SouID:         9,
				IsSmartCenter: true,
			},
			Active: true,
		},
		{
			Description: "when ASNEF successfully is activated",
			Lead: Lead{
				SouID:         58,
				IsSmartCenter: true,
			},
			Active: true,
		},
		{
			Description: "when ASNEF is not activated",
			Lead: Lead{
				SouID:         99,
				IsSmartCenter: true,
			},
			Active: false,
		},
		{
			Description: "when ASNEF is not activated",
			Lead: Lead{
				SouID:         0,
				IsSmartCenter: true,
			},
			Active: false,
		},
		{
			Description: "when IsSmartCenter is false",
			Lead: Lead{
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

func TestAsnefValues(t *testing.T) {
	assert := assert.New(t)

	var asnef Asnef
	fakedb := FakeDb{
		OpenFunc:     func() error { return nil },
		CloseFunc:    func() error { return nil },
		UpdateFunc:   func(a interface{}, wCond string, wFields []string) error { return nil },
		InsertFunc:   func(lead interface{}) error { return nil },
		InstanceFunc: func() *gorm.DB { return nil },
	}

	type ExpectedResponse struct {
		SCState  bool
		Response HookResponse
	}

	tests := []struct {
		Description      string
		Lead             Lead
		ExpectedResponse ExpectedResponse
	}{
		{
			Description: "when ASNEF check is selected. Client activates the limitation",
			Lead: Lead{
				Creditea: &Creditea{
					ASNEF:         true,
					AlreadyClient: false,
				},
			},
			ExpectedResponse: ExpectedResponse{
				SCState: false,
				Response: HookResponse{
					StatusCode: http.StatusOK,
					Err:        nil,
				},
			},
		},
		{
			Description: "when AlreadyClient check is selected. Client activates the limitation",
			Lead: Lead{
				Creditea: &Creditea{
					ASNEF:         false,
					AlreadyClient: true,
				},
			},
			ExpectedResponse: ExpectedResponse{
				SCState: false,
				Response: HookResponse{
					StatusCode: http.StatusOK,
					Err:        nil,
				},
			},
		},
		{
			Description: "when AlreadyClient and ASNEF checks are selected. Client activates the limitation",
			Lead: Lead{
				Creditea: &Creditea{
					ASNEF:         true,
					AlreadyClient: true,
				},
			},
			ExpectedResponse: ExpectedResponse{
				SCState: false,
				Response: HookResponse{
					StatusCode: http.StatusOK,
					Err:        nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			cont := Handler{
				Storer: &fakedb,
				Lead:   test.Lead,
			}
			response := asnef.Perform(&cont)

			assert.Equal(cont.Lead.IsSmartCenter, test.ExpectedResponse.SCState)
			assert.Equal(response, test.ExpectedResponse.Response)
		})
	}
}

func TestResponseAsnef(t *testing.T) {
	assert := assert.New(t)

	phone := HelperRandstring(9)
	dni := HelperRandstring(9)

	tests := []struct {
		Description      string
		Lead             Lead
		ExpectedStatus   int
		ExpectedResponse Asnef
		ExpectedResult   Lead
	}{
		{
			Description: "when a lead has a positive asnef validation (exists on db)",
			Lead: Lead{
				SouID:    999,
				LeaDNI:   &dni,
				LeaPhone: &phone,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: Asnef{
				Result: true,
			},
			ExpectedResult: Lead{
				Creditea: &Creditea{
					ASNEF:         true,
					AlreadyClient: true,
				},
			},
		},
		{
			Description: "when a lead has a negative asnef validation (not exists on db)",
			Lead: Lead{
				SouID:    999,
				LeaDNI:   &dni,
				LeaPhone: &phone,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: Asnef{
				Result: false,
			},
			ExpectedResult: Lead{
				Creditea: &Creditea{
					ASNEF:         false,
					AlreadyClient: false,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				testresp, err := json.Marshal(test.ExpectedResponse)
				assert.NoError(err)

				res.WriteHeader(test.ExpectedStatus)
				res.Write(testresp)
			}))
			defer func() { testServer.Close() }()

			asnef := Asnef{}

			data := InputData{
				SouID: test.Lead.SouID,
				DNI:   *test.Lead.LeaDNI,
				Phone: *test.Lead.LeaPhone,
			}

			bytevalues, err := json.Marshal(data)
			assert.NoError(err)

			req, err := http.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(bytevalues))
			assert.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(err)

			rawdata, _ := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(rawdata, &asnef)
			assert.NoError(err)

			assert.Equal(asnef.Result, test.ExpectedResponse.Result)

			if test.ExpectedResponse.Result {
				assert.True(test.ExpectedResult.Creditea.ASNEF)
				assert.True(test.ExpectedResult.Creditea.AlreadyClient)
			} else {
				assert.False(test.ExpectedResult.Creditea.ASNEF)
				assert.False(test.ExpectedResult.Creditea.AlreadyClient)
			}
		})
	}
}

func TestGetSourceDescription(t *testing.T) {
	assert := assert.New(t)

	lead := Lead{
		SouID: 99,
	}

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_0")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_0")
	defer db.Close()

	rs := mock.NewRows([]string{"sou_id", "sou_description", "sou_idcrm"}).
		FromCSVString("5,hello world,3")

	mock.ExpectQuery("SELECT (.+)").
		WithArgs(lead.SouID).
		WillReturnRows(rs)

	var asnef Asnef

	str, err := asnef.GetSourceDescription(db, &lead)
	assert.NoError(err)
	assert.NotNil(str)
}

func TestGetSourcesFromDescription(t *testing.T) {
	assert := assert.New(t)

	description := "test"

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_1")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_1")
	defer db.Close()

	rs := mock.NewRows([]string{"sou_id", "sou_description", "sou_idcrm"}).
		FromCSVString("5,hello world,3")

	mock.ExpectQuery("SELECT (.+)").
		WithArgs(description).
		WillReturnRows(rs)

	var asnef Asnef

	arrstr, err := asnef.GetSourcesFromDescription(description, db)
	assert.NoError(err)
	assert.NotNil(arrstr)
}

func TestHasAsnef(t *testing.T) {
	assert := assert.New(t)

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_2")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_2")
	defer db.Close()

	phone := HelperRandstring(9)
	dni := HelperRandstring(9)

	lead := Lead{
		LeaPhone: &phone,
		LeaDNI:   &dni,
	}

	tests := []struct {
		Description string
		Rs          *sqlmock.Rows
		Expected    bool
	}{
		{
			Description: "When results are returned, return true",
			Rs: mock.NewRows([]string{"sou_id", "sou_description", "sou_idcrm"}).
				FromCSVString("5,hello world,3"),
			Expected: true,
		},
		{
			Description: "When NO results are returned, return false",
			Rs:          mock.NewRows([]string{"sou_id", "sou_description", "sou_idcrm"}),
			Expected:    false,
		},
	}

	oml := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	sources := "'1','2','3'"

	dnival := fmt.Sprintf("%s%s%s", "%", *lead.LeaDNI, "%")
	var asnef Asnef

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			mock.ExpectQuery("SELECT (.+)").
				WithArgs(oml, sources, 0, dnival, lead.LeaPhone, 1, 1).
				WillReturnRows(test.Rs)

			result, err := asnef.HasAsnef(sources, db, &lead)

			assert.NoError(err)
			assert.Equal(result, test.Expected)
		})
	}
}
