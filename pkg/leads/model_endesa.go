package leads

import "github.com/jinzhu/gorm"

// Endesa struct represents the fields of the Adeslas table
type Endesa struct {
	gorm.Model

	LeadID         uint    `sql:"column:lea_id"`
	TypeEnergy     *string `json:"type_energy,omitempty"`
	HomeSize       *string `json:"home_size,omitempty"`
	HomePopulation *string `json:"home_population,omitempty"`
	TypeHeating    *string `json:"type_heating,omitempty"`
	TypeKitchen    *string `json:"type_kitchen,omitempty"`
	TypeWater      *string `json:"type_washer,omitempty"`
	WashingMachine *string `json:"washing_machine,omitempty"`
	Dryer          *string `json:"dryer,omitempty"`
	Dishwasher     *string `json:"dish_washer,omitempty"`
	Owner          *string `json:"owner,omitempty"`
	Company        *string `json:"company,omitempty"`
	PostalCode     *string `json:"postal_code,omitempty"`
	Age            *string `json:"age,omitempty"`
	Surname        *string `json:"surname,omitempty"`
	ExternalID     *string `json:"external_id,omitempty"`
}

// TableName sets the default table name
func (Endesa) TableName() string {
	return "endesa"
}
