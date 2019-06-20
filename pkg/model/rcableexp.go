package leads

import "github.com/jinzhu/gorm"

//RcableExp struct represents the fields of the rcableexp table
type RcableExp struct {
	gorm.Model

	LeadID    uint    `sql:"column:lea_id"`
	Location  *string `json:"location,omitempty"`
	Answer    *string `json:"answer,omitempty"`
	Values    *string `json:"values,omitempty"`
	Coverture *string `json:"coverture,omitempty"`
}

// TableName sets the default table name
func (RcableExp) TableName() string {
	return "rcableexp"
}
