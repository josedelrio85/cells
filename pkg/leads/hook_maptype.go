package leads

import (
	"log"
	"net/http"
)

// MapType is a is a struct that represents a maptype entity
type MapType struct{}

// Active implents the Hookable interface, so when checking
// for active hooks will trigger the hook
// when the TypeId matches a closed list.
//
// lead: The lead to check on.
//
// Returns true if the hook gets activated.
func (a MapType) Active(lead Lead) bool {
	switch lead.SouID {
	case 74, 75, 76:
		switch lead.LeatypeID {
		case 1, 2, 8, 24:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// Perform returns the result of MapType validation
// cont: pointer to Handler struct
// Returns a HookReponse with 200 Status and updates lea_type value
func (a MapType) Perform(cont *Handler) HookResponse {
	log.Println("perform maptype hook")
	listSource := map[int64]bool{
		74: true,
		75: true,
		76: true,
	}

	listType := map[int64]bool{
		1:  true,
		2:  true,
		8:  true,
		24: true,
	}

	if listSource[cont.Lead.SouID] && listType[cont.Lead.LeatypeID] {
		cont.Lead.LeatypeID = 9
		if err := cont.Lead.GetLeontelValues(cont.Storer.Instance()); err != nil {
			return HookResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	log.Printf("leatype id %d", cont.Lead.LeatypeID)
	log.Printf("leatypeLeontel id %d", cont.Lead.LeatypeIDLeontel)
	return HookResponse{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}
