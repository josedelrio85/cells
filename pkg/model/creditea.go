package leads

import "github.com/jinzhu/gorm"

// Creditea struct represents the fields of the creditea table
type Creditea struct {
	gorm.Model

	LeadID             uint    `sql:"column:lea_id"`
	Cantidadsolicitada *string `json:"cantidadsolicitada,omitempty"`
	Motivo             *string `json:"motivo"`
	Tipocontrato       *string `json:"tipocontrato,omitempty"`
	Ingresosnetos      *string `json:"ingresosnetos,omitempty"`
	Validacionlp       *string `json:"validacionlp,omitempty"`
	Fuerahorario       *string `json:"fuerahorario,omitempty"`
}

// TableName sets the default table name
func (Creditea) TableName() string {
	return "creditea"
}
