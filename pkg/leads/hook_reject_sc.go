package leads

import (
	"log"
	"net/http"
)

// RejectSC is a is a struct that represents a reject smartcenter entity
type RejectSC struct{}

// Active sets activation of this hook
func (r RejectSC) Active(lead Lead) bool {
	listSource := map[int64]bool{
		71: true, // R CABLE WEB CARTERA
	}
	if listSource[lead.SouID] {
		log.Println("Reject SC hook activated")
		return true
	}
	return false
}

// Perform returns the result of RejectSC validation
// cont: pointer to Handler struct
// Returns a HookReponse with 200 Status and sets smart center value to false
func (r RejectSC) Perform(cont *Handler) HookResponse {

	cont.Lead.IsSmartCenter = false

	return HookResponse{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}
