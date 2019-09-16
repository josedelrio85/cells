package leads

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	model "github.com/bysidecar/leads/pkg/model"
	"github.com/jinzhu/gorm"
)

// Allowed is a struct
type Allowed struct {
	Result bool   `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
	Redis  model.Redis
}

// Active implents the Hooable interface, so when checking
// for active hooks will trigger the Asnef hook
// when the SouID matches a closed list.
//
// lead: The lead to check Asneff on.
//
// Returns true if the Asnef Hook gets activated.
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
// lead: The lead to check Asnef on.
// db: not used in this implementation
// Returns a HookReponse with the allowed check result.
func (a Allowed) Perform(db *gorm.DB, lead *model.Lead) HookResponse {

	phone := *lead.LeaPhone
	key := fmt.Sprintf("%s-%d-%d", phone, lead.SouID, lead.LeatypeID)
	_, err := a.Redis.Get(key)
	if err != nil {
		return HookResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
		}
	}

	a.Redis.Set(key, phone)

	return HookResponse{
		StatusCode: http.StatusOK,
		Err:        nil,
	}
}

func generateCode() (*string, error) {
	b := make([]byte, 8)
	n, err := io.ReadAtLeast(rand.Reader, b, 8)
	if n != 8 {
		return nil, err
	}
	for i := 0; i < len(b); i++ {
		b[i] = codeTable[int(b[i])%len(codeTable)]
	}

	code := string(b)
	return &code, nil
}

var codeTable = [...]byte{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
	'A', 'C', 'E', 'G', 'H', 'K', 'M', 'N', 'R', 'P', 'Q', 'R', 'S', 'X', 'Y', 'Z',
}
