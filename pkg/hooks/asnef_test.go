package leads

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/bysidecar/leads/pkg/model"
)

func TestActive(t *testing.T) {
	assert := assert.New(t)

	var asnef Asnef

	tests := []struct {
		Description string
		Lead        model.Lead
		Active      bool
	}{
		{
			Description: "when Asnef successfully is activated",
			Lead: model.Lead{
				SouID:         9,
				IsSmartCenter: true,
			},
			Active: true,
		},
		{
			Description: "when Asnef successfully is activated",
			Lead: model.Lead{
				SouID:         58,
				IsSmartCenter: true,
			},
			Active: true,
		},
		{
			Description: "when Asnef is not activated",
			Lead: model.Lead{
				SouID:         99,
				IsSmartCenter: true,
			},
			Active: false,
		},
		{
			Description: "when Asnef is not activated",
			Lead: model.Lead{
				SouID:         0,
				IsSmartCenter: true,
			},
			Active: false,
		},
		{
			Description: "when IsSmartCenter is false",
			Lead: model.Lead{
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

	database := helperDb()
	database.Open()
	defer database.Close()

	lead := model.Lead{
		SouID: 9,
	}
	candidates := GetCandidates(lead)

	tests := []struct {
		Description string
		Lead        model.Lead
		Response    HookResponse
	}{
		{
			Description: "when Asnef successfully checks and gives ok to the lead",
			Lead: model.Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: model.Creditea{
					Asnef:     false,
					Yacliente: false,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		},
		{
			Description: "when Asnef checks is clicked. Client activates the limitation",
			Lead: model.Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: model.Creditea{
					Cantidadsolicitada: &cantidad,
					Asnef:              true,
					Yacliente:          false,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		},
		{
			Description: "when Yacliente check is clicked. Client activates the limitation",
			Lead: model.Lead{
				SouID:         9,
				LeaPhone:      &phone,
				LeaDNI:        &dni,
				IsSmartCenter: true,
				Creditea: model.Creditea{
					Cantidadsolicitada: &cantidad,
					Asnef:              false,
					Yacliente:          true,
				},
			},
			Response: HookResponse{
				Err:        nil,
				StatusCode: http.StatusOK,
			},
		},
		{
			Description: "when Asnef validations is not passed",
			Lead: model.Lead{
				SouID:         9,
				LeaPhone:      &candidates[0].Telefono,
				LeaDNI:        &candidates[0].DNI,
				IsSmartCenter: false,
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
			Lead        model.Lead
			Response    HookResponse
		}{
			Description: "when Asnef pre validation is not passed",
			Lead: model.Lead{
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
			log.Println(test.Description)
			log.Println("-----------------")
			response := asnef.Perform(database.DB, &test.Lead)

			assert.Equal(test.Response, response)
			assert.Equal(test.Response.StatusCode, response.StatusCode)
			assert.Equal(test.Response.Err, response.Err)
		})
	}
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

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}
