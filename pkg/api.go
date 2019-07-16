package leads

import (
	"fmt"
	"net/http"
	"strconv"

	hooks "github.com/bysidecar/leads/pkg/hooks"
	model "github.com/bysidecar/leads/pkg/model"
)

// Handler is a struct created to use its ch property as element that implements
// http.Handler.Neededed to call HandleFunction as param in router Handler function.
type Handler struct {
	ch          http.Handler
	Storer      Storer
	Lead        model.Lead
	ActiveHooks []hooks.Hookable
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

		if ch.Lead.SouID == 0 {
			err := fmt.Errorf("Request does not contain sou_id value")
			message := fmt.Sprintf("Error decoding lead => no sou_id!, Err: %v", err)
			responseError(w, message, err)
			return
		}

		if ch.Lead.LeatypeID == 0 {
			ch.Lead.LeatypeID = 1
		}

		for _, hook := range ch.ActiveHooks {

			if !hook.Active(ch.Lead) {
				continue
			}

			hookResponse := hook.Perform(&ch.Lead)

			if hookResponse.Err != nil {
				responseError(w, hookResponse.Err.Error(), hookResponse.Err)
			}
		}

		// TODO think about hibernated campaings, should reject them?

		// TODO set lea_destiny value
		if err := ch.Lead.GetLeontelValues(ch.Storer.Instance()); err != nil {
			message := fmt.Sprintf("Error retrieving Leontel values, Err: %v", err)
			responseError(w, message, err)
			return
		}

		// TODO think about passport relation with lead. New properties on Lead entity?
		passport := Passport{}
		passport.Get(ch.Lead)

		// todo delete this, for dev purposes only
		ch.Lead.LeadToLeontel()

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

		// json.NewEncoder(w).Encode(&ch.Lead)
		id := fmt.Sprintf("%d", ch.Lead.ID)
		responseOk(w, id)
	})
}

// TODO think about a way to handle test request in prod environment.
// TODO One task is the skill to "cancel" sending the lead to Leontel queue
// TODO other task is the skill to substitute requested sou_id for test queue sou_id (15)
// TODO another task is to set the property lea_destiny from 'LEONTEL' to 'TEST'
// TODO think about what kind of response want to offer. Must respect the dependant code distributed on prod environment
