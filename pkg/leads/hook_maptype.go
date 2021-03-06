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
	// kinkon R only
	case 64:
		switch lead.LeatypeID {
		case 3, 8, 24, 27, 30:
			return true
		default:
			return false
		}
	// kinkon-empresas
	case 74, 75, 76:
		switch lead.LeatypeID {
		case 2, 8, 24:
			return true
		default:
			return false
		}
	// adeslas + endesa + virgin
	case 77, 78, 79:
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
		64: true,
		74: true,
		75: true,
		76: true,
		77: true,
		78: true,
		79: true,
	}

	listType := map[int64]bool{
		2:  true,
		3:  true,
		8:  true,
		24: true,
		27: true,
		30: true,
	}

	if listSource[cont.Lead.SouID] && listType[cont.Lead.LeatypeID] {
		getNewType(&cont.Lead)
		if err := cont.Lead.GetSourceValues(cont.Storer.Instance()); err != nil {
			return HookResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	// log.Printf("leatypeid %d", cont.Lead.LeatypeID)
	// log.Printf("LeatypeIDLeontel %d", cont.Lead.LeatypeIDLeontel)

	return HookResponse{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}

func getNewType(lead *Lead) {
	switch lead.SouID {
	case 64:
		switch lead.LeatypeID {
		case 3, 8, 24, 27, 30:
			lead.LeatypeID = 9
		}
	case 74, 75, 76:
		switch lead.LeatypeID {
		case 2, 8, 24:
			lead.LeatypeID = 9
		}
	case 77:
		switch lead.LeatypeID {
		case 8:
			lead.LeatypeID = 20
		}
	case 78, 79:
		switch lead.LeatypeID {
		case 8:
			lead.LeatypeID = 9
		}
	}
}
