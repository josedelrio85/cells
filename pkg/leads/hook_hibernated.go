package leads

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
)

// Hibernated is a list of campaigns that should be rejected
type Hibernated struct {
	List map[int64]bool
	Db   *gorm.DB
}

// Active implents the Hooable interface, so when checking for active hooks will trigger the hibernated campaign hook when the SouID matches a closed list.
//
// lead: The lead to check Asneff on.
//
// Returns true if the hibernated campaign Hook gets activated.
func (h Hibernated) Active(lead Lead) bool {
	hibernated := Hibernated{
		List: map[int64]bool{
			2:  true,
			10: true,
			11: true,
			17: true,
			18: true,
			19: true,
			20: true,
			24: true,
			26: true,
			27: true,
			28: true,
			29: true,
			30: true,
			31: true,
			32: true,
			33: true,
			34: true,
			35: true,
			36: true,
			37: true,
			38: true,
			39: true,
			40: true,
			41: true,
			42: true,
			43: true,
			44: true,
			47: true,
		},
	}
	// 4, 12, 16,
	h = hibernated

	if h.List[lead.SouID] || lead.SouIDLeontel == 0 {
		return true
	}
	return false
}

// Perform returns the result of hibernated campaign validation
// lead: The lead to check the hibernated campaign.
// Returns a HookReponse with the hibernated campaign check result
func (h Hibernated) Perform(cont *Handler) HookResponse {
	err := fmt.Errorf("An hibernated campaign was detected! => %d", cont.Lead.SouID)

	return HookResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        err,
	}
}
