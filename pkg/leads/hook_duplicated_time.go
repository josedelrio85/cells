package leads

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// DuplicatedTime is a is a struct that represents a DuplicatedTime entity
type DuplicatedTime struct {
	ExpirationTime int
}

// Active implents the Hookable interface, so when checking
// for active hooks will trigger the hook
// when the SouID matches a closed list.
//
// lead: The lead to check on.
//
// Returns true if the hook gets activated.
func (d DuplicatedTime) Active(lead Lead) bool {
	// TODO should we set this condition to all campaigns?
	switch lead.SouID {
	case 15:
		return true
	case 64, 65, 66:
		return true
	case 63:
		return true
	case 70:
		return true
	default:
		return false
	}
}

// Perform returns the result of duplicated validation
// cont: pointer to Handler struct
// Returns a HookReponse with the duplicated time check result.
func (d DuplicatedTime) Perform(cont *Handler) HookResponse {
	lead := &cont.Lead
	phone := *lead.LeaPhone
	key := fmt.Sprintf("%s-%d-%d", phone, lead.SouID, lead.LeatypeID)

	redisvalue, err := cont.Redis.Get(key)
	if err != nil {
		return HookResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
		}
	}

	// if a value is returned from Get than its length > 0
	// then there is a match in redis environment and we must reject the lead
	if len(*redisvalue) > 0 {
		message := fmt.Sprintf("Max attempts limit reached %s", key)
		return HookResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        errors.New(message),
		}
	}

	expirationtime := d.getExpirationTime(lead.SouID)
	// if ther isn't a value we set it
	cont.Redis.Set(key, phone, expirationtime)

	return HookResponse{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}

// getExpirationTime retrieves the ammount of time in which a key will expire
// souid the value of campaign
// Returns an integer
func (d DuplicatedTime) getExpirationTime(souid int64) int {
	switch souid {
	case 64, 65, 66:
		return 180
	case 63:
		return 180
	case 15:
		return 3
	default:
		return 60
	}
}
