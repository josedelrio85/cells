package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dbInstance Database

func TestMain(m *testing.M) {
	dbInstance = helperDb()

	code := m.Run()

	setDownDb()

	os.Exit(code)
}

func TestHandlerFunction(t *testing.T) {
	assert := assert.New(t)

	phoneTest := "666666666"
	ipTest := "127.0.0.1"

	database := helperDb()

	tests := []struct {
		Description    string
		Storer         Storer
		TypeRequest    string
		StatusCode     int
		Lead           Lead
		ExpectedResult bool
	}{
		{
			Description:    "when HandleFunction receive a POST request with no data",
			TypeRequest:    http.MethodPost,
			StatusCode:     http.StatusInternalServerError,
			Storer:         nil,
			Lead:           Lead{},
			ExpectedResult: false,
		},
		{
			Description: "when HandleFunction receive a POST request without sou_id value",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Storer:      nil,
			Lead: Lead{
				LeatypeID:     1,
				LeaPhone:      &phoneTest,
				LeaIP:         &ipTest,
				IsSmartCenter: false,
			},
			ExpectedResult: false,
		},
		{
			Description: "when HandleFunction receive a POST request without lea_type value",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			Storer:      &database,
			Lead: Lead{
				SouID:         15,
				LeaPhone:      &phoneTest,
				LeaIP:         &ipTest,
				IsSmartCenter: false,
			},
			ExpectedResult: true,
		},
	}

	for _, test := range tests {
		ch := Handler{
			Storer: test.Storer,
		}

		if test.Storer != nil {
			err := test.Storer.Open()
			defer test.Storer.Close()
			assert.NoError(err)
		}

		ts := httptest.NewServer(ch.HandleFunction())
		defer ts.Close()

		body, err := json.Marshal(test.Lead)
		if err != nil {
			t.Errorf("error marshalling test json: Err: %v", err)
			return
		}

		req, err := http.NewRequest(test.TypeRequest, ts.URL, bytes.NewBuffer(body))
		if err != nil {
			t.Errorf("error creating the test Request: err %v", err)
			return
		}

		http := &http.Client{}
		resp, err := http.Do(req)
		if err != nil {
			t.Errorf("error sending test request: Err %v", err)
			return
		}

		assert.Equal(resp.StatusCode, test.StatusCode)

		response := ResponseAPI{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Errorf("error decoding response. Err %v", err)
			return
		}

		assert.Equal(response.Success, test.ExpectedResult)
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

	database := helperDb()
	err := database.Open()
	defer database.Close()

	assert.NoError(err)

	tests := []struct {
		Index          int
		Description    string
		TypeRequest    string
		StatusCode     int
		Lead           Lead
		ExpectedResult LeadLeontel
	}{
		{
			Index:       1,
			Description: "check data returned for sou_id 9 Creditea EndToEnd",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:     9,
				LeatypeID: 1,
				LeaDNI:    &t1,
				Creditea: &Creditea{
					RequestedAmount: &t2,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     13,
				LeaType:       2,
				Dninie:        &t1,
				Observaciones: &obsTest,
			},
		},
		{
			Index:       2,
			Description: "check data returned for sou_id 11 Creditea Rastreator",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:     11,
				LeatypeID: 1,
				LeaDNI:    &t1,
				Creditea: &Creditea{
					RequestedAmount: &t2,
					NetIncome:       &t3,
					ContractType:    &t4,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     15,
				LeaType:       2,
				Dninie:        &t1,
				Observaciones: &obsTest2,
			},
		},
		{
			Index:       3,
			Description: "check data returned for sou_id 46-49 Microsoft Hazelcambio + Recomendador",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:        46,
				LeatypeID:    1,
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
				LeaType:        2,
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
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:        48,
				LeatypeID:    1,
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
				LeaType:        2,
				Observaciones2: &obsTest3,
			},
		},
		{
			Index:       5,
			Description: "check data returned for sou_id 50 Microsoft Ofertas",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:         50,
				LeatypeID:     1,
				IsSmartCenter: false,
				Microsoft:     &Microsoft{},
			},
			ExpectedResult: LeadLeontel{
				LeaSource: 61,
				LeaType:   2,
			},
		},
		{
			Index:       6,
			Description: "check data returned for sou_id 51 Microsoft Ficha Producto",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:        51,
				LeatypeID:    1,
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
				LeaType:        2,
				Observaciones2: &obsTest4,
			},
		},
		{
			Index:       7,
			Description: "check data returned for sou_id 54 R Cable Expansion",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:     54,
				LeatypeID: 1,
				RcableExp: &RcableExp{
					RespValues: &t1,
					Location:   &t2,
					Answer:     &t3,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     63,
				LeaType:       2,
				Observaciones: &obsTest5,
			},
		},
		{
			Index:       8,
			Description: "check data returned for sou_id 64 R Cable End To End (Kinkon) C2C",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:         64,
				LeatypeID:     1,
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
				LeaType:       2,
				Telefono:      &t1,
				Observaciones: &t8,
			},
		},
		{
			Index:       9,
			Description: "check data returned for sou_id 64 R Cable End To End (Kinkon) Coverture check",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:     64,
				LeatypeID: 24,
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
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: Lead{
				SouID:     66,
				LeatypeID: 27,
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
		ch := Handler{
			Lead: test.Lead,
		}

		err := ch.Lead.GetLeontelValues(database.DB)
		assert.NoError(err)

		leontel := ch.Lead.LeadToLeontel()
		assert.Equal(test.ExpectedResult, leontel)
	}
}

func TestGetLeontelValues(t *testing.T) {
	assert := assert.New(t)
	database := helperDb()
	err := database.Open()
	defer database.Close()

	assert.NoError(err)

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
				SouIDLeontel:     13,
				LeatypeIDLeontel: 2,
			},
		},
		{
			Description: "EVO BANCO 3 => 4 | INACTIVIDAD 2 => 3",
			Lead: Lead{
				SouID:     3,
				LeatypeID: 2,
			},
			ExpectedResult: Lead{
				SouIDLeontel:     4,
				LeatypeIDLeontel: 3,
			},
		},
		{
			Description: "R CABLE EXPANSION END TO END 54 => 63 | FDH 8 => 12",
			Lead: Lead{
				SouID:     54,
				LeatypeID: 8,
			},
			ExpectedResult: Lead{
				SouIDLeontel:     63,
				LeatypeIDLeontel: 12,
			},
		},
	}

	for _, test := range tests {
		err := test.Lead.GetLeontelValues(database.DB)
		assert.NoError(err)

		assert.Equal(test.ExpectedResult.SouIDLeontel, test.Lead.SouIDLeontel)
		assert.Equal(test.ExpectedResult.LeatypeIDLeontel, test.Lead.LeatypeIDLeontel)
	}
}

type LeontelRespTest struct {
	Success bool  `json:"success"`
	ID      int64 `json:"id"`
}

func TestSendLeadToLeontel(t *testing.T) {
	assert := assert.New(t)

	t3 := "c"
	t4 := "D"
	t5 := "000000000"
	t6 := "99997896Z"

	database := helperDb()

	tests := []struct {
		Description    string
		Lead           Lead
		ExpectedResult LeontelRespTest
		Storer         Storer
	}{
		{
			Description: "When send a valid lead to Leontel",
			Lead: Lead{
				SouID:            15,
				SouIDLeontel:     23,
				LeatypeID:        1,
				LeatypeIDLeontel: 2,
				LeaURL:           &t3,
				LeaIP:            &t4,
				LeaPhone:         &t5,
				LeaDNI:           &t6,
			},
			ExpectedResult: LeontelRespTest{
				Success: true,
			},
			Storer: &database,
		},
	}

	for _, test := range tests {

		if test.Storer != nil {
			err := test.Storer.Open()
			defer test.Storer.Close()
			assert.NoError(err)
		}

		result, err := test.Lead.SendLeadToLeontel()

		leontelID := strconv.FormatInt(result.ID, 10)
		test.Lead.LeaSmartcenterID = leontelID

		test.Storer.Insert(&test.Lead)

		cond := fmt.Sprintf("ID=%d", test.Lead.ID)
		fields := []string{"LeaSmartcenterID"}
		test.Storer.Update(&test.Lead, cond, fields)

		assert.NoError(err)
		assert.Equal(test.ExpectedResult.Success, result.Success)
		assert.NotNil(test.Lead.LeaSmartcenterID)

		test.Storer.Instance().Delete(test.Lead)
	}
}

func TestOpenDb(t *testing.T) {
	assert := assert.New(t)

	err := dbInstance.Open()

	assert.NoError(err)
}

func helperDb() Database {

	port := GetSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string the Redshift's port %s, Err: %s", port, err)
	}

	database := Database{
		Host:      GetSetting("DB_HOST"),
		Port:      portInt,
		User:      GetSetting("DB_USER"),
		Pass:      GetSetting("DB_PASS"),
		Dbname:    GetSetting("DB_NAME"),
		Charset:   "utf8",
		ParseTime: "True",
		Loc:       "Local",
	}
	return database
}

func setDownDb() {

	if err := dbInstance.Open(); err != nil {
		log.Printf("error opening database connection. err: %s", err)
	}
	defer dbInstance.DB.Close()
}

// GetSetting reads an ENV VAR setting, it does crash the service if with an
// error message if any setting is not found.
//
// - setting: The setting (ENV VAR) to read.
//
// Returns the setting value.
func GetSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}

// HelperRandstring is a helper function to generate random strings
func HelperRandstring(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
