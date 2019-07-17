package leads

import (
	"net/http"

	model "github.com/bysidecar/leads/pkg/model"
)

type Asnef struct{}

// Active implents the Hooable interface, so when checking for active hooks will trigger the Asnef hook when the SouID matches a closed list.
//
// lead: The lead to check Asneff on.
//
// Returns true if the Asnef Hook gets activated.
func (a Asnef) Active(lead model.Lead) bool {
	if lead.SouID == 1 {
		return true
	}

	return false
}

// Perform is currently a fake function that has been created for hook explanaton purposes that will return always an OK HookResponse.
//
// TODO: implement the real Asnef check.
//
// lead: The lead to check Asnef on.
//
// Returns a HookReponse with the asnef check result.
func (a Asnef) Perform(lead model.Lead) HookResponse {
	return HookResponse{
		StatusCode: http.StatusOK,
	}
}
