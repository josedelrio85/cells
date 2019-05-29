package apic2c

import (
	"log"
	"net/http"
)

// Handler is a struct created to use its ch property as element that implements
// http.Handler.Neededed to call HandleFunction as param in router Handler function.
type Handler struct {
	ch http.Handler
}

// HandleFunction is a function used to manage all received requests.
// Only POST method accepted.
// Decode the identity json request ____
// Returns an StatusMethodNotAllowed state if other kind of request is received.
// Returns StatusInternalServerError when decoding the body content fails.
func (ch *Handler) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Method not allowed ", http.StatusMethodNotAllowed)
			http.Error(w, "Method not allowed ", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode()
	})
}
