package leads

import (
	"log"
	"net/http"
	"strconv"
	"testing"

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

func TestPerform(t *testing.T) {
	assert := assert.New(t)

	var asnef Asnef
	phone := "665932356"
	dni := "21528205K"
	cantidad := "1000"

	database := helperDbAsnef()
	database.Open()
	defer database.Close()

	lead := Lead{
		SouID: 9,
	}
	candidates := GetCandidates(lead)

	tests := []struct {
		Description string
		Lead        Lead
		Response    HookResponse
	}{
		{
			Description: "when ASNEF successfully checks and gives ok to the lead",
			Lead: Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: &Creditea{
					ASNEF:         false,
					AlreadyClient: false,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		},
		{
			Description: "when ASNEF checks is clicked. Client activates the limitation",
			Lead: Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: &Creditea{
					RequestedAmount: &cantidad,
					ASNEF:           true,
					AlreadyClient:   false,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		},
		{
			Description: "when AlreadyClient checks is clicked. Client activates the limitation",
			Lead: Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: &Creditea{
					RequestedAmount: &cantidad,
					ASNEF:           false,
					AlreadyClient:   true,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		},
		{
			Description: "when ASNEF validations is not passed",
			Lead: Lead{
				SouID:         9,
				LeaPhone:      &candidates[0].Telefono,
				LeaDNI:        &candidates[0].DNI,
				IsSmartCenter: false,
				Creditea: &Creditea{
					ASNEF:         true,
					AlreadyClient: true,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		},
	}

	candidatespre := GetCandidatesPreasnef(database.DB)
	if len(candidatespre) > 0 {
		newtest := struct {
			Description string
			Lead        Lead
			Response    HookResponse
		}{
			Description: "when Asnef pre validation is not passed",
			Lead: Lead{
				SouID:         9,
				LeaPhone:      candidatespre[0].LeaPhone,
				LeaDNI:        candidatespre[0].LeaDNI,
				IsSmartCenter: false,
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		}
		tests = append(tests, newtest)
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			cont := Handler{
				Storer: &database,
				Lead:   test.Lead,
			}
			response := asnef.Perform(&cont)

			assert.Equal(test.Response, response)
			assert.Equal(test.Response.StatusCode, response.StatusCode)
			assert.Equal(test.Response.Err, response.Err)
		})
	}
}

func helperDbAsnef() Database {

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