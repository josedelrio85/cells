package leads

import (
	"fmt"
	"time"

	"math/rand"
	"net/http"
	"strconv"

	redis "github.com/bysidecar/leads/pkg/leads/redis"

	"github.com/tomasen/realip"
)

// Handler is a struct created to use its ch property as element that implements
// http.Handler.Neededed to call HandleFunction as param in router Handler function.
type Handler struct {
	ch          http.Handler
	Storer      Storer
	Reporter    Storer
	Lead        Lead
	ActiveHooks []Hookable
	Redis       redis.Redis
	Dev         bool
}

// HandleFunction is a function used to manage all received requests.
func (ch *Handler) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch.Lead = Lead{}

		if err := ch.Lead.Decode(r.Body); err != nil {
			message := fmt.Sprintf("Error decoding lead, Err: %v", err)
			responseError(w, message, err)
			return
		}

		if ch.Lead.SouID == 0 {
			err := fmt.Errorf("Request does not contain sou_id value")
			message := fmt.Sprintf("Error decoding lead => no sou_id!, Err: %v", err)
			responseError(w, message, err)
			return
		}

		if ch.Lead.LeaIP == nil {
			requestIP := realip.FromRequest(r)
			ch.Lead.LeaIP = &requestIP
		}

		if ch.Lead.LeatypeID == 0 {
			ch.Lead.LeatypeID = 1
		}

		if err := ch.Lead.GetLeontelValues(ch.Storer.Instance()); err != nil {
			message := fmt.Sprintf("Error retrieving Leontel values, Err: %v", err)
			responseError(w, message, err)
			return
		}

		for _, hook := range ch.ActiveHooks {
			if !hook.Active(ch.Lead) {
				continue
			}

			hookResponse := hook.Perform(ch)
			if hookResponse.StatusCode == http.StatusUnprocessableEntity {
				message := fmt.Sprintf("An Unprocessable Entity was detected, Err: %v", hookResponse.Err)
				responseUnprocessable(w, message, hookResponse.Err)
				return
			}

			if hookResponse.Err != nil {
				responseError(w, hookResponse.Err.Error(), hookResponse.Err)
				return
			}
		}

		if err := ch.Lead.GetPassport(); err != nil {
			message := fmt.Sprintf("Error retrieving passport, Err: %v", err)
			responseError(w, message, err)
			return
		}

		if err := ch.Storer.Insert(&ch.Lead); err != nil {
			message := fmt.Sprintf("Error inserting lead in BD, Err: %v", err)
			responseError(w, message, err)
			return
		}

		if !ch.Dev {
			if err := ch.Reporter.Insert(&ch.Lead); err != nil {
				message := fmt.Sprintf("Error inserting lead into lead-report, Err: %v", err)
				responseError(w, message, err)
				return
			}
		}

		if ch.Lead.IsSmartCenter {
			leonresp, err := ch.Lead.SendLeadToLeontel()
			if err != nil {
				message := fmt.Sprintf("Error sending lead to SmartCenter, Err: %v", err)
				// TODO should break the flow? Maybe pass some info to responseOK method and handle the response in client
				responseUnprocessable(w, message, err)
			}
			leontelID := strconv.FormatInt(leonresp.ID, 10)
			ch.Lead.LeaSmartcenterID = leontelID

			cond := fmt.Sprintf("ID=%d", ch.Lead.ID)
			fields := []string{"LeaSmartcenterID", leontelID}
			ch.Storer.Update(Lead{}, cond, fields)
		}

		id := fmt.Sprintf("%d", ch.Lead.ID)
		responseOk(w, id, ch.Lead.IsSmartCenter)
	})
}

//HelperRandstring lalalal
func HelperRandstring(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
