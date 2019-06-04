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

	// ch.Storer.Instance().Order("lea_id desc").Limit(1).Find(&ch.Lead)
	// ch.Lead.UpdatePostLeontel(ch.Storer.Instance(), 99999)
	if err := ch.Lead.Decode(req.Body); err != nil {
		message := fmt.Sprintf("Error decoding lead, Err: %v", err)
		responseError(w, message, err)
		return
	}

	leontel := *ch.Lead.LeaAux6

	fmt.Println(leontel)

	if leontel == "true"{
		fmt.Println("it is true!")
	}else {
		fmt.Println("it is false!")
	}

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
	// check if the time of the request complains with the on time requirementes. Returns true if ok, false otherwise.
	// ch.Lead.LeatypeID = isInSchedule(ch.Lead.SouID) ? 1 : 20;
	ch.Lead.LeatypeID = 1

	ch.proccessLead(w, req)
}

// RcableExpHandler blablabla
func (ch *Handler) RcableExpHandler(w http.ResponseWriter, req *http.Request) {
	if err := ch.Lead.Decode(req.Body); err != nil {
		message := fmt.Sprintf("Error decoding lead, Err: %v", err)
		responseError(w, message, err)
		return
	}

	// Data received from LP
	// phone
	// lea_url
	// typeCob		=>  when it is not C2C can be 2(form) or 24(form coverture) 		=> can be setted in lea_type
	// values 		=> 	only when typeCob is null, setted in observations
	// leontel		=>  not setted, used as flag => must be mapped in lea_aux6 as text (temporaly? create new field bool type in webservice.leads?)
	// test				=> 	?
	// name				=> 	setted in observaciones but must be mapped in name
	// location		=> 	setted in observaciones but must be mapped in	lea_aux4
	// answer			=> 	setted in observaciones but must be mapped in lea_aux5

	// TODO changes will be made in JS code in r-euskaltel-web project
	// TODO keys of the object sended by POST must be renamed to
	// TODO  typeCob 		=> lea_type
	// TODO  values 		=> observations
	// TODO  name 			=> name
	// TODO  location 	=> lea_aux4
	// TODO  answer 		=> lea_aux5
	// TODO  leontel 		=> lea_aux6

	leatypecob := ch.Lead.LeatypeID

	leontel := "false"
	if ch.Lead.LeaAux6 != nil {
		leontel = *ch.Lead.LeaAux6
	}

	answer := ""
	if ch.Lead.LeaAux5 != nil {
		answer = *ch.Lead.LeaAux5
	}

	observations := ""
	if ch.Lead.Observations !=  nil {
		observations = *ch.Lead.Observations
	}

	if leatypecob == 2 {
		observations := fmt.Sprintf("%v -- %v -- %v -- %v", *ch.Lead.LeaName, *ch.Lead.LeaAux4, observations, answer)
		ch.Lead.Observations = &observations
	}
	// it is not neccessary set values in observations for other situations because it was mapped automatically
	
	if leatypecob == 24 {
		coverture := "COBERTURA KO"	
		if leontel == "true" {
			coverture = "COBERTURA OK"	
		}
		ch.Lead.LeaAux3 = &coverture
	}

	if leatypecob == 0 {
		// TODO implement isInSchedule function (not this project)
		// check if the time of the request complains with the on time requirementes. Returns true if ok, false otherwise.
		// ch.Lead.LeatypeID = isInSchedule(ch.Lead.SouID) ? 1 : 8;
		ch.Lead.LeatypeID = 1
	}

	if leontel != "true" {
		destiny := "---"
		if len(answer) > 0 {
			destiny = 	"CONCURSO"
		}
		ch.Lead.LeaDestiny = &destiny
	}

	// Data sended to Leontel 
	// observaciones
	// observaciones2 (leatype_description)


	// 54
	ch.Lead.SouID = 15

	ch.proccessLead(w, req)
}

// TODO think about a way to handle test request in prod environment. 
// TODO One task is the skill to "cancel" sending the lead to Leontel queue
// TODO other task is the skill to substitute requested sou_id for test queue sou_id (15)
// TODO another task is to set the property lea_destiny from 'LEONTEL' to 'TEST'

// TODO think about what kind of response want to offer. Must respect the dependant code distributed on prod environment 


// MicrosoftHandler blablabla
func (ch *Handler) MicrosoftHandler(w http.ResponseWriter, req *http.Request) {

	if err := ch.Lead.Decode(req.Body); err != nil {
		message := fmt.Sprintf("Error decoding lead, Err: %v", err)
		responseError(w, message, err)
		return
	}
	ch.Lead.SouID = 15	
	
	ch.proccessLead(w, req)
}

// proccessLead handles all the proccess to set unattended values, insert the lead 
// in webservice DB and send the lead to Leontel environment.
func (ch *Handler) proccessLead(w http.ResponseWriter, req *http.Request) {
	if ch.Lead.LeaDestiny == nil {
		destiny := "LEONTEL"
		ch.Lead.LeaDestiny = &destiny		
	}

	if err := ch.Lead.GetLeontelValues(ch.Storer.Instance()); err != nil {
		responseError(w, "Error retrieving Leontel values ", err)
		return
	}
	if err:= ch.Lead.GetParams(w, req); err != nil {
		message := fmt.Sprintf("Error handling URL and IP params, Err: %v", err)
		responseError(w, message, err)
		return	
	}

	newLead := ch.Lead
	newLead.LeaID = 0

	if err := ch.Storer.Insert(&newLead); err != nil {
		message := fmt.Sprintf("Error inserting lead in BD, Err: %v", err)
		responseError(w, message, err)
		return	
	}

	if false {
		leonresp, err := newLead.SendLeadToLeontel()
		if err != nil {
			message := fmt.Sprintf("Error sending lead to Leontel, Err: %v", err)
			responseError(w, message, err)
			return
		}
	
		if leonresp.Success {
			if err := ch.Lead.UpdatePostLeontel(ch.Storer.Instance(), leonresp.ID); err !=  nil {
				message := fmt.Sprintf("Error updating lead post insert in Leontel, Err: %v", err)
				responseError(w, message, err)
				return
			}
		}	
		responseLeontel(w, leonresp)
	}
	responseOk(w, "OK")	
}


// TODO maybe move Middleware function to its own file. 

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
 