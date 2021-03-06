package leads

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/josedelrio85/voalarm"
)

// ResponseAPI represents the data structure needed to create a response
type ResponseAPI struct {
	Code        int
	Message     string `json:"message"`
	Success     bool   `json:"success"`
	SmartCenter bool   `json:"smartcenter"`
}

// response sets the params to generate a JSON response
func response(w http.ResponseWriter, ra ResponseAPI) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ra.Code)

	result := struct {
		Success     bool   `json:"success"`
		Message     string `json:"message"`
		SmartCenter bool   `json:"smartcenter"`
	}{
		Success:     ra.Success,
		Message:     ra.Message,
		SmartCenter: ra.SmartCenter,
	}

	json.NewEncoder(w).Encode(result)
}

// responseError generates log, alarm and response when an error occurs
func responseError(w http.ResponseWriter, message string, err error) {
	log.Println(message)

	sendAlarm(message, http.StatusInternalServerError, err)

	ra := ResponseAPI{
		Code:        http.StatusInternalServerError,
		Message:     message,
		Success:     false,
		SmartCenter: false,
	}
	response(w, ra)
}

// responseOk calls response function with proper data to generate an OK response
func responseOk(w http.ResponseWriter, message string, smartcenter bool) {
	ra := ResponseAPI{
		Code:        http.StatusOK,
		Message:     message,
		Success:     true,
		SmartCenter: smartcenter,
	}
	response(w, ra)
}

// responseUnprocessable calls response function to inform user of something does not work 100% OK
func responseUnprocessable(w http.ResponseWriter, message string, err error) {
	sendAlarm(message, http.StatusUnprocessableEntity, err)

	ra := ResponseAPI{
		Code:        http.StatusUnprocessableEntity,
		Message:     message,
		Success:     true,
		SmartCenter: false,
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
func sendAlarm(message string, status int, err error) {
	fancyHandleError(err)

	mstype := voalarm.Acknowledgement
	switch status {
	case http.StatusInternalServerError:
		mstype = voalarm.Warning
	case http.StatusUnprocessableEntity:
		mstype = voalarm.Info
	}

	alarm := voalarm.NewClient("")
	_, err = alarm.SendAlarm(message, mstype, err)
	if err != nil {
		fancyHandleError(err)
	}
}
