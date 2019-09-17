package leads

import (
	"fmt"
	"net/http"

	container "github.com/bysidecar/leads/pkg/container"
	model "github.com/bysidecar/leads/pkg/model"

	"github.com/pkg/errors"
)

// Allowed is a is a struct that represents a Redis entity
type Allowed struct {
	Redis model.Redis
}

// Active implents the Hooable interface, so when checking
// for active hooks will trigger the hook
// when the SouID matches a closed list.
//
// lead: The lead to check on.
//
// Returns true if the hook gets activated.
func (a Allowed) Active(lead model.Lead) bool {
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

// Perform returns the result of allowed validation
// lead: The lead to check on.
// db: not used in this implementation
// Returns a HookReponse with the allowed check result.
func (a Allowed) Perform(cont container.Container) HookResponse {
	lead := cont.Lead
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
