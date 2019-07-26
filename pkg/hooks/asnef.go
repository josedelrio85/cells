package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	model "github.com/bysidecar/leads/pkg/model"
	"github.com/jinzhu/gorm"
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
	Db     *gorm.DB
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
func (a Asnef) Perform(db *gorm.DB, lead *model.Lead) HookResponse {

	if lead.Creditea.ASNEF || lead.Creditea.AlreadyClient {
		lead.IsSmartCenter = false
		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	}

	preresult, err := helper(db, lead)
	if err != nil {
		lead.IsSmartCenter = false
		return HookResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if preresult {
		lead.IsSmartCenter = false
		lead.Creditea.ASNEF = true
		lead.Creditea.AlreadyClient = true
		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	}

	asnefdata := InputData{
		SouID: lead.SouID,
		DNI:   *lead.LeaDNI,
		Phone: *lead.LeaPhone,
	}

	url := "https://ws.bysidecar.es/lead/asnef/check"
	var statuscode int

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
		// TODO how to difference if yacliente OR asnef ?
		lead.Creditea.ASNEF = true
		lead.Creditea.AlreadyClient = true
	}

	return HookResponse{
		StatusCode: statuscode,
		Err:        err,
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

// helper makes a prevalidation in leads BD to check for any match in the las month.
// If the conditions are matched, returns true.
func helper(db *gorm.DB, lead *model.Lead) (bool, error) {
	// si hay resultados => asnef positivo  || si no => sigue comprobando otra validación
	leadalt := model.Lead{}

	source := model.Source{}
	if result := db.Where("sou_id = ?", lead.SouID).First(&source); result.Error != nil {
		log.Fatalf("Error retrieving SouIDLeontel value: %v", result.Error)
	}
	soudesc := fmt.Sprintf("%s%s%s", "%", source.SouDescription[:5], "%")

	sources := []model.Source{}
	db.Where("sou_description like ?", soudesc).Find(&sources)

	stringsources := make([]string, 0)
	for _, s := range sources {
		stringsources = append(stringsources, fmt.Sprintf("%d", s.SouID))
	}

	dni := fmt.Sprintf("%s%s%s", "%", *lead.LeaDNI, "%")
	oneMonthLess := time.Now().AddDate(0, -1, 0)
	datecontrol := oneMonthLess.Format("2006-01-02")

	query := db.Debug().Table("leadnew").Select("leadnew.ID")
	query = query.Joins("JOIN creditea on leadnew.id = creditea.lea_id")
	query = query.Where("leadnew.lea_ts > ?", datecontrol)
	query = query.Where("leadnew.sou_id IN (?)", stringsources)
	query = query.Where("leadnew.is_smart_center = ?", 0)
	query = query.Where("leadnew.lea_dni like ? or leadnew.lea_phone = ?", dni, lead.LeaPhone)
	query = query.Where("creditea.asnef = ? or creditea.already_client = ?", 1, 1)
	err := query.First(&leadalt).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return false, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	return true, nil
}

// GetCandidatesPreasnef retrieves a list of candidates to match a positive asnef prevalidation.
// Used in test method only.
func GetCandidatesPreasnef(db *gorm.DB) []model.Lead {
	candidates := []model.Lead{}

	sources := []model.Source{}
	db.Where("sou_description like ?", "%CREDI%").Find(&sources)

	stringsources := make([]string, 0)
	for _, s := range sources {
		stringsources = append(stringsources, fmt.Sprintf("%d", s.SouID))
	}

	oneMonthLess := time.Now().AddDate(0, -1, 0)
	datecontrol := oneMonthLess.Format("2006-01-02")

	query := db.Table("leadnew").Select("leadnew.lea_phone, leadnew.lea_dni")
	query = query.Joins("JOIN creditea on leadnew.id = creditea.lea_id")
	query = query.Where("leadnew.lea_ts > ?", datecontrol)
	query = query.Where("leadnew.sou_id IN (?)", stringsources)
	query = query.Where("leadnew.is_smart_center = ?", 0)
	query = query.Where("creditea.asnef = ? or creditea.already_client = ?", 1, 1)
	query = query.Find(&candidates).Group("lea_dni")

	return candidates
}
