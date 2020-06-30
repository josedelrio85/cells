package leads

import "github.com/jinzhu/gorm"

// Virgin struct represents the fields of the Virgin table
type Virgin struct {
	gorm.Model

	LeadID     uint    `sql:"column:lea_id"`
	Optin      *string `json:"optin,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Age        *string `json:"age,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	ExternalID *string `json:"external_id,omitempty"`
	DataMonth  *string `json:"data_month,omitempty"`
	HaveDSL    *string `json:"have_dsl,omitempty"`
	WhenHiring *string `json:"when_hiring,omitempty"`
}

// TableName sets the default table name
func (Virgin) TableName() string {
	return "virgin"
}
