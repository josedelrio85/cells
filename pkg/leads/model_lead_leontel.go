package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

// LeadLeontel represents the data structure of lea_leads Leontel table
type LeadLeontel struct {
	LeaID                    int64      `json:"-"`
	LeaType                  int64      `json:"lea_type,omitempty"`
	LeaTs                    *time.Time `json:"-"`
	LeaSource                int64      `json:"lea_source,omitempty"`
	LeaLot                   int64      `json:"-"`
	LeaAssigned              int64      `json:"-"`
	LeaScheduled             *time.Time `json:"lea_scheduled,omitempty"`
	LeaScheduledAuto         int64      `json:"lea_scheduled_auto,omitempty"`
	LeaCost                  float64    `json:"-"`
	LeaSeen                  int64      `json:"-"`
	LeaClosed                int64      `json:"lea_closed,omitempty"`
	LeaNew                   int64      `json:"-"`
	LeaClient                int64      `json:"-"`
	LeaDateControl           *time.Time `json:"-"`
	Telefono                 *string    `json:"TELEFONO,omitempty"`
	Nombre                   *string    `json:"nombre,omitempty"`
	Apellido1                *string    `json:"apellido1,omitempty"`
	Apellido2                *string    `json:"apellido2,omitempty"`
	Dninie                   *string    `json:"dninie,omitempty"`
	Observaciones            *string    `json:"observaciones,omitempty"`
	URL                      *string    `json:"url,omitempty"`
	Asnef                    *string    `json:"asnef,omitempty"`
	Wsid                     uint       `json:"wsid,omitempty"`
	IP                       *string    `json:"ip,omitempty"`
	Email                    *string    `json:"Email,omitempty"`
	Observaciones2           *string    `json:"observaciones2,omitempty"`
	Hashid                   *string    `json:"hashid,omitempty"`
	Nombrecompleto           *string    `json:"nombrecompleto,omitempty"`
	Movil                    *string    `json:"movil,omitempty"`
	Tipocliente              *string    `json:"tipocliente,omitempty"`
	Escliente                *string    `json:"escliente,omitempty"`
	Tiposolicitud            *string    `json:"tiposolicitud,omitempty"`
	Fechasolicitud           *time.Time `json:"fechasolicitud,omitempty"`
	Poblacion                *string    `json:"poblacion,omitempty"`
	Provincia                *string    `json:"provincia,omitempty"`
	Direccion                *string    `json:"direccion,omitempty"`
	Cargo                    *string    `json:"cargo,omitempty"`
	Cantidaddeseada          *string    `json:"cantidaddeseada,omitempty"`
	Cantidadofrecida         *string    `json:"cantidadofrecida,omitempty"`
	Ncliente                 *string    `json:"ncliente,omitempty"`
	Calle                    *string    `json:"calle,omitempty"`
	CP                       *string    `json:"cp,omitempty"`
	Interesadoen             *string    `json:"interesadoen,omitempty"`
	Numero                   *string    `json:"numero,omitempty"`
	Compaiaactualfibraadsl   *string    `json:"compaiaactualfibraadsl,omitempty"`
	Companiaactualmovil      *string    `json:"companiaactualmovil,omitempty"`
	Fibraactual              *string    `json:"fibraactual,omitempty"`
	Moviactuallineaprincipal *string    `json:"moviactuallineaprincipal,omitempty"`
	Numerolineasadicionales  int64      `json:"numerolineasadicionales,omitempty"`
	Tarifaactualsiniva       float64    `json:"tarifaactualsiniva,omitempty"`
	Motivocambio             *string    `json:"motivocambio,omitempty"`
	Importeaumentado         int64      `json:"importeaumentado,omitempty"`
	Importeretirado          int64      `json:"importeretirado,omitempty"`
	Tipoordenador            *string    `json:"tipoordenador,omitempty"`
	Sector                   *string    `json:"sector,omitempty"`
	Presupuesto              *string    `json:"presupuesto,omitempty"`
	Rendimiento              *string    `json:"rendimiento,omitempty"`
	Movilidad                *string    `json:"movilidad,omitempty"`
	Tipouso                  *string    `json:"tipouso,omitempty"`
	Office365                *string    `json:"Office365,omitempty"`
}

// LeontelResp represents a structure with the response returned
// by Leontel endpoint.
type LeontelResp struct {
	Success bool   `json:"success"`
	ID      int64  `json:"id"`
	Error   string `json:"error,omitempty"`
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
		Ncliente:  lead.GaClientID,
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
		// AdeslasOLD
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
	case 64, 65, 66, 74, 75, 76:
		// kinkon + r-empresas
		args := []*string{}
		args = append(args, lead.Observations)

		if lead.Kinkon != nil {
			coverture := "Cobertura"
			args = append(args, &coverture)
			args = append(args, lead.Kinkon.Coverture)

			product := "Producto"
			args = append(args, &product)
			args = append(args, lead.Kinkon.Product)

			if lead.Kinkon.CovData != (CovData{}) {
				leontel.Provincia = lead.Kinkon.CovData.State
				leontel.Poblacion = lead.Kinkon.CovData.Town
				cargs := []*string{}
				cargs = append(cargs, lead.Kinkon.CovData.Street)
				cargs = append(cargs, lead.Kinkon.CovData.Number)
				cargs = append(cargs, lead.Kinkon.CovData.Floor)
				covargs := concatPointerStrs(cargs...)
				leontel.Direccion = &covargs
			}

			if lead.Kinkon.Portability != (Portability{}) {
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

			if lead.Kinkon.HolderData != (HolderData{}) {
				fullname := fmt.Sprintf("%s %s", *lead.Kinkon.HolderData.Name, *lead.Kinkon.HolderData.Surname)
				leontel.Nombrecompleto = &fullname
				leontel.Dninie = lead.Kinkon.HolderData.Idnumber
				leontel.Email = lead.Kinkon.HolderData.Mail

				contactphone := "Teléfono contacto"
				args = append(args, &contactphone)
				args = append(args, lead.Kinkon.HolderData.ContactPhone)
			}

			if lead.Kinkon.BillingInfo != (BillingInfo{}) {
				accountholder := "Titular cuenta"
				args = append(args, &accountholder)
				args = append(args, lead.Kinkon.BillingInfo.AccountHolder)

				ccc := "CCC"
				args = append(args, &ccc)
				args = append(args, lead.Kinkon.BillingInfo.AccountNumber)
			}

			if lead.Kinkon.Mvf != (Mvf{}) {
				q1 := "Lead Reference Number"
				q2 := "Distribution ID"
				q3 := "¿Ya tiene una centralita telefónica?"
				q4 := "¿Cuantas extensiones necesita?"
				q5 := "Nº exacto de teléfonos"
				q6 := "¿Cuántos empleados tiene su empresa?"
				q7 := "¿Qué funcionalidad de centralita necesita?"
				q8 := "Apellidos"
				q9 := "Código Postal"

				args = append(args, &q1, lead.Kinkon.Mvf.LeadReferenceNumber)
				args = append(args, &q2, lead.Kinkon.Mvf.DistributionID)
				args = append(args, &q3, lead.Kinkon.Mvf.HasSwitchboard)
				args = append(args, &q4, lead.Kinkon.Mvf.ExtensionsNumber)
				args = append(args, &q5, lead.Kinkon.Mvf.PhoneAmount)
				args = append(args, &q6, lead.Kinkon.Mvf.EmployeeNumber)
				args = append(args, &q7, lead.Kinkon.Mvf.SwitchboardFunctionality)
				args = append(args, &q8, lead.Kinkon.Mvf.Surname)
				args = append(args, &q9, lead.Kinkon.Mvf.PostalCode)
			}

			if lead.Kinkon.Ignium != (Ignium{}) {
				q1 := "Optin"
				q2 := "Código postal"
				q3 := "Edad"
				q4 := "Apellidos"
				q5 := "External ID"
				q6 := "Datos al mes"
				q7 := "¿Tienes actualmente ADSL/Fibra?"
				q8 := "Cuando lo vas a contratar"
				q9 := "Hora preferida de contacto"
				q10 := "¿Tienes permanencia?"
				q11 := "¿De que compañia eres?"
				q12 := "Tarifa"

				args = append(args, &q1, lead.Kinkon.Ignium.Optin)
				args = append(args, &q2, lead.Kinkon.Ignium.PostalCode)
				args = append(args, &q3, lead.Kinkon.Ignium.Age)
				args = append(args, &q4, lead.Kinkon.Ignium.Surname)
				args = append(args, &q5, lead.Kinkon.Ignium.ExternalID)
				args = append(args, &q6, lead.Kinkon.Ignium.DataMonth)
				args = append(args, &q7, lead.Kinkon.Ignium.HaveDSL)
				args = append(args, &q8, lead.Kinkon.Ignium.WhenHiring)
				args = append(args, &q9, lead.Kinkon.Ignium.ContacTime)
				args = append(args, &q10, lead.Kinkon.Ignium.Permanence)
				args = append(args, &q11, lead.Kinkon.Ignium.ActualCompany)
				args = append(args, &q12, lead.Kinkon.Ignium.Rate)
			}
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones = &observations

	case 69:
		// Alterna
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
	case 77:
		// Adeslas
		args := []*string{}
		args = append(args, lead.Observations)
		if lead.Adeslas != nil {
			args = append(args, lead.Adeslas.Product)
			args = append(args, lead.Adeslas.Landing)
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones = &observations
	case 78:
		// Endesa
		args := []*string{}
		args = append(args, lead.Observations)

		if lead.Endesa != nil {
			q1 := "¿Qué tipo de energía tienes en tu hogar?"
			q2 := "¿Cuál es el tamaño de tu vivienda?"
			q3 := "¿Cuántas personas viven en casa?"
			q4 := "¿Qué tipo de energía usas en la calefacción?"
			q5 := "¿Qué tipo de energía usas en la en la cocina?"
			q6 := "¿Qué tipo de energía usas en el agua caliente?"
			q7 := "¿Cada cuanto pones la lavadora?"
			q8 := "¿Cada cuanto pones la secadora?"
			q9 := "¿Cada cuanto pones el lavavajillas?"
			q10 := "¿Eres el propietario de la vivienda?"
			q11 := "¿Cuál es tu compañía actual??"
			q12 := "Código postal"
			q13 := "Edad"
			q14 := "Apellidos"
			q15 := "External ID"

			args = append(args, &q14, lead.Endesa.Surname)
			args = append(args, &q1, lead.Endesa.TypeEnergy)
			args = append(args, &q2, lead.Endesa.HomeSize)
			args = append(args, &q3, lead.Endesa.HomePopulation)
			args = append(args, &q4, lead.Endesa.TypeHeating)
			args = append(args, &q5, lead.Endesa.TypeKitchen)
			args = append(args, &q6, lead.Endesa.TypeWater)
			args = append(args, &q7, lead.Endesa.WashingMachine)
			args = append(args, &q8, lead.Endesa.Dryer)
			args = append(args, &q9, lead.Endesa.Dishwasher)
			args = append(args, &q10, lead.Endesa.Owner)
			args = append(args, &q11, lead.Endesa.Company)
			args = append(args, &q12, lead.Endesa.PostalCode)
			args = append(args, &q13, lead.Endesa.Age)
			args = append(args, &q15, lead.Endesa.ExternalID)
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones = &observations
	case 79:
		// Virgin TODO delete when evolution is ready
		args := []*string{}
		args = append(args, lead.Observations)

		if lead.Virgin != nil {
			q1 := "Optin"
			q2 := "Código postal"
			q3 := "Edad"
			q4 := "Apellidos"
			q5 := "External ID"
			q6 := "Datos al mes"
			q7 := "¿Tienes actualmente ADSL/Fibra?"
			q8 := "Cuando lo vayas a contratar"
			q9 := "Hora preferida de contacto"
			q10 := "¿Tienes permanencia?"
			q11 := "¿De que compañia eres?"

			args = append(args, &q1, lead.Virgin.Optin)
			args = append(args, &q2, lead.Virgin.PostalCode)
			args = append(args, &q3, lead.Virgin.Age)
			args = append(args, &q4, lead.Virgin.Surname)
			args = append(args, &q5, lead.Virgin.ExternalID)
			args = append(args, &q6, lead.Virgin.DataMonth)
			args = append(args, &q7, lead.Virgin.HaveDSL)
			args = append(args, &q8, lead.Virgin.WhenHiring)
			args = append(args, &q9, lead.Virgin.ContacTime)
			args = append(args, &q10, lead.Virgin.Permanence)
			args = append(args, &q11, lead.Virgin.ActualCompany)
		}
		observations := concatPointerStrs(args...)
		leontel.Observaciones = &observations
	default:
	}
	return leontel
}

// Active is an implementation of Active method from Scable interface
func (ll LeadLeontel) Active(lead Lead, dev bool) bool {
	// TODO (delete) keep this hack to use virgin as Leontel campaign in pro environment
	if !dev {
		log.Printf("TEMPORAL HACK LEONTEL-EVOLUTION: ALL CAMPAIGNS ARE LEONTEL WHEN DEV IS %t", dev)
		return true
	}
	// (TODO for now, discard 79 (virgin) ||
	if lead.SouID != 79 {
		log.Printf("souid %d Leontel active", lead.SouID)
		return true
	}
	return false
}

// Send is an implementation of Send method from Scable interface
func (ll LeadLeontel) Send(lead Lead) ScResponse {
	leadLeontel := lead.LeadToLeontel()
	bytevalues, err := json.Marshal(leadLeontel)
	if err != nil {
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}

	endpoint, ok := os.LookupEnv("LEAD_LEONTEL_ENDPOINT")
	if !ok {
		err := errors.New("unable to load Lead Leontel URL endpoint")
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
	if err != nil {
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	leontelresp := []LeontelResp{}
	if err := json.Unmarshal(data, &leontelresp); err != nil {
		return ScResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			ID:         0,
			Error:      err,
		}
	}

	err = nil
	status := http.StatusOK
	if leontelresp[0].Error != "" {
		err = errors.New(leontelresp[0].Error)
		status = http.StatusUnprocessableEntity
	}

	return ScResponse{
		Success:    leontelresp[0].Success,
		StatusCode: status,
		ID:         leontelresp[0].ID,
		Error:      err,
	}
}
