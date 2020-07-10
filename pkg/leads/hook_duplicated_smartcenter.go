package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// DuplicatedSmartCenter is a struct that represents a DuplicatedSmartCenter entity
type DuplicatedSmartCenter struct{}

// DataLeontelResp represents the data content in the response of Leontel environment
type DataLeontelResp struct {
	Closed      string `json:"lea_closed"`
	LeaID       string `json:"lea_id"`
	LeaTs       string `json:"lea_ts"`
	Description string `json:"sub_description"`
	SubID       string `json:"sub_id"`
}

// RespSC is a struct that representes Leontel environment response
type RespSC struct {
	Success bool              `json:"success"`
	Data    []DataLeontelResp `json:"data"`
	Error   interface{}       `json:"error"`
}

// Active implents the Hookable interface, so when checking
// for active hooks will trigger the hook
// when the SouID matches a closed list.
//
// lead: The lead to check on.
//
// Returns true if the hook gets activated.
func (t DuplicatedSmartCenter) Active(lead Lead) bool {
	switch lead.SouID {
	case 15:
		return true
	case 64, 65, 66:
		return true
	case 74, 75, 76:
		return true
	case 63:
		return true
	case 69:
		return true
	case 70:
		return true
	case 77:
		return true
	case 78:
		return true
	default:
		return false
	}
}

// Perform returns the result of duplicated smartcenter validation
// cont: pointer to Handler struct
// Returns a HookReponse with the duplicated smartcenter check result.
func (t DuplicatedSmartCenter) Perform(cont *Handler) HookResponse {
	lead := &cont.Lead

	endpoint, ok := os.LookupEnv("CHECK_LEAD_LEONTEL_ENDPOINT")
	if !ok {
		message := "unable to load Check Lead Leontel URL endpoint"
		sendAlarm(message, http.StatusInternalServerError, errors.New("unable to load Check Lead Leontel URL endpoint"))
		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	}

	// Pass the leontel data for souid and type
	data := struct {
		Phone   string `json:"TELEFONO"`
		SouID   int64  `json:"lea_source"`
		LeaType int64  `json:"lea_type"`
	}{
		Phone:   *lead.LeaPhone,
		SouID:   lead.SouIDLeontel,
		LeaType: lead.LeatypeIDLeontel,
	}

	bytevalues, err := json.Marshal(data)
	if err != nil {
		message := "unable to marshal data"
		sendAlarm(message, http.StatusInternalServerError, err)

		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
	if err != nil {
		message := "unable to make POST request"
		sendAlarm(message, http.StatusInternalServerError, err)
		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	}
	defer resp.Body.Close()

	rawdata, _ := ioutil.ReadAll(resp.Body)
	structdata := RespSC{}

	if err := json.Unmarshal(rawdata, &structdata); err != nil {
		message := "unable to unmarshal lead status response"
		sendAlarm(message, http.StatusInternalServerError, err)
		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	}

	// if success is false the lead is not open and it is a positive result
	if !structdata.Success {
		return HookResponse{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	}

	message := fmt.Sprintf("Not allowed, lead %s already open ", structdata.Data[0].LeaID)
	return HookResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        errors.New(message),
	}
}
