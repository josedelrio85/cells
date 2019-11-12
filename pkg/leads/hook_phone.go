package leads

import (
	"fmt"
	"net/http"

	"github.com/nyaruka/phonenumbers"
	"github.com/pkg/errors"
)

// Phone is a struct that represents a Phone entity
type Phone struct{}

// Active implents the Hookable interface, so when checking
// for active hooks will trigger the hook
// when the SouID matches a closed list.
// This hook checks for the presence of phone parameter in request
//
// lead: The lead to check on.
//
// Returns true if the hook gets activated.
func (p Phone) Active(lead Lead) bool {
	// TODO liable to add campaing discriminator
	return true
}

// Perform returns the result of phone validation
// cont: pointer to Handler struct
// Returns a HookReponse
func (p Phone) Perform(cont *Handler) HookResponse {
	if cont.Lead.LeaPhone == nil {
		return HookResponse{
			Err:        errors.New(fmt.Sprintf("Not allowed, phone param needed")),
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	num, err := phonenumbers.Parse(*cont.Lead.LeaPhone, "ES")
	if err != nil {
		return HookResponse{
			Err:        errors.New(fmt.Sprintf("Not allowed => %v", err)),
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	if !phonenumbers.IsValidNumber(num) {
		return HookResponse{
			Err:        errors.New(fmt.Sprintf("Not allowed, phone number not valid")),
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	list := map[string]bool{
		"666666666": true,
		"999999999": true,
	}
	if list[*cont.Lead.LeaPhone] {
		return HookResponse{
			Err:        errors.New(fmt.Sprintf("Not allowed, quarantine phone number")),
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	return HookResponse{
		Err:        nil,
		StatusCode: http.StatusOK,
	}
}
