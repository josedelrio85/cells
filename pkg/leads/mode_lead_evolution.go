package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Evolution represents data structure needed in SC
type Evolution struct {
	Properties     Properties     `json:"propiedades"`
	AdditionalData AdditionalData `json:"datosAdicionales,omitempty"`
	Localizators   Localizators   `json:"localizadores,omitempty"`
}

// Properties bla
type Properties struct {
	SubjectID       string `json:"idsujeto,omitempty"`
	OriginalID      string `json:"idoriginal"`
	CampaingID      string `json:"idcampanya"`
	Name            string `json:"nombre,omitempty"`
	Surname         string `json:"apellido,omitempty"`
	Surname2        string `json:"apellido2,omitempty"`
	Phone           string `json:"telefono"`
	Phone2          string `json:"telefono2,omitempty"`
	PhoneWork       string `json:"telefonoTrabajo,omitempty"`
	MobilePhone     string `json:"movil,omitempty"`
	MobilePhone2    string `json:"movil2,omitempty"`
	Address         string `json:"direccion,omitempty"`
	PostalCode      string `json:"codigoPostal,omitempty"`
	Town            string `json:"poblacion,omitempty"`
	State           string `json:"provincia,omitempty"`
	Country         string `json:"pais,omitempty"`
	Fax             string `json:"fax,omitempty"`
	Email           string `json:"email,omitempty"`
	Email2          string `json:"emaiL2,omitempty"`
	BirthDate       string `json:"fechA_NACIMIENTO,omitempty"`
	SignupDate      string `json:"fechA_ALTA,omitempty"`
	LanguageID      int64  `json:"iD_IDIOMA,omitempty"`
	Observations    string `json:"observaciones,omitempty"`
	LocatableSince  int64  `json:"localizablE_DESDE,omitempty"`
	LocatableFrom   int64  `json:"localizablE_HASTA,omitempty"`
	DNI             string `json:"sDNI,omitempty"`
	FullName        string `json:"sNombre_Completo,omitempty"`
	Company         string `json:"sEmpresa,omitempty"`
	Sex             string `json:"cSexo,omitempty"`
	Text1           string `json:"textO1,omitempty"`
	Text2           string `json:"textO2,omitempty"`
	Text3           string `json:"textO3,omitempty"`
	FavSource       int64  `json:"nCanalPreferencial,omitempty"`
	Num1            int64  `json:"nuM1,omitempty"`
	Num2            int64  `json:"nuM2,omitempty"`
	Num3            int64  `json:"nuM3,omitempty"`
	SegmentAttibute string `json:"atributo_Segmento,omitempty"`
	Priority        int64  `json:"prioridad,omitempty"`
	NextContact     string `json:"tProximo_Contacto,omitempty"`
	NState          string `json:"nEstado,omitempty"`
	// Skill           string `json:"atributo_Skill,omitempty"`
	// NList           string `json:"nLista,omitempty"`
}

// AdditionalData bla
type AdditionalData struct {
	AddProp1 string `json:"additionalProp1,omitempty"`
	AddProp2 string `json:"additionalProp2,omitempty"`
	AddProp3 string `json:"additionalProp3,omitempty"`
}

// Localizators bla
type Localizators struct {
	AddProp1 int64 `json:"additionalProp1,omitempty"`
	AddProp2 int64 `json:"additionalProp2,omitempty"`
	AddProp3 int64 `json:"additionalProp3,omitempty"`
}

// EvolutionResp represents response from Evolution API
type EvolutionResp struct {
	Success bool
}

// LeadToEvolution maps leads values to Evolution struct
func (lead Lead) LeadToEvolution() Evolution {
	evolution := Evolution{
		Properties: Properties{
			OriginalID: *lead.LeaPhone,
			CampaingID: lead.SouIDEvolution,
			Phone:      *lead.LeaPhone,
			Name:       *lead.LeaName,
			Email:      *lead.LeaMail,
			DNI:        *lead.LeaDNI,
		},
		AdditionalData: AdditionalData{
			AddProp1: *lead.LeaURL,
			AddProp2: *lead.LeaIP,
			AddProp3: *lead.GaClientID,
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

// GetEvolutionValues queries for the Evolution sou_id equivalences
func (lead *Lead) GetEvolutionValues(db *gorm.DB) error {
	// source := Source{}
	// leatype := Leatype{}

	// if result := db.Where("sou_id = ?", lead.SouID).First(&source); result.Error != nil {
	// 	return fmt.Errorf("Error retrieving SouIDLeontel value: %#v", result.Error)
	// }
	// if result := db.Where("leatype_id = ?", lead.LeatypeID).First(&leatype); result.Error != nil {
	// 	return fmt.Errorf("error retrieving LeatypeIDLeontel value: %#v", result.Error)
	// }
	// lead.SouIDLeontel = source.SouIdcrm
	// lead.SouDescLeontel = source.SouDescription
	// lead.LeatypeIDLeontel = leatype.LeatypeIdcrm
	// lead.LeatypeDescLeontel = leatype.LeatypeDescription
	return nil
}

// Active is an implementation of Active method from Scable interface
func (e Evolution) Active(lead Lead) bool {
	return true
}

// Send is an implementation of Send method from Scable interface
func (e Evolution) Send(lead Lead) ScResponse {
	leadevolution := lead.LeadToEvolution()
	bytevalues, err := json.Marshal(leadevolution)
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

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
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
	evolutionresp := EvolutionResp{}
	if err := json.Unmarshal(data, &evolutionresp); err != nil {
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}

	if !evolutionresp.Success {
		return ScResponse{
			Success:    evolutionresp.Success,
			StatusCode: http.StatusUnprocessableEntity,
			ID:         0,
			Error:      errors.New(fmt.Sprintf("Evolution response %t", evolutionresp.Success)),
		}
	}

	return ScResponse{
		Success:    evolutionresp.Success,
		StatusCode: http.StatusOK,
		// TODO should be SubjectID field, not implemented by provider
		ID:    0,
		Error: nil,
	}
}
