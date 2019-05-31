package apic2c

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Handler is a struct created to use its ch property as element that implements
// http.Handler.Neededed to call HandleFunction as param in router Handler function.
type Handler struct {
	ch       http.Handler
	Storer   Storer
	Lead     Lead
	LeadTest LeadTest
}

// HandleFunction is a function used to manage all received requests.
// Only POST method accepted.
// Decode the identity json request ____
// Returns an StatusMethodNotAllowed state if other kind of request is received.
// Returns StatusInternalServerError when decoding the body content fails.
func (ch *Handler) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := ch.Lead.decode(r.Body); err != nil {
			message := fmt.Sprintf("Error decoding lead, Err: %v", err)
			responseError(w, message, err)
			return
		}
		responseOk(w, "OK")
	})
}

// TestHandler blablabla
func (ch *Handler) TestHandler(w http.ResponseWriter, req *http.Request) {
	responseOk(w, "TestHandler")
}

// RcableHandler blablabla
func (ch *Handler) RcableHandler(w http.ResponseWriter, req *http.Request) {

	if err := ch.Lead.decode(req.Body); err != nil {
		message := fmt.Sprintf("Error decoding lead, Err: %v", err)
		responseError(w, message, err)
		return
	}

	ch.Storer.Instance()

	ch.Storer.Insert(&ch.Lead)

	// leadLeontel := ch.Lead.LeadToLeontel()
	// log.Println(leadLeontel)
	responseOk(w, "OK")
}

// Middleware to preproccess requests. Checks for the correct content-type
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Method not allowed ", http.StatusMethodNotAllowed)
			http.Error(w, "Method not allowed ", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (lead *Lead) decode(body io.ReadCloser) error {
	if err := json.NewDecoder(body).Decode(lead); err != nil {
		return err
	}
	return nil
}
