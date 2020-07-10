package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActiveEvolution(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Index          int
		Description    string
		Lead           Lead
		Evolution      Evolution
		ExpectedResult bool
	}{
		{
			Index:       1,
			Description: "when a lead not belongs to Evolution (R Cable End to End)",
			Lead: Lead{
				SouID: 64,
			},
			Evolution:      Evolution{},
			ExpectedResult: false,
		},
		{
			Index:       2,
			Description: "when a lead belongs to Evolution (Virgin)",
			Lead: Lead{
				SouID: 79,
			},
			Evolution:      Evolution{},
			ExpectedResult: true,
		},
	}

	for _, test := range tests {
		result := test.Evolution.Active(test.Lead)
		assert.Equal(result, test.ExpectedResult)
	}
}

func TestLeadToEvolution(t *testing.T) {
	assert := assert.New(t)

	phone := "676767676"
	ip := "::1"
	sep := "--"
	empty := ""

	virginobs := []string{
		fmt.Sprintf("Optin %s %s %s", sep, empty, sep),
		fmt.Sprintf(" Código postal %s %s %s", sep, empty, sep),
		fmt.Sprintf(" Edad %s %s %s", sep, empty, sep),
		fmt.Sprintf(" Apellidos %s %s %s", sep, empty, sep),
		fmt.Sprintf(" External ID %s %s %s", sep, empty, sep),
		fmt.Sprintf(" Datos al mes %s %s %s", sep, empty, sep),
		fmt.Sprintf(" ¿Tienes actualmente ADSL/Fibra? %s %s %s", sep, empty, sep),
		fmt.Sprintf(" Cuando lo vayas a contratar %s %s", sep, empty),
	}

	tests := []struct {
		Index          int
		Description    string
		Lead           Lead
		ExpectedResult Evolution
	}{
		{
			Index:       1,
			Description: "check data returned for sou_id 79 Virgin campaign",
			Lead: Lead{
				SouID:          79,
				SouIDEvolution: 100000006,
				LeaPhone:       &phone,
				LeaIP:          &ip,
				IsSmartCenter:  false,
				Virgin: &Virgin{
					Optin:      &empty,
					PostalCode: &empty,
					Age:        &empty,
					Surname:    &empty,
					ExternalID: &empty,
					DataMonth:  &empty,
					HaveDSL:    &empty,
					WhenHiring: &empty,
				},
			},
			ExpectedResult: Evolution{
				Properties: Properties{
					CampaignID:   100000006,
					OriginalID:   phone,
					Phone:        phone,
					Name:         empty,
					Email:        empty,
					DNI:          empty,
					Observations: strings.Join(virginobs, ""),
				},
				AdditionalData: AdditionalData{
					AddProp1: empty,
					AddProp2: ip,
					AddProp3: empty,
				},
				Localizators: Localizators{
					AddProp1: 0,
				},
			},
		},
	}

	for _, test := range tests {
		result := test.Lead.LeadToEvolution()
		assert.Equal(test.ExpectedResult, result)
	}
}

func TestSenLeadToEvolution(t *testing.T) {
	assert := assert.New(t)

	phone := HelperRandstring(9)
	empty := ""

	tests := []struct {
		Description      string
		TypeRequest      string
		StatusCode       int
		LeadEvolution    Evolution
		ExpectedResult   bool
		ShouldResponse   string
		ExpectedResponse EvolutionResp
	}{
		{
			Description:    "when HandleFunction receive a POST request with no data",
			TypeRequest:    http.MethodPost,
			StatusCode:     http.StatusInternalServerError,
			LeadEvolution:  Evolution{},
			ShouldResponse: "false",
			ExpectedResponse: EvolutionResp{
				Success: false,
			},
		},
		{
			Description: "when HandleFunction receive a POST request with onlye a part of needed data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			LeadEvolution: Evolution{
				Properties: Properties{
					OriginalID: phone,
					CampaignID: 1,
					Name:       empty,
				},
			},
			ShouldResponse: "false",
			ExpectedResponse: EvolutionResp{
				Success: false,
			},
		},
		{
			Description: "when HandleFunction receive a POST request with all needed data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			LeadEvolution: Evolution{
				Properties: Properties{
					OriginalID:      phone,
					CampaignID:      1,
					Name:            empty,
					Surname:         empty,
					Surname2:        empty,
					Phone:           phone,
					Phone2:          empty,
					PhoneWork:       empty,
					MobilePhone:     empty,
					MobilePhone2:    empty,
					Address:         empty,
					PostalCode:      empty,
					Town:            empty,
					State:           empty,
					Country:         empty,
					Fax:             empty,
					Email:           empty,
					Email2:          empty,
					BirthDate:       empty,
					SignupDate:      empty,
					LanguageID:      0,
					Observations:    empty,
					LocatableSince:  0,
					LocatableFrom:   0,
					DNI:             empty,
					FullName:        empty,
					Company:         empty,
					Sex:             empty,
					Text1:           empty,
					Text2:           empty,
					Text3:           empty,
					FavSource:       0,
					Num1:            0,
					Num2:            0,
					Num3:            0,
					SegmentAttibute: empty,
					Priority:        0,
					NextContact:     empty,
					NState:          empty,
				},
				AdditionalData: AdditionalData{
					AddProp1: empty,
					AddProp2: empty,
					AddProp3: empty,
				},
				Localizators: Localizators{
					AddProp1: 0,
					AddProp2: 0,
					AddProp3: 0,
				},
			},
			ShouldResponse: "true",
			ExpectedResponse: EvolutionResp{
				Success: true,
			},
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(test.StatusCode)
			resptxt := []byte(test.ShouldResponse)
			res.Write(resptxt)
		}))
		defer func() { ts.Close() }()

		body, err := json.Marshal(test.LeadEvolution)
		assert.NoError(err)

		req, err := http.NewRequest(test.TypeRequest, ts.URL, bytes.NewBuffer(body))
		assert.NoError(err)
		req.Header.Add("accept", "text/plain")
		req.Header.Add("Content-Type", "application/json")

		userid, ok := os.LookupEnv("EVOLUTION_AUTH_USER")
		assert.True(ok)
		password, ok := os.LookupEnv("EVOLUTION_AUTH_PASS")
		assert.True(ok)
		req.SetBasicAuth(userid, password)

		http := &http.Client{}
		resp, err := http.Do(req)
		assert.NoError(err)

		assert.Equal(resp.StatusCode, test.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		assert.NoError(err)
		txtresp := string(data)

		evolutionresp := EvolutionResp{
			Success: false,
		}
		if txtresp == "true" {
			evolutionresp.Success = true
		}

		assert.Equal(evolutionresp.Success, test.ExpectedResponse.Success)
	}
}
