package leads

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	model "github.com/bysidecar/leads/pkg/model"
)

// Passport struct represents the values of the Passport returned from Passport Service
type Passport struct {
	Idgroup string `json:"idgroup"`
	ID      string `json:"idbysidecar"`
}

// Interaction represents the structure needed to obtain a passport
type Interaction struct {
	Provider    string `json:"provider"`
	Application string ` json:"application"`
	IP          string `json:"ip"`
}

// Get function retrieves a passport for the incoming lead
func (p *Passport) Get(lead model.Lead) error {
	url := "https://passport.bysidecar.me/id/settle"

	interaction := Interaction{
		Provider:    lead.SouDescLeontel,
		Application: lead.LeatypeDescLeontel,
		IP:          *lead.LeaIP,
	}

	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(interaction); err != nil {
		log.Fatalf("Error on encoding struct data.  %s, Err: %s", interaction, err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		log.Fatalf("Error on creating request object.  %s, Err: %s", url, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error on making request. Err: %s", err)
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(p); err != nil {
		log.Fatalf("Error on decoding response from Passport.  %s, Err: %s", res.Body, err)
		return err
	}

	return nil
}
