package leads

import "github.com/jinzhu/gorm"

// Alterna represents the data structure for Alterna project
type Alterna struct {
	gorm.Model

	LeadID      uint    `sql:"column:lea_id"`
	InstallType *string `json:"install_type,omitempty"`
	CPUS        *string `json:"cpus,omitempty"`
	Street      *string `json:"street,omitempty"`
	Number      *string `json:"number,omitempty"`
	PostalCode  *string `json:"postal_code,omitempty"`
}

// TableName sets the default table name
func (Alterna) TableName() string {
	return "alterna"
}
