package apic2c

import (
	"fmt"
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

// TestHandler blablabla
func (ch *Handler) TestHandler(w http.ResponseWriter, req *http.Request) {

	ch.Storer.Instance().Order("lea_id desc").Limit(1).Find(&ch.Lead)
	ch.Lead.UpdatePostLeontel(ch.Storer.Instance(), 99999)

	responseOk(w, "TestHandler")
}

// RcableHandler blablabla
func (ch *Handler) RcableHandler(w http.ResponseWriter, req *http.Request) {

	if err := ch.Lead.Decode(req.Body); err != nil {
		message := fmt.Sprintf("Error decoding lead, Err: %v", err)
		responseError(w, message, err)
		return
	}

	ch.Lead.SouID = 15
	destiny := "LEONTEL"
	ch.Lead.LeaDestiny = &destiny
	// check if the time of the request complains with the on time requirementes. Returns true if ok, false otherwise.
	// ch.Lead.LeatypeID = isInSchedule(ch.Lead.SouID) ? 1 : 20;
	ch.Lead.LeatypeID = 1
	if err := ch.Lead.GetLeontelValues(ch.Storer.Instance()); err != nil {
		responseError(w, "Error retrieving Leontel values ", err)
		return
	}
	if err:= ch.Lead.GetParams(w, req); err != nil {
		responseError(w, "Error handling URL and IP params", err)
		return	
	}

	if err := ch.Storer.Insert(&ch.Lead); err != nil {
		responseError(w, "Error inserting lead in BD ", err)
		return	
	}

	leonresp, err := ch.Lead.SendLeadToLeontel()
	if err != nil {
		responseError(w, "Error sending lead to Leontel, Err: %v", err)
		return
	}

	if leonresp.Success {
		if err := ch.Lead.UpdatePostLeontel(ch.Storer.Instance(), leonresp.ID); err !=  nil {
			responseError(w, "Error sending lead to Leontel, Err: %v", err)
			return
		}
	}

	responseLeontel(w, leonresp)
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
 
// TODO think about a way to handle test request in prod environment. 
// TODO One task is the skill to "cancel" sending the lead to Leontel queue
// TODO other task is the skill to substitute requested sou_id for test queue sou_id (15)
// TODO another task is to set the property lea_destiny from 'LEONTEL' to 'TEST'

// TODO think about what kind of response want to offer. Must respect the dependant code distributed on prod environment 
