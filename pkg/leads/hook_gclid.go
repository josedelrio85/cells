package leads

import (
	"net/http"
)

// Gclid is a is a struct that represents a Gclid entity
type Gclid struct{}

// Active implents the Hookable interface, so when checking
// for active hooks will trigger the hook
// when the SouID matches a closed list.
//
// lead: The lead to check on.
//
// Returns true if the hook gets activated.
func (a Gclid) Active(lead Lead) bool {
	switch lead.SouID {
	case 15:
		return true
	case 9:
		return true
	default:
		return false
	}
}

// Perform returns the result of gclid validation
// cont: pointer to Handler struct
// Returns a HookReponse with 200 Status and the updated sou_id value
func (a Gclid) Perform(cont *Handler) HookResponse {
	lead := &cont.Lead
	if lead.Gclid != nil {
		cont.Lead.SouID = a.getGclidSouID(cont.Lead.SouID)
		if err := cont.Lead.GetLeontelValues(cont.Storer.Instance()); err != nil {
			return HookResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}
	return HookResponse{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}

// getGclidSouID returns the equivalent Gclid source value for the incoming campaign
func (a Gclid) getGclidSouID(souid int64) int64 {
	switch souid {
	case 15:
		return 15
	case 9:
		// TODO change for the correct value
		return 9
	default:
		return souid
	}
}
