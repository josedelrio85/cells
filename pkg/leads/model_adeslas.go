package leads

import "github.com/jinzhu/gorm"

// Adeslas struct represents the fields of the Adeslas table
type Adeslas struct {
	gorm.Model

	LeadID  uint    `sql:"column:lea_id"`
	Product *string `json:"product,omitempty"`
	Landing *string `json:"landing,omitempty"`
}

// TableName sets the default table name
func (Adeslas) TableName() string {
	return "adeslas"
}
