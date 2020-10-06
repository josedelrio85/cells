package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// Lead struct represents the fields of the leads table
type Lead struct {
	gorm.Model
	LegacyID           int64      `json:"-"`
	LeaTs              *time.Time `sql:"DEFAULT:current_timestamp" json:"-" `
	LeaSmartcenterID   string     `json:"-"`
	PassportID         string     `json:"passport_id"`
	PassportIDGrp      string     `json:"passport_id_group"`
	SouID              int64      `json:"sou_id"`
	LeatypeID          int64      `json:"lea_type"`
	UtmSource          *string    `json:"utm_source,omitempty"`
	SubSource          *string    `json:"sub_source,omitempty"`
	LeaPhone           *string    `json:"phone,omitempty"`
	LeaMail            *string    `json:"mail,omitempty"`
	LeaName            *string    `json:"name,omitempty"`
	LeaDNI             *string    `json:"dni,omitempty"`
	LeaURL             *string    `sql:"type:text" json:"url,omitempty"`
	LeaIP              *string    `json:"ip,omitempty"`
	GaClientID         *string    `json:"ga_client_id,omitempty"`
	IsSmartCenter      bool       `json:"smartcenter,omitempty"`
	SouIDLeontel       int64      `sql:"-" json:"sou_id_leontel"`
	SouDescLeontel     string     `sql:"-" json:"sou_desc_leontel"`
	LeatypeIDLeontel   int64      `sql:"-" json:"lea_type_leontel"`
	LeatypeDescLeontel string     `sql:"-" json:"lea_type_desc_leontel"`
	SouIDEvolution     int64      `sql:"-" json:"sou_id_evolution"`
	Gclid              *string    `json:"gclid,omitempty"`
	Domain             *string    `json:"domain,omitempty"`
	Observations       *string    `sql:"type:text" json:"observations,omitempty"`
	RequestID          string     `json:"-"`
	RcableExp          *RcableExp `json:"rcableexp"`
	Microsoft          *Microsoft `json:"microsoft"`
	Creditea           *Creditea  `json:"creditea"`
	Kinkon             *Kinkon    `json:"kinkon"`
	Alterna            *Alterna   `json:"alterna"`
	Adeslas            *Adeslas   `json:"adeslas"`
	Endesa             *Endesa    `json:"endesa"`
	Virgin             *Virgin    `json:"virgin"`
}

// TableName sets the default table name
func (Lead) TableName() string {
	return "leads"
}

// Decode reques's body into a Lead struct
func (lead *Lead) Decode(body io.ReadCloser) error {
	if err := json.NewDecoder(body).Decode(lead); err != nil {
		return err
	}
	return nil
}

// GetParams retrieves values for ip, port and url properties
func (lead *Lead) GetParams(w http.ResponseWriter, req *http.Request) error {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return fmt.Errorf("Error parsing IP value: %#v", err)
	}
	// forward := req.Header.Get("X-Forwarded-For")

	lead.LeaIP = &ip
	url := fmt.Sprintf("%s%s", req.Host, req.URL.Path)
	lead.LeaURL = &url

	return nil
}

// concatPointerStrs concats an undefined number
// of *string params into string separated by --
func concatPointerStrs(args ...*string) string {
	var buffer bytes.Buffer
	tam := len(args) - 1
	for i, arg := range args {
		if arg != nil {
			buffer.WriteString(*arg)
			if i < tam {
				buffer.WriteString(" -- ")
			}
		}
	}
	return buffer.String()
}

// GetPassport gets a passport and sets it into lead properties
func (lead *Lead) GetPassport() error {
	passport := Passport{}

	// TODO maybe use desc from leads instead Leontel
	interaction := Interaction{
		Provider:    lead.SouDescLeontel,
		Application: lead.LeatypeDescLeontel,
		IP:          *lead.LeaIP,
	}

	if err := passport.Get(interaction); err != nil {
		return err
	}

	lead.PassportID = passport.PassportID
	lead.PassportIDGrp = passport.PassportIDGrp

	return nil
}

// GetSourceValues queries for the SmartCenter equivalences
// of sou_id and lea_type values
func (lead *Lead) GetSourceValues(db *gorm.DB) error {
	source := Source{}
	leatype := Leatype{}

	if result := db.Where("sou_id = ?", lead.SouID).First(&source); result.Error != nil {
		return fmt.Errorf("Error retrieving SouIDLeontel value: %#v", result.Error)
	}
	if result := db.Where("leatype_id = ?", lead.LeatypeID).First(&leatype); result.Error != nil {
		return fmt.Errorf("error retrieving LeatypeIDLeontel value: %#v", result.Error)
	}
	lead.SouIDLeontel = source.SouIdcrm
	lead.SouDescLeontel = source.SouDescription
	lead.SouIDEvolution = source.SouIDEvolution
	lead.LeatypeIDLeontel = leatype.LeatypeIdcrm
	lead.LeatypeDescLeontel = leatype.LeatypeDescription
	return nil
}
