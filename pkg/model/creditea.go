package leads

import "github.com/jinzhu/gorm"

// Creditea struct represents the fields of the creditea table
type Creditea struct {
	gorm.Model

	LeadID          uint    `sql:"column:lea_id"`
	RequestedAmount *string `json:"requested_amount,omitempty"`
	ContractType    *string `json:"contract_type,omitempty"`
	NetIncome       *string `json:"net_income,omitempty"`
	OutOfSchedule   *string `json:"out_of_schedule,omitempty"`
	ASNEF           bool    `json:"asnef"`
	AlreadyClient   bool    `json:"already_client"`
}

// TableName sets the default table name
func (Creditea) TableName() string {
	return "creditea"
}
