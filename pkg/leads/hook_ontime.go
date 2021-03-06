package leads

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// Ontime represents the struct in which store the result of the validations
type Ontime struct {
	Result        bool `json:"result"`
	ResultOntime  bool
	ResultHoliday bool
	Db            *gorm.DB
}

// InputDataOntime is the data to check asnef/already client validation
type InputDataOntime struct {
	SouID int64 `json:"sou_id"`
}

// Active implents the Hookable interface, so when checking for active hooks will trigger the Ontime validation hook when the LeatypeID matches a closed list.
// lead: The lead to check Ontime validation on.
// Returns true if the Ontime validation Hook gets activated.
func (a Ontime) Active(lead Lead) bool {
	switch lead.LeatypeID {
	case 1, 3, 4, 9:
		return true
	default:
		return false
	}
}

// Perform returns the result of ontime validation.
// If not on time => LeatypeID setted to 8 (FDH)
// lead: The lead to check ontime validation on.
// Returns a HookReponse with the ontime check result.
// True => ontime | false => holiday || out of time
func (a Ontime) Perform(cont *Handler) HookResponse {
	lead := &cont.Lead
	statuscode := http.StatusOK

	var err error
	inputholiday := InputDataOntime{
		SouID: lead.SouIDLeontel,
	}

	if err = a.checkHoliday(inputholiday); err != nil {
		log.Fatalf("Error checking holiday.Err: %s", err)
		statuscode = http.StatusInternalServerError
	}

	inputontime := InputDataOntime{
		SouID: lead.SouIDLeontel,
	}

	if err = a.checkOntime(inputontime); err != nil {
		log.Fatalf("Error checking ontime.Err: %s", err)
		statuscode = http.StatusInternalServerError
	}

	if a.ResultHoliday || !a.ResultOntime {
		lead.LeatypeID = 8
		lead.LeatypeIDLeontel = 12
	}

	return HookResponse{
		StatusCode: statuscode,
		Err:        err,
	}
}

// checkHoliday gets the result of holiday validation
func (a *Ontime) checkHoliday(input InputDataOntime) error {
	url := "https://ws.josedelrio85.es/smartcenter/timetable/isHoliday"

	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(input); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return err
	}
	a.ResultHoliday = a.Result
	return nil
}

// checkOntime gets the result of ontime validation
func (a *Ontime) checkOntime(input InputDataOntime) error {
	url := "https://ws.josedelrio85.es/smartcenter/timetable/isCampaignOnTime"

	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(input); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return err
	}
	a.ResultOntime = a.Result
	return nil
}
