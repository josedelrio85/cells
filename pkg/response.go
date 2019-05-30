package apic2c

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

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
	fancyHandleError(err)

	alarm := voalarm.NewClient("")
	_, err = alarm.SendAlarm(voalarm.Acknowledgement, err)
	if err != nil {
		fancyHandleError(err)
	}

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
