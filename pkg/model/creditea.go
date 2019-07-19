package leads

import "github.com/jinzhu/gorm"

// Creditea struct represents the fields of the creditea table
type Creditea struct {
	gorm.Model

	LeadID             uint    `sql:"column:lea_id"`
	Cantidadsolicitada *string `json:"cantidadsolicitada,omitempty"`
	Tipocontrato       *string `json:"tipocontrato,omitempty"`
	Ingresosnetos      *string `json:"ingresosnetos,omitempty"`
	Fuerahorario       *string `json:"fuerahorario,omitempty"`
	Asnef              bool    `json:"asnef"`
	Yacliente          bool    `json:"yacliente"`

	// Validacionlp *string `json:"validacionlp,omitempty"`
	// Motivo       *string `json:"motivo"`
}

// TableName sets the default table name
func (Creditea) TableName() string {
	return "creditea"
}
