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

	model "github.com/bysidecar/leads/pkg/model"

	"github.com/stretchr/testify/assert"
)

var dbInstance model.Database

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
		Storer         model.Storer
		TypeRequest    string
		StatusCode     int
		Lead           model.Lead
		ExpectedResult bool
	}{
		{
			Description:    "when HandleFunction receive a POST request with no data",
			TypeRequest:    http.MethodPost,
			StatusCode:     http.StatusInternalServerError,
			Storer:         nil,
			Lead:           model.Lead{},
			ExpectedResult: false,
		},
		{
			Description: "when HandleFunction receive a POST request without sou_id value",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Storer:      nil,
			Lead: model.Lead{
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
			Lead: model.Lead{
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
	observations := "H"

	obsTest := "A -- B"
	obsTest2 := "A -- C -- D -- B"
	obsTest3 := "A -- B -- C -- D -- E -- F"
	obsTest4 := "A -- B -- C -- D -- E -- F -- G -- A -- B -- C -- D -- E"
	obsTest5 := "A -- B -- C"

	database := helperDb()
	err := database.Open()
	defer database.Close()

	assert.NoError(err)

	tests := []struct {
		Description    string
		TypeRequest    string
		StatusCode     int
		Lead           model.Lead
		ExpectedResult model.LeadLeontel
	}{
		{
			Description: "check data returned for sou_id 9 Creditea EndToEnd",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: model.Lead{
				SouID:     9,
				LeatypeID: 1,
				LeaDNI:    &t1,
				Creditea: &model.Creditea{
					RequestedAmount: &t2,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: model.LeadLeontel{
				LeaSource:     13,
				LeaType:       2,
				Dninie:        &t1,
				Observaciones: &obsTest,
			},
		},
		{
			Description: "check data returned for sou_id 11 Creditea Rastreator",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: model.Lead{
				SouID:     11,
				LeatypeID: 1,
				LeaDNI:    &t1,
				Creditea: &model.Creditea{
					RequestedAmount: &t2,
					NetIncome:       &t3,
					ContractType:    &t4,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: model.LeadLeontel{
				LeaSource:     15,
				LeaType:       2,
				Dninie:        &t1,
				Observaciones: &obsTest2,
			},
		},
		{
			Description: "check data returned for sou_id 46-49 Microsoft Hazelcambio + Recomendador",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: model.Lead{
				SouID:        46,
				LeatypeID:    1,
				Observations: &observations,
				Microsoft: &model.Microsoft{
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
			ExpectedResult: model.LeadLeontel{
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
			Description: "check data returned for sou_id 48 Microsoft Calculadora",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: model.Lead{
				SouID:        48,
				LeatypeID:    1,
				Observations: &observations,
				Microsoft: &model.Microsoft{
					DevicesAverageAge:      &t1,
					DevicesOperatingSystem: &t2,
					DevicesHangFrequency:   &t3,
					DevicesNumber:          &t4,
					DevicesLastYearRepairs: &t5,
					DevicesStartupTime:     &t6,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: model.LeadLeontel{
				LeaSource:      61,
				LeaType:        2,
				Observaciones2: &obsTest3,
			},
		},
		{
			Description: "check data returned for sou_id 50 Microsoft Ofertas",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: model.Lead{
				SouID:         50,
				LeatypeID:     1,
				IsSmartCenter: false,
				Microsoft:     &model.Microsoft{},
			},
			ExpectedResult: model.LeadLeontel{
				LeaSource: 61,
				LeaType:   2,
			},
		},
		{
			Description: "check data returned for sou_id 51 Microsoft Ficha Producto",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: model.Lead{
				SouID:        51,
				LeatypeID:    1,
				Observations: &observations,
				Microsoft: &model.Microsoft{
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
			ExpectedResult: model.LeadLeontel{
				LeaSource:      61,
				LeaType:        2,
				Observaciones2: &obsTest4,
			},
		},
		{
			Description: "check data returned for sou_id 54 R Cable Expansion",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			Lead: model.Lead{
				SouID:     54,
				LeatypeID: 1,
				RcableExp: &model.RcableExp{
					RespValues: &t1,
					Location:   &t2,
					Answer:     &t3,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: model.LeadLeontel{
				LeaSource:     63,
				LeaType:       2,
				Observaciones: &obsTest5,
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
		Lead           model.Lead
		ExpectedResult model.Lead
	}{
		{
			Description: "CREDITEA END TO END	9 => 13 | C2C 1 => 2",
			Lead: model.Lead{
				SouID:     9,
				LeatypeID: 1,
			},
			ExpectedResult: model.Lead{
				SouIDLeontel:     13,
				LeatypeIDLeontel: 2,
			},
		},
		{
			Description: "EVO BANCO 3 => 4 | INACTIVIDAD 2 => 3",
			Lead: model.Lead{
				SouID:     3,
				LeatypeID: 2,
			},
			ExpectedResult: model.Lead{
				SouIDLeontel:     4,
				LeatypeIDLeontel: 3,
			},
		},
		{
			Description: "R CABLE EXPANSION END TO END 54 => 63 | FDH 8 => 12",
			Lead: model.Lead{
				SouID:     54,
				LeatypeID: 8,
			},
			ExpectedResult: model.Lead{
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
		Lead           model.Lead
		ExpectedResult LeontelRespTest
		Storer         model.Storer
	}{
		{
			Description: "When send a valid lead to Leontel",
			Lead: model.Lead{
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
		test.Lead.LeaSmartcenterID = &leontelID

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

func helperDb() model.Database {

	port := getSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string the Redshift's port %s, Err: %s", port, err)
	}

	database := model.Database{
		Host:      getSetting("DB_HOST"),
		Port:      portInt,
		User:      getSetting("DB_USER"),
		Pass:      getSetting("DB_PASS"),
		Dbname:    getSetting("DB_NAME"),
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

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}

func helperRandstring(length int) string {
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
