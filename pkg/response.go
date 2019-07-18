package leads

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"

	model "github.com/bysidecar/leads/pkg/model"
	"github.com/bysidecar/voalarm"
)

// ResponseAPI represents the data structure needed to create a response
type ResponseAPI struct {
	Code    int
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// response sets the params to generate a JSON response
func response(w http.ResponseWriter, ra ResponseAPI) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ra.Code)

	result := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: ra.Success,
		Message: ra.Message,
	}

	json.NewEncoder(w).Encode(result)
}

// responseError generates log, alarm and response when an error occurs
func responseError(w http.ResponseWriter, message string, err error) {
	sendAlarm(message, err)

	ra := ResponseAPI{
		Code:    http.StatusInternalServerError,
		Message: message,
		Success: false,
	}
	response(w, ra)
}

// responseOk calls response function with proper data to generate an OK response
func responseOk(w http.ResponseWriter, message string) {
	ra := ResponseAPI{
		Code:    http.StatusOK,
		Message: message,
		Success: true,
	}
	response(w, ra)
}

// responseLeontel calls response function passing the data obtained from Leontel proxy
func responseLeontel(w http.ResponseWriter, resp *model.LeontelResp) {
	ra := ResponseAPI{
		Code:    http.StatusOK,
		Message: strconv.FormatInt(resp.ID, 10),
		Success: resp.Success,
	}
	response(w, ra)
}

// responseUnprocessable calls response function to inform user of something does not work 100% OK
func responseUnprocessable(w http.ResponseWriter, message string, err error) {
	sendAlarm(message, err)

	ra := ResponseAPI{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Success: true,
	}
	response(w, ra)
}

// fancyHandleError logs the error and indicates the line and function
func fancyHandleError(err error) (b bool) {
	if err != nil {
		// using 1 => it will actually log where the error happened, 0 = this function.
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
		b = true
	}
	return
}

// sendAlarm to VictorOps plattform and format the error for more info
func sendAlarm(message string, err error) {
	fancyHandleError(err)

	// TODO need to add another param into SendAlarm method ("leads") but first a vendorize process is needed
	alarm := voalarm.NewClient("")
	_, err = alarm.SendAlarm(voalarm.Acknowledgement, err)
	if err != nil {
		fancyHandleError(err)
	}
}
