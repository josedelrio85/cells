package leads

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

// Lead struct represents the fields of the leads table
type Lead struct {
	gorm.Model
	LegacyID           int64      `json:"-"`
	LeaTs              *time.Time `sql:"DEFAULT:current_timestamp" json:"-" `
	LeaSmartcenterID   string     `json:"-"`
	PassportID         string     `json:"passport_id"`
	PassportIDGrp      string     `json:"passport_id_group"`
	SouID              int64      `json:"sou_id"`
	LeatypeID          int64      `json:"lea_type"`
	UtmSource          *string    `json:"utm_source,omitempty"`
	SubSource          *string    `json:"sub_source,omitempty"`
	LeaPhone           *string    `json:"phone,omitempty"`
	LeaMail            *string    `json:"mail,omitempty"`
	LeaName            *string    `json:"name,omitempty"`
	LeaDNI             *string    `json:"dni,omitempty"`
	LeaURL             *string    `sql:"type:text" json:"url,omitempty"`
	LeaIP              *string    `json:"ip,omitempty"`
	GaClientID         *string    `json:"ga_client_id,omitempty"`
	IsSmartCenter      bool       `json:"smartcenter,omitempty"`
	SouIDLeontel       int64      `sql:"-" json:"sou_id_leontel"`
	SouDescLeontel     string     `sql:"-" json:"sou_desc_leontel"`
	LeatypeIDLeontel   int64      `sql:"-" json:"lea_type_leontel"`
	LeatypeDescLeontel string     `sql:"-" json:"lea_type_desc_leontel"`
	Gclid              *string    `json:"gclid,omitempty"`
	Domain             *string    `json:"domain,omitempty"`
	Observations       *string    `sql:"type:text" json:"observations,omitempty"`
	RcableExp          *RcableExp `json:"rcableexp"`
	Microsoft          *Microsoft `json:"microsoft"`
	Creditea           *Creditea  `json:"creditea"`
	Kinkon             *Kinkon    `json:"kinkon"`
	Alterna            *Alterna   `json:"alterna"`
}

// TableName sets the default table name
func (Lead) TableName() string {
	return "leads"
}

// Decode reques's body into a Lead struct
func (lead *Lead) Decode(body io.ReadCloser) error {
	if err := json.NewDecoder(body).Decode(lead); err != nil {
		return err
	}
	return nil
}

// LeadToLeontel maps values from Lead struct (webservice.leads table)
// as an appropiate lead Leontel input for each campaign
func (lead *Lead) LeadToLeontel() LeadLeontel {
	leontel := LeadLeontel{
		LeaSource: lead.SouIDLeontel,
		LeaType:   lead.LeatypeIDLeontel,
		Telefono:  lead.LeaPhone,
		Nombre:    lead.LeaName,
		URL:       lead.LeaURL,
		IP:        lead.LeaIP,
		Email:     lead.LeaMail,
		Dninie:    lead.LeaDNI,
		Wsid:      lead.ID,
		// TODO set GaClientID value into LeadLeontel struct
	}

	switch souid := lead.SouID; souid {
	case 1, 21, 22:
		// Creditea Abandonos
		// lea_aux1 => 1 | 0 ??
		// lea_surname => apellido1
		// lea_aux2 => dninie
		// lea_aux3 => asnef
		// En teoría se recogen name, surname, dninie, asnef pero en BD solo hay lea_aux1 = 1 | 0
	case 2:
		//Creditea Stand
		// lea_surname => apellido1
		// lea_aux2 => dninie
		// lea_aux3 => asnef
	case 4:
		// Creditea Timeout
	case 5:
	case 14:
		// R Cable + R Cable Empresas
	case 7:
		// Hercules
	case 8:
		// Seguro para movil
	case 9, 58:
		// Creditea EndToEnd + CREDITEA HM CORTO
		// C2C_Creditea_validaDNI_telf =>	lea_aux2 => cantidadsolicitada || lea_aux1 => dninie
		// almacenaLeadNoValido => lea_aux2 => cantidadsolicitada || lea_aux1 => dninie || lea_aux3 => motivo
		// HM Creditea (webpack) => lea_aux2 => cantidadsolicitada || lea_aux1 => dninie || lea_aux3 => motivo
		// leontel observaciones => (DNI: $dninie Cantidad solicitada: $cantidadSolicitada)
		args := []*string{
			lead.LeaDNI,
			lead.Creditea.RequestedAmount,
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones = &observations
	case 10:
		// Creditea FB
	case 11:
		// Creditea Rastreator (no activo)
		// lea_aux1 (cantidad solicitada)
		// lea_aux2 (Tipo contrato)
		// lea_aux3 (ingresos netos)
		// lea_aux4 (dni) => dninie
		args := []*string{
			lead.LeaDNI,
			lead.Creditea.NetIncome,
			lead.Creditea.ContractType,
			lead.Creditea.RequestedAmount,
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones = &observations
	case 12:
		// Euskaltel
	case 13:
		// Adeslas
	case 17, 18, 19, 20:
		// Yoigo
		// producto, cobertura, impuesto
		// lea_aux2 => producto
		// lea_aux3 => cobertura // impuesto
		// leontel.Observaciones = producto
		// leontel.Observaciones2 = cobertura // impuesto
	case 23:
		// Creditea pago recurrente (Cron, tratamiento diferente, NO API)
	case 24:
		// Sanal
	case 25:
		// Microsoft Mundo R
	case 46, 49:
		// Microsoft Hazelcambio + Recomendador
		leontel.Tipoordenador = lead.Microsoft.ComputerType
		leontel.Sector = lead.Microsoft.Sector
		leontel.Tipouso = lead.Microsoft.Usecase

		leontel.Presupuesto = lead.Microsoft.Budget
		leontel.Rendimiento = lead.Microsoft.Performance
		leontel.Movilidad = lead.Microsoft.Movility
		leontel.Office365 = lead.Microsoft.Office365
		leontel.Observaciones2 = lead.Observations

		// Setear Microsoft Global
		lead.Microsoft.Oldsouid = lead.SouID
		lead.SouID = 52
		lead.SouIDLeontel = 61
		leontel.LeaSource = lead.SouIDLeontel
	case 48:
		// Microsoft Calculadora
		// => observaciones2
		// anos_ordenadores_media: {$anos_ordenadores_media}
		// sistema_operativo_instalado: {$sistema_operativo_instalado}
		// frecuencia_bloqueo_ordenadores: {$frecuencia_bloqueo_ordenadores}
		// num_dispositivos_empresa: {$num_dispositivos_empresa}
		// reparaciones_ultimo_ano: {$reparaciones_ultimo_ano}
		// tiempo_arrancar_dispositivos: {$tiempo_arrancar_dispositivos}

		args := []*string{
			lead.Microsoft.DevicesAverageAge,
			lead.Microsoft.DevicesOperatingSystem,
			lead.Microsoft.DevicesHangFrequency,
			lead.Microsoft.DevicesNumber,
			lead.Microsoft.DevicesLastYearRepairs,
			lead.Microsoft.DevicesStartupTime,
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones2 = &observations

		// Setear Microsoft Global
		lead.Microsoft.Oldsouid = lead.SouID
		lead.SouID = 52
		lead.SouIDLeontel = 61
		leontel.LeaSource = lead.SouIDLeontel
	case 50:
		// Microsoft Ofertas
		// Setear Microsoft Global
		lead.Microsoft.Oldsouid = lead.SouID
		lead.SouID = 52
		lead.SouIDLeontel = 61
		leontel.LeaSource = lead.SouIDLeontel
	case 51:
		// Microsoft FichaProducto
		// => observaciones2

		// Tipo: {$productType}
		// Producto: {$name}
		// idProducto: {$id}
		// precioOriginal: {$originalPrice}
		// Precio: {$price}
		// Marca: {$brand}
		// %Descuento: {$discountPercentage}
		// Cod. descuento: {$discountCode}
		// Tipo Procesador: {$typeOfProcessor}
		// Capacidad HDD: {$hardDiskCapacity}
		// Gráfica: {$graphics}
		// Wireless: {$wirelessInterface}

		args := []*string{
			lead.Microsoft.ProductType,
			lead.Microsoft.ProductName,
			lead.Microsoft.ProductID,
			lead.Microsoft.OriginalPrice,
			lead.Microsoft.Price,
			lead.Microsoft.Brand,
			lead.Microsoft.DiscountPercentage,
			lead.Microsoft.DiscountCode,
			lead.Microsoft.ProcessorType,
			lead.Microsoft.DiskCapacity,
			lead.Microsoft.Graphics,
			lead.Microsoft.WirelessInterface,
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones2 = &observations

		// Setear Microsoft Global
		lead.Microsoft.Oldsouid = lead.SouID
		lead.SouID = 52
		lead.SouIDLeontel = 61
		leontel.LeaSource = lead.SouIDLeontel
	case 53:
		// IPF
		// "lea_aux1" => $lead->nameId,
		// "lea_aux2" => $lead->productAmountTaken,
		// "lea_aux4" => $lead->clientId,
		// "observations" => $observations
	case 54:
		// R Cable Expansion
		args := []*string{lead.LeaName, lead.RcableExp.RespValues, lead.RcableExp.Location, lead.RcableExp.Answer}
		observations := concatPointerStrs(args...)
		lead.Observations = &observations
		leontel.Observaciones = &observations
	case 55:
		// R Cable Expansion Entrante
	case 56:
		// Creditea BO
	case 57:
		// Sanitas
		// lea_destiny =  GSS => we must have IsLeontel = true
	case 64, 65, 66:
		args := []*string{}

		if lead.Kinkon != nil {
			coverture := "Cobertura"
			args = append(args, &coverture)
			args = append(args, lead.Kinkon.Coverture)

			product := "Producto"
			args = append(args, &product)
			args = append(args, lead.Kinkon.Product)

			if *lead.Kinkon.CovData != (CovData{}) {
				leontel.Provincia = lead.Kinkon.CovData.State
				leontel.Poblacion = lead.Kinkon.CovData.Town
				cargs := []*string{}
				cargs = append(cargs, lead.Kinkon.CovData.Street)
				cargs = append(cargs, lead.Kinkon.CovData.Number)
				cargs = append(cargs, lead.Kinkon.CovData.Floor)
				covargs := concatPointerStrs(cargs...)
				leontel.Direccion = &covargs
			}

			if *lead.Kinkon.Portability != (Portability{}) {
				leontel.Compaiaactualfibraadsl = lead.Kinkon.Portability.PhoneProvider
				leontel.Companiaactualmovil = lead.Kinkon.Portability.MobilePhoneProvider

				phone := "Teléfono fijo portabilidad:"
				args = append(args, &phone)
				args = append(args, lead.Kinkon.Portability.Phone)

				mobile := "Teléfono movil portabilidad:"
				args = append(args, &mobile)
				args = append(args, lead.Kinkon.Portability.MobilePhone)

				phone2 := "Teléfono movil 2 portabilidad:"
				args = append(args, &phone2)
				args = append(args, lead.Kinkon.Portability.MobilePhone2)

				provider := "Operador movil portabilidad:"
				args = append(args, &provider)
				args = append(args, lead.Kinkon.Portability.MobilePhoneProvider2)
			}

			if *lead.Kinkon.HolderData != (HolderData{}) {
				fullname := fmt.Sprintf("%s %s", *lead.Kinkon.HolderData.Name, *lead.Kinkon.HolderData.Surname)
				leontel.Nombrecompleto = &fullname
				leontel.Dninie = lead.Kinkon.HolderData.Idnumber
				leontel.Email = lead.Kinkon.HolderData.Mail

				contactphone := "Teléfono contacto"
				args = append(args, &contactphone)
				args = append(args, lead.Kinkon.HolderData.ContactPhone)
			}

			if *lead.Kinkon.BillingInfo != (BillingInfo{}) {
				accountholder := "Titular cuenta"
				args = append(args, &accountholder)
				args = append(args, lead.Kinkon.BillingInfo.AccountHolder)

				ccc := "CCC"
				args = append(args, &ccc)
				args = append(args, lead.Kinkon.BillingInfo.AccountNumber)
			}

			observations := concatPointerStrs(args...)
			leontel.Observaciones = &observations
		}

	case 69:
		if lead.Alterna != nil {
			leontel.CP = lead.Alterna.PostalCode
			leontel.Calle = lead.Alterna.Street
			leontel.Numero = lead.Alterna.Number
			leontel.Tiposolicitud = lead.Alterna.InstallType
			leontel.Observaciones = lead.Observations
			if lead.Alterna.CPUS != nil {
				leontel.Observaciones = lead.Alterna.CPUS
			}
		}
	case 70:
		leontel.Observaciones = lead.Observations

	default:

	}
	return leontel
}

// SendLeadToLeontel sends the lead to Leontel endpoint
// Returns the response sended by the endpoint
func (lead *Lead) SendLeadToLeontel() (*LeontelResp, error) {
	leadLeontel := lead.LeadToLeontel()
	bytevalues, err := json.Marshal(leadLeontel)
	if err != nil {
		return nil, err
	}

	endpoint, ok := os.LookupEnv("LEAD_LEONTEL_ENDPOINT")
	if !ok {
		err := errors.New("unable to load Lead Leontel URL endpoint")
		return nil, err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	leontelresp := []LeontelResp{}
	if err := json.Unmarshal(data, &leontelresp); err != nil {
		return nil, err
	}
	return &leontelresp[0], nil
}

// GetLeontelValues queries for the Leontel equivalences
// of sou_id and lea_type values
func (lead *Lead) GetLeontelValues(db *gorm.DB) error {
	source := Source{}
	leatype := Leatype{}

	if result := db.Where("sou_id = ?", lead.SouID).First(&source); result.Error != nil {
		return fmt.Errorf("Error retrieving SouIDLeontel value: %#v", result.Error)
	}
	if result := db.Where("leatype_id = ?", lead.LeatypeID).First(&leatype); result.Error != nil {
		return fmt.Errorf("error retrieving LeatypeIDLeontel value: %#v", result.Error)
	}
	lead.SouIDLeontel = source.SouIdcrm
	lead.SouDescLeontel = source.SouDescription
	lead.LeatypeIDLeontel = leatype.LeatypeIdcrm
	lead.LeatypeDescLeontel = leatype.LeatypeDescription
	return nil
}

// GetParams retrieves values for ip, port and url properties
func (lead *Lead) GetParams(w http.ResponseWriter, req *http.Request) error {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return fmt.Errorf("Error parsing IP value: %#v", err)
	}
	// forward := req.Header.Get("X-Forwarded-For")

	lead.LeaIP = &ip
	url := fmt.Sprintf("%s%s", req.Host, req.URL.Path)
	lead.LeaURL = &url

	return nil
}

// concatPointerStrs concats an undefined number
// of *string params into string separated by --
func concatPointerStrs(args ...*string) string {
	var buffer bytes.Buffer
	tam := len(args) - 1
	for i, arg := range args {
		if arg != nil {
			buffer.WriteString(*arg)
			if i < tam {
				buffer.WriteString(" -- ")
			}
		}
	}
	return buffer.String()
}

// GetPassport gets a passport and sets it into lead properties
func (lead *Lead) GetPassport() error {
	passport := Passport{}

	interaction := Interaction{
		Provider:    lead.SouDescLeontel,
		Application: lead.LeatypeDescLeontel,
		IP:          *lead.LeaIP,
	}

	if err := passport.Get(interaction); err != nil {
		return err
	}

	lead.PassportID = passport.PassportID
	lead.PassportIDGrp = passport.PassportIDGrp

	return nil
}
