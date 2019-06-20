package leads

import "github.com/jinzhu/gorm"

// Microsoft struct represents the fields of the microsoft table
type Microsoft struct {
	gorm.Model

	LeadID        uint    `sql:"column:lea_id"`
	Tipoordenador *string `json:"tipoordenador"`
	Sector        *string `json:"sector"`
	Presupuesto   *string `json:"presupuesto"`
	Rendimiento   *string `json:"rendimiento"`
	Movilidad     *string `json:"movilidad"`
	Office365     *string `json:"office365"`
	Observaciones *string `json:"observaciones"`

	Producttype        *string `json:"producttype"`
	Productname        *string `json:"productname"`
	ProductID          *string `json:"productid"`
	Originalprice      *string `json:"originalprice"`
	Price              *string `json:"price"`
	Brand              *string `json:"brand"`
	Discountpercentage *string `json:"discountPercentage"`
	Discountcode       *string `json:"discountCode"`
	Typeofprocessor    *string `json:"typeofprocessor"`
	Harddiskcapacity   *string `json:"harddiskcapacity"`
	Graphics           *string `json:"graphics"`
	Wirelessinterface  *string `json:"wirelessinterface"`

	Anosordenadoresmedia        *string `json:"anos_ordenadores_media"`
	SistemaOperativoInstalado   *string `json:"sistema_operativo_instalado"`
	FrecuenciaBloqueOrdenadores *string `json:"frecuencia_bloqueo_ordenadores"`
	NumeroDispositivosEmpresa   *string `json:"num_dispositivos_empresa"`
	ReparacionesUltimoAno       *string `json:"reparaciones_ultimo_ano"`
	TiempoArrancarDispositivos  *string `json:"tiempo_arrancar_dispositivos"`

	Index bool `json:"index"`
}

// TableName sets the default table name
func (Microsoft) TableName() string {
	return "microsoft"
}

func (lead *Lead) beforeIncidence() (int64, int64) {
	// this function is a rememberance of the assignations
	// of sou_id's in function of domain and gclid values
	// it was used before an incidence ocurred in Leontel
	// environment
	var tipo int64
	var souid int64
	switch domain := *lead.Domain; domain {
	case "microsoftbusiness.es":
		tipo = 1
		souid = 49
	case "microsoftprofesional.es":
		tipo = 2
		souid = 50
		if !lead.Microsoft.Index {
			tipo = 3
			souid = 51
		}
	case "ofertas.mundo-r.com":
		tipo = 4
		souid = 25
	case "microsoftbusiness.es/hazelcambio":
		tipo = 5
		souid = 46
	case "microsoftnegocios.es":
		tipo = 6
		souid = 48
	}
	return tipo, souid
}
