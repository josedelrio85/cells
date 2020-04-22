package leads

import (
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
	case 77:
		switch lead.LeatypeID {
		case 8:
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
	listSource := map[int64]bool{
		74: true,
		75: true,
		76: true,
		77: true,
	}

	listType := map[int64]bool{
		1:  true,
		2:  true,
		8:  true,
		24: true,
	}

	if listSource[cont.Lead.SouID] && listType[cont.Lead.LeatypeID] {
		getNewType(&cont.Lead)
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

func getNewType(lead *Lead) {
	switch lead.SouID {
	case 74, 75, 76:
		switch lead.LeatypeID {
		case 1, 2, 8, 24:
			lead.LeatypeID = 9
		}
	case 77:
		switch lead.LeatypeID {
		case 8:
			lead.LeatypeID = 20
		}
	}
}
