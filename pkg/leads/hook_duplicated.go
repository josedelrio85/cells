package leads

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Duplicated is a is a struct that represents a Redis entity
type Duplicated struct{}

// Active implents the Hooable interface, so when checking
// for active hooks will trigger the hook
// when the SouID matches a closed list.
//
// lead: The lead to check on.
//
// Returns true if the hook gets activated.
func (a Duplicated) Active(lead Lead) bool {
	switch lead.SouID {
	case 64:
		return true
	case 65:
		return true
	case 66:
		return true
	default:
		return false
	}
}

// Perform returns the result of duplicated validation
// lead: The lead to check on.
// db: not used in this implementation
// Returns a HookReponse with the duplicated check result.
func (a Duplicated) Perform(cont *Handler) HookResponse {
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

	// if ther isn't a value we set it
	cont.Redis.Set(key, phone)

	return HookResponse{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}
