package leads

import "github.com/jinzhu/gorm"

// Creditea struct represents the fields of the creditea table
type Creditea struct {
	gorm.Model

	LeadID             uint    `sql:"column:lea_id"`
	Cantidadsolicitada *string `json:"cantidadsolicitada"`
	Motivo             *string `json:"motivo"`
	Tipocontrato       *string `json:"tipocontrato"`
	Ingresosnetos      *string `json:"ingresosnetos"`
}
