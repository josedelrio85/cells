package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// Evolution represents data structure needed in SC
type Evolution struct {
	Properties     Properties     `json:"propiedades"`
	AdditionalData AdditionalData `json:"datosAdicionales"`
	Localizators   Localizators   `json:"localizadores"`
}

// Properties bla
type Properties struct {
	// SubjectID       string `json:"idsujeto"`
	OriginalID      string `json:"idoriginal"`
	CampaingID      int64  `json:"idcampanya"`
	Name            string `json:"nombre"`
	Surname         string `json:"apellido"`
	Surname2        string `json:"apellidO2"`
	Phone           string `json:"telefono"`
	Phone2          string `json:"telefono2"`
	PhoneWork       string `json:"telefonoTrabajo"`
	MobilePhone     string `json:"movil"`
	MobilePhone2    string `json:"movil2"`
	Address         string `json:"direccion"`
	PostalCode      string `json:"codigO_POSTAL"`
	Town            string `json:"poblacion"`
	State           string `json:"provincia"`
	Country         string `json:"pais"`
	Fax             string `json:"fax"`
	Email           string `json:"email"`
	Email2          string `json:"emaiL2"`
	BirthDate       string `json:"fechA_NACIMIENTO"`
	SignupDate      string `json:"fechA_ALTA"`
	LanguageID      int64  `json:"iD_IDIOMA"`
	Observations    string `json:"observaciones"`
	LocatableSince  int64  `json:"localizablE_DESDE"`
	LocatableFrom   int64  `json:"localizablE_HASTA"`
	DNI             string `json:"sDNI"`
	FullName        string `json:"sNombre_Completo"`
	Company         string `json:"sEmpresa"`
	Sex             string `json:"cSexo"`
	Text1           string `json:"textO1"`
	Text2           string `json:"textO2"`
	Text3           string `json:"textO3"`
	FavSource       int64  `json:"nCanalPreferencial"`
	Num1            int64  `json:"nuM1"`
	Num2            int64  `json:"nuM2"`
	Num3            int64  `json:"nuM3"`
	SegmentAttibute string `json:"atributo_Segmento"`
	Priority        int64  `json:"prioridad"`
	NextContact     string `json:"tProximo_Contacto"`
	NState          string `json:"nEstado"`
	// Skill           string `json:"atributo_Skill"`
	// NList           string `json:"nLista"`
}

// AdditionalData bla
type AdditionalData struct {
	AddProp1 string `json:"additionalProp1"`
	AddProp2 string `json:"additionalProp2"`
	AddProp3 string `json:"additionalProp3"`
}

// Localizators bla
type Localizators struct {
	AddProp1 int64 `json:"additionalProp1"`
	AddProp2 int64 `json:"additionalProp2"`
	AddProp3 int64 `json:"additionalProp3"`
}

// EvolutionResp represents response from Evolution API
type EvolutionResp struct {
	Success bool
}

// LeadToEvolution maps leads values to Evolution struct
func (lead Lead) LeadToEvolution() Evolution {
	// TODO this code should be removed when Evolution API
	// allows to send empty data in this fields
	f1 := "2020-07-08"
	f2 := "2020-07-12"

	evolution := Evolution{
		Properties: Properties{
			OriginalID:  *lead.LeaPhone,
			CampaingID:  lead.SouIDEvolution,
			Phone:       *lead.LeaPhone,
			Name:        checkPointerValue(lead.LeaName),
			Email:       checkPointerValue(lead.LeaMail),
			DNI:         checkPointerValue(lead.LeaDNI),
			BirthDate:   f1,
			SignupDate:  f1,
			NextContact: f2,
		},
		AdditionalData: AdditionalData{
			AddProp1: checkPointerValue(lead.LeaURL),
			AddProp2: *lead.LeaIP,
			AddProp3: checkPointerValue(lead.GaClientID),
		},
		Localizators: Localizators{
			AddProp1: int64(lead.ID),
		},
	}

	switch souid := lead.SouID; souid {
	case 79:
		args := []*string{}
		args = append(args, lead.Observations)

		if lead.Virgin != nil {
			q1 := "Optin"
			q2 := "Código postal"
			q3 := "Edad"
			q4 := "Apellidos"
			q5 := "External ID"
			q6 := "Datos al mes"
			q7 := "¿Tienes actualmente ADSL/Fibra?"
			q8 := "Cuando lo vayas a contratar"

			args = append(args, &q1, lead.Virgin.Optin)
			args = append(args, &q2, lead.Virgin.PostalCode)
			args = append(args, &q3, lead.Virgin.Age)
			args = append(args, &q4, lead.Virgin.Surname)
			args = append(args, &q5, lead.Virgin.ExternalID)
			args = append(args, &q6, lead.Virgin.DataMonth)
			args = append(args, &q7, lead.Virgin.HaveDSL)
			args = append(args, &q8, lead.Virgin.WhenHiring)
		}
		observations := concatPointerStrs(args...)
		evolution.Properties.Observations = observations
	default:
	}
	return evolution
}

// Active is an implementation of Active method from Scable interface
func (e Evolution) Active(lead Lead) bool {
	switch lead.SouID {
	// virgin
	case 79:
		log.Printf("souid %d Evolution active", lead.SouID)
		return true
	default:
		return false
	}
}

// Send is an implementation of Send method from Scable interface
func (e Evolution) Send(lead Lead) ScResponse {
	leadevolution := lead.LeadToEvolution()
	bytevalues, err := json.Marshal(leadevolution)
	// fmt.Println(string(bytevalues))
	if err != nil {
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}

	endpoint, ok := os.LookupEnv("EVOLUTION_ENDPOINT")
	if !ok {
		err := errors.New("unable to load Lead Evolution URL endpoint")
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(bytevalues))
	if err != nil {
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}
	req.Header.Add("accept", "text/plain")
	req.Header.Add("Content-Type", "application/json")
	// TODO set env_var
	userid := "bySideCar"
	password := "47Qh5qQy5JsRZQCzAiRi"
	req.SetBasicAuth(userid, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	txtresp := string(data)
	evolutionresp := EvolutionResp{
		Success: false,
	}
	if txtresp == "true" {
		evolutionresp.Success = true
	}

	err = nil
	status := http.StatusOK
	if !evolutionresp.Success {
		err = errors.New(fmt.Sprintf("Evolution response %t", evolutionresp.Success))
		status = http.StatusUnprocessableEntity
	}

	return ScResponse{
		Success:    evolutionresp.Success,
		StatusCode: status,
		// TODO should be SubjectID field, not implemented by provider
		ID:    0,
		Error: err,
	}
}

func checkPointerValue(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}
