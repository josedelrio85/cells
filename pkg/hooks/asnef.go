package leads

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	model "github.com/bysidecar/leads/pkg/model"
)

// Candidate blablabla
type Candidate struct {
	LeaID     string `json:"lea_id,omitempty"`
	Telefono  string `json:"telefono,omitempty"`
	DNI       string `json:"dninie,omitempty"`
	LeaSource string `json:"lea_source,omitempty"`
}

// Asnef is a struct that represents the result provided by asnef/already client validation
type Asnef struct {
	Result bool   `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// InputData is the data to check asnef/already client validation
type InputData struct {
	SouID int64  `json:"sou_id"`
	DNI   string `json:"dni"`
	Phone string `json:"phone"`
}

// Active implents the Hooable interface, so when checking for active hooks will trigger the Asnef hook when the SouID matches a closed list.
//
// lead: The lead to check Asneff on.
//
// Returns true if the Asnef Hook gets activated.
func (a Asnef) Active(lead model.Lead) bool {
	if lead.IsSmartCenter {
		switch lead.SouID {
		case 9:
			return true
		case 58:
			return true
		default:
			return false
		}
	}
	return false
}

// Perform returns the result of asnef/already client validation
// lead: The lead to check Asnef on.
// Returns a HookReponse with the asnef check result.
func (a Asnef) Perform(lead *model.Lead) HookResponse {

	url := "https://ws.bysidecar.es/lead/asnef/check"
	var statuscode int

	if lead.Creditea.Motivo != nil {
		lead.IsSmartCenter = false
		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
			Result:     true,
		}
	}

	asnefdata := InputData{
		SouID: lead.SouID,
		DNI:   *lead.LeaDNI,
		Phone: *lead.LeaPhone,
	}

	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(asnefdata); err != nil {
		log.Fatalf("Error encoding asnef data. %v, Err: %s", asnefdata, err)
		statuscode = http.StatusInternalServerError
	}

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		log.Fatalf("Error on creating request object.  %s, Err: %s", url, err)
		statuscode = http.StatusInternalServerError
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request.  %s, Err: %s", url, err)
		statuscode = http.StatusInternalServerError
	}

	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		log.Fatalf("Error decoding asnef response. Err: %s", err)
		statuscode = http.StatusInternalServerError
	}

	statuscode = http.StatusOK

	if a.Result {
		lead.IsSmartCenter = false
		motivo := "Asnef/Ya cliente positivo"
		lead.Creditea.Motivo = &motivo
	}

	return HookResponse{
		StatusCode: statuscode,
		Err:        err,
		Result:     a.Result,
	}
}

// GetCandidates retrives a list of asnef candidates
func GetCandidates(lead model.Lead) []Candidate {

	url := "https://ws.bysidecar.es/lead/asnef/getcandidates"

	candidate := make([]Candidate, 0)
	data := new(bytes.Buffer)
	input := InputData{
		SouID: lead.SouID,
	}

	if err := json.NewEncoder(data).Encode(input); err != nil {
		log.Fatalf("Error encoding asnef data. Err: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		log.Fatalf("Error on creating request object.  %s, Err: %s", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request.  %s, Err: %s", url, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&candidate); err != nil {
		log.Fatalf("Error decoding asnef response. Err: %s", err)
	}

	return candidate
}
