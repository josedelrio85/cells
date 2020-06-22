package leads

// Evolution represents data structure needed in SC
type Evolution struct {
	Properties     Properties     `json:"propiedades"`
	AdditionalData AdditionalData `json:"datosAdicionales,omitempty"`
	Localizators   Localizators   `json:"localizadores,omitempty"`
}

// Properties bla
type Properties struct {
	SubjectID       string `json:"idsujeto,omitempty"`
	OriginalID      string `json:"idoriginal"`
	CampaingID      string `json:"idcampanya"`
	Name            string `json:"nombre,omitempty"`
	Surname         string `json:"apellido,omitempty"`
	Surname2        string `json:"apellido2,omitempty"`
	Phone           string `json:"telefono"`
	Phone2          string `json:"telefono2,omitempty"`
	PhoneWork       string `json:"telefonoTrabajo,omitempty"`
	MobilePhone     string `json:"movil,omitempty"`
	MobilePhone2    string `json:"movil2,omitempty"`
	Address         string `json:"direccion,omitempty"`
	PostalCode      string `json:"codigoPostal,omitempty"`
	Town            string `json:"poblacion,omitempty"`
	State           string `json:"provincia,omitempty"`
	Country         string `json:"pais,omitempty"`
	Fax             string `json:"fax,omitempty"`
	Email           string `json:"email,omitempty"`
	Email2          string `json:"emaiL2,omitempty"`
	BirthDate       string `json:"fechA_NACIMIENTO,omitempty"`
	SignupDate      string `json:"fechA_ALTA,omitempty"`
	LanguageID      string `json:"iD_IDIOMA,omitempty"`
	Observations    string `json:"observaciones,omitempty"`
	LocatableSince  string `json:"localizablE_DESDE,omitempty"`
	LocatableFrom   string `json:"localizablE_HASTA,omitempty"`
	DNI             string `json:"sDNI,omitempty"`
	FullName        string `json:"sNombre_Completo,omitempty"`
	Company         string `json:"sEmpresa,omitempty"`
	Sex             string `json:"cSexo,omitempty"`
	Text1           string `json:"textO1,omitempty"`
	Text2           string `json:"textO2,omitempty"`
	Text3           string `json:"textO3,omitempty"`
	FavSource       string `json:"nCanalPreferencial,omitempty"`
	Num1            string `json:"nuM1,omitempty"`
	Num2            string `json:"nuM2,omitempty"`
	Num3            string `json:"nuM3,omitempty"`
	SegmentAttibute string `json:"atributo_Segmento,omitempty"`
	Priority        string `json:"prioridad,omitempty"`
	NextContact     string `json:"tProximo_Contacto,omitempty"`
	Skill           string `json:"atributo_Skill,omitempty"`
	NState          string `json:"nEstado,omitempty"`
	NList           string `json:"nLista,omitempty"`
}

// AdditionalData bla
type AdditionalData struct {
	AddProp1 string `json:"additionalProp1,omitempty"`
	AddProp2 string `json:"additionalProp2,omitempty"`
	AddProp3 string `json:"additionalProp3,omitempty"`
}

// Localizators bla
type Localizators struct {
	AddProp1 string `json:"additionalProp1,omitempty"`
	AddProp2 string `json:"additionalProp2,omitempty"`
	AddProp3 string `json:"additionalProp3,omitempty"`
}
