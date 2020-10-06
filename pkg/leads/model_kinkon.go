package leads

import "github.com/jinzhu/gorm"

// Kinkon represents the data structure for Kinkon project
type Kinkon struct {
	gorm.Model

	LeadID      uint        `sql:"column:lea_id"`
	Coverture   *string     `json:"coverture,omitempty"`
	CovData     CovData     `gorm:"embedded" json:"cov_data"`
	Portability Portability `gorm:"embedded" json:"portability"`
	HolderData  HolderData  `gorm:"embedded" json:"holder_data"`
	BillingInfo BillingInfo `gorm:"embedded" json:"billing_info"`
	Product     *string     `json:"product,omitempty"`
	Mvf         Mvf         `gorm:"embedded" json:"mvf"`
	Ignium      Ignium      `gorm:"embedded" json:"ignium"`
}

// CovData represents the data structure for coverture data
type CovData struct {
	State    *string `json:"state,omitempty"`
	Town     *string `json:"town,omitempty"`
	Street   *string `json:"street,omitempty"`
	Number   *string `json:"number,omitempty"`
	Floor    *string `json:"floor,omitempty"`
	CovPhone *string `json:"phone,omitempty"`
}

// Portability represents the data structure for portability data
type Portability struct {
	Phone                *string `json:"phone,omitempty"`
	PhoneProvider        *string `json:"phone_provider,omitempty"`
	MobilePhone          *string `json:"mobile_phone,omitempty"`
	MobilePhoneProvider  *string `json:"mobile_phone_provider,omitempty"`
	MobilePhone2         *string `json:"mobile_phone_2,omitempty"`
	MobilePhoneProvider2 *string `json:"mobile_phone_provider_2,omitempty"`
}

// HolderData represents the data structure for holder data
type HolderData struct {
	Name         *string `json:"name,omitempty"`
	Surname      *string `json:"surname,omitempty"`
	Idnumber     *string `json:"idnumber,omitempty"`
	Mail         *string `json:"mail,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
}

// BillingInfo represents the data structure for billing info data
type BillingInfo struct {
	AccountHolder *string `json:"account_holder,omitempty"`
	AccountNumber *string `json:"account_number,omitempty"`
}

// Mvf represents the data structure for mvf lead provider
type Mvf struct {
	LeadReferenceNumber      *string `sql:"column:lead_reference_number_mvf" json:"lrf,omitempty"`
	DistributionID           *string `sql:"column:distribution_id_mvf" json:"dist_id,omitempty"`
	HasSwitchboard           *string `sql:"column:has_switchboard_mvf" json:"has_switchboard,omitempty"`
	ExtensionsNumber         *string `sql:"column:extensions_number_mvf" json:"extensions_number,omitempty"`
	PhoneAmount              *string `sql:"column:phone_amount_mvf" json:"phone_amount,omitempty"`
	EmployeeNumber           *string `sql:"column:employee_number_mvf" json:"employee_number,omitempty"`
	SwitchboardFunctionality *string `sql:"column:switchboard_func_mvf" json:"switchboard_functionality,omitempty"`
	Surname                  *string `sql:"column:surname_mvf" json:"surname,omitempty"`
	PostalCode               *string `sql:"column:postal_code_mvf" json:"postal_code,omitempty"`
}

// Ignium represents the data structure for Ignium lead provider
type Ignium struct {
	Optin         string  `sql:"optin_ignium" json:"optin_ignium,omitempty"`
	PostalCode    *string `sql:"postal_code_ignium" json:"postal_code_ignium,omitempty"`
	Age           *string `sql:"age_ignium" json:"age_ignium,omitempty"`
	Surname       *string `sql:"surname_ignium" json:"surname_ignium,omitempty"`
	ExternalID    *string `sql:"external_id_ignium" json:"external_id_ignium,omitempty"`
	DataMonth     *string `sql:"data_month_ignium" json:"data_month_ignium,omitempty"`
	HaveDSL       *string `sql:"have_dsl_ignium" json:"have_dsl_ignium,omitempty"`
	WhenHiring    *string `sql:"when_hiring_ignium" json:"when_hiring_ignium,omitempty"`
	ContacTime    *string `sql:"contact_time_ignium" json:"contact_time_ignium,omitempty"`
	Permanence    *string `sql:"permanence_ignium" json:"permanence_ignium,omitempty"`
	ActualCompany *string `sql:"actual_company_ignium" json:"actual_company_ignium,omitempty"`
	Rate          *string `sql:"rate_ignium" json:"rate_ignium,omitempty"`
}

// TableName sets the default table name
func (Kinkon) TableName() string {
	return "kinkon"
}
