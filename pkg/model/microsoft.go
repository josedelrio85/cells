package leads

import "github.com/jinzhu/gorm"

// Microsoft struct represents the fields of the microsoft table
type Microsoft struct {
	gorm.Model

	LeadID       uint    `sql:"column:lea_id"`
	ComputerType *string `json:"computer_type"`
	Sector       *string `json:"sector"`
	Budget       *string `json:"budget"`
	Performance  *string `json:"performance"`
	Movility     *string `json:"movility"`
	Office365    *string `json:"office365"`
	Usecase      *string `json:"usecase"`
	Comments     *string `json:"comments"`

	ProductType        *string `json:"product_type"`
	ProductName        *string `json:"product_name"`
	ProductID          *string `json:"product_id"`
	OriginalPrice      *string `json:"original_price"`
	Price              *string `json:"price"`
	Brand              *string `json:"brand"`
	DiscountPercentage *string `json:"discount_percentage"`
	DiscountCode       *string `json:"discount_code"`
	ProcessorType      *string `json:"processor_type"`
	DiskCapacity       *string `json:"disk_capacity"`
	Graphics           *string `json:"graphics"`
	WirelessInterface  *string `json:"wireless_interface"`

	DevicesAverageAge      *string `json:"devices_average_age"`
	DevicesOperatingSystem *string `json:"devices_operating_system"`
	DevicesHangFrequency   *string `json:"devices_hang_frequency"`
	DevicesNumber          *string `json:"devices_number"`
	DevicesLastYearRepairs *string `json:"devices_last_year_repairs"`
	DevicesStartupTime     *string `json:"devices_startup_time"`

	Pageindex bool  `json:"index"`
	Oldsouid  int64 `json:"oldsouid"`
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
		if !lead.Microsoft.Pageindex {
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
