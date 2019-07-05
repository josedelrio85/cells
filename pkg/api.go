package leads

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	model "github.com/bysidecar/leads/pkg/model"
)

// Handler is a struct created to use its ch property as element that implements
// http.Handler.Neededed to call HandleFunction as param in router Handler function.
type Handler struct {
	ch     http.Handler
	Storer Storer
	Lead   model.Lead
}

// HandleFunction is a function used to manage all received requests.
func (ch *Handler) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch.Lead = model.Lead{}
		if err := ch.Lead.Decode(r.Body); err != nil {
			message := fmt.Sprintf("Error decoding lead, Err: %v", err)
			responseError(w, message, err)
			return
		}

		// todo set lea_destiny value

		if err := ch.Lead.GetLeontelValues(ch.Storer.Instance()); err != nil {
			message := fmt.Sprintf("Error retrieving Leontel values, Err: %v", err)
			responseError(w, message, err)
			return
		}

		// todo delete this if when you want to test leontel insert
		if false {
			leonresp, err := ch.Lead.SendLeadToLeontel()
			if err != nil {
				message := fmt.Sprintf("Error sending lead to Leontel, Err: %v", err)
				// todo| think about what to do if leontel insert fails. I think we should log
				// todo| the error but the flow should continue.
				responseError(w, message, err)
				return
			}
			leontelID := strconv.FormatInt(leonresp.ID, 10)
			ch.Lead.LeaCrmid = &leontelID
		}

		if err := ch.Storer.Insert(&ch.Lead); err != nil {
			message := fmt.Sprintf("Error inserting lead in BD, Err: %v", err)
			responseError(w, message, err)
			return
		}

		fmt.Println(&ch.Lead)
		json.NewEncoder(w).Encode(&ch.Lead)
	})
}

// TODO think about a way to handle test request in prod environment.
// TODO One task is the skill to "cancel" sending the lead to Leontel queue
// TODO other task is the skill to substitute requested sou_id for test queue sou_id (15)
// TODO another task is to set the property lea_destiny from 'LEONTEL' to 'TEST'
// TODO think about what kind of response want to offer. Must respect the dependant code distributed on prod environment
