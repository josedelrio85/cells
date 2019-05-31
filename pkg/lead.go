package apic2c

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// LeadLeontel represents the data structure of lea_leads Leontel table
type LeadLeontel struct {
	// lea_ fields
	LeaID            int64     `json:"-"`
	LeaType          int64     `json:"lea_type,omitempty"`
	LeaTs            time.Time `json:"-"`
	LeaSource        int64     `json:"lea_source,omitempty"`
	LeaLot           int64     `json:"-"`
	LeaAssigned      int64     `json:"-"`
	LeaScheduled     time.Time `json:"lea_scheduled,omitempty"`
	LeaScheduledAuto int64     `json:"lea_scheduled_auto,omitempty"`
	LeaCost          float64   `json:"-"`
	LeaSeen          int64     `json:"-"`
	LeaClosed        int64     `json:"lea_closed,omitempty"`
	LeaNew           int64     `json:"-"`
	LeaClient        int64     `json:"-"`
	LeaDateControl   time.Time `json:"-"`

	// most used fields
	Telefono       *string `json:"TELEFONO,omitempty"`
	Nombre         *string `json:"nombre,omitempty"`
	Apellido1      *string `json:"apellido1,omitempty"`
	Apellido2      *string `json:"apellido2,omitempty"`
	Dninie         *string `json:"dninie,omitempty"`
	Observaciones  *string `json:"observaciones,omitempty"`
	URL            *string `json:"url,omitempty"`
	Asnef          *string `json:"asnef,omitempty"`
	Wsid           int64   `json:"wsid,omitempty"`
	IP             *string `json:"ip,omitempty"`
	Email          *string `json:"Email,omitempty"`
	Observaciones2 *string `json:"observaciones2,omitempty"`

	Hashid                   *string   `json:"hashid,omitempty"`
	Nombrecompleto           *string   `json:"nombrecompleto,omitempty"`
	Movil                    *string   `json:"movil,omitempty"`
	Tipocliente              *string   `json:"tipocliente,omitempty"`
	Escliente                *string   `json:"escliente,omitempty"`
	Tiposolicitud            *string   `json:"tiposolicitud,omitempty"`
	Fechasolicitud           time.Time `json:"fechasolicitud,omitempty"`
	Poblacion                *string   `json:"poblacion,omitempty"`
	Provincia                *string   `json:"provincia,omitempty"`
	Direccion                *string   `json:"direccion,omitempty"`
	Cargo                    *string   `json:"cargo,omitempty"`
	Cantidaddeseada          *string   `json:"cantidaddeseada,omitempty"`
	Cantidadofrecida         *string   `json:"cantidadofrecida,omitempty"`
	Ncliente                 *string   `json:"ncliente,omitempty"`
	Calle                    *string   `json:"calle,omitempty"`
	CP                       *string   `json:"cp,omitempty"`
	Interesadoen             *string   `json:"interesadoen,omitempty"`
	Numero                   *string   `json:"numero,omitempty"`
	Compaiaactualfibraadsl   *string   `json:"compaiaactualfibraadsl,omitempty"`
	Companiaactualmovil      *string   `json:"companiaactualmovil,omitempty"`
	Fibraactual              *string   `json:"fibraactual,omitempty"`
	Moviactuallineaprincipal *string   `json:"moviactuallineaprincipal"`
	Numerolineasadicionales  int64     `json:"numerolineasadicionales,omitempty"`
	Tarifaactualsiniva       float64   `json:"tarifaactualsiniva,omitempty"`
	Motivocambio             *string   `json:"motivocambio,omitempty"`
	Importeaumentado         int64     `json:"importeaumentado,omitempty"`
	Importeretirado          int64     `json:"importeretirado,omitempty"`
	Tipoordenador            *string   `json:"tipoordenador,omitempty"`
	Sector                   *string   `json:"sector,omitempty"`
	Presupuesto              *string   `json:"presupuesto,omitempty"`
	Rendimiento              *string   `json:"rendimiento,omitempty"`
	Movilidad                *string   `json:"movilidad,omitempty"`
	Tipouso                  *string   `json:"tipouso,omitempty"`
	Office365                *string   `json:"Office365,omitempty"`
}

// Lead struct represents the fields of the leads table
type Lead struct {
	LeaID        int64      `sql:"primary_key" json:"-"`
	LeaTs        *time.Time `sql:"DEFAULT:current_timestamp" json:"-" `
	LeaDestiny   *string    `json:"lea_destiny,omitempty"`
	LeaExtracted *time.Time `sql:"default:current_timestamp" json:"-"`
	LeaStatus    *string    `json:"-"`
	LeaExtid     *string    `json:"-"`
	LeaCrmid     *string    `json:"-"`
	SouID        int64      `json:"sou_id,omitempty"`
	LeatypeID    int64      `json:"lea_type,omitempty"`
	LeaMedium    *string    `json:"-"`
	UtmSource    *string    `json:"utm_source,omitempty"`
	SubSource    *string    `json:"sub_source,omitempty"`
	LeaCampa     *string    `json:"-"`
	LeaPhone     *string    `json:"phone,omitempty"`
	LeaMail      *string    `json:"mail,omitempty"`
	LeaName      *string    `json:"name,omitempty"`
	LeaSurname   *string    `json:"surname,omitempty"`
	LeaURL       *string    `json:"url,omitempty"`
	Observations *string    `json:"observations,omitempty"`
	LeaIP        *string    `json:"ip,omitempty"`
	LeaAux1      *string    `json:"lea_aux1,omitempty"`
	LeaAux2      *string    `json:"lea_aux2,omitempty"`
	LeaAux3      *string    `json:"lea_aux3,omitempty"`
	LeaAux4      *string    `json:"lea_aux4,omitempty"`
	LeaAux5      *string    `json:"lea_aux5,omitempty"`
	LeaAux6      *string    `json:"lea_aux6,omitempty"`
	LeaAux7      *string    `json:"lea_aux7,omitempty"`
	LeaAux8      *string    `json:"lea_aux8,omitempty"`
	LeaAux9      *string    `json:"lea_aux9,omitempty"`
	LeaAux10     *string    `sql:"type:text  default null" json:"lea_aux10,omitempty"`

	SouIDLeontel     int64 `sql:"-" json:"sou_id_leontel"`
	LeatypeIDLeontel int64 `sql:"-" json:"lea_type_leontel"`
}

// LeadTest struct represents the fields of the leads table
type LeadTest struct {
	LeaID      int64 `sql:"primary_key"`
	LeaPhone   string
	LeaDestiny *string
	LeaTs      *time.Time `sql:"DEFAULT:current_timestamp"`
}

// LeontelResp represents a structure with the response returned
// by Leontel endpoint.
type LeontelResp struct {
	Success bool  `json:"success"`
	ID      int64 `json:"id"`
}

// SendLeadToLeontel sends the lead to Leontel endpoint
// Returns the response sended by the endpoint
func (l *Lead) SendLeadToLeontel() (*LeontelResp, error) {

	leadLeontel := l.LeadToLeontel()

	bytevalues, err := json.Marshal(leadLeontel)
	if err != nil {
		return nil, err
	}

	endpoint := "http://localhost:8888/leads/store"
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	leontelresp := LeontelResp{}
	json.Unmarshal(data, &leontelresp)

	return &leontelresp, nil
}

// LeadToLeontel maps values from Lead struct (webservice.leads table)
// as an appropiate lead Leontel input for each campaign
func (l *Lead) LeadToLeontel() LeadLeontel {
	log.Println("LeadToLeontel")
	log.Println(l)

	leontel := LeadLeontel{
		LeaSource: l.SouIDLeontel,
		LeaType:   l.LeatypeIDLeontel,
		Telefono:  l.LeaPhone,
		Nombre:    l.LeaName,
		URL:       l.LeaURL,
		IP:        l.LeaIP,
		Wsid:      l.LeaID,
	}
	log.Println("---------------------")

	// TODO avoid use of huge switch. How??
	switch souid := l.SouID; souid {
	case 1, 21, 22:
		// Creditea Abandonos
		// lea_aux1 => 1 | 0 ??
		// lea_surname => apellido1
		// lea_aux2 => dninie
		// lea_aux3 => asnef
		leontel.Apellido1 = l.LeaSurname
		leontel.Dninie = l.LeaAux2
		leontel.Asnef = l.LeaAux3
	case 9, 58:
		// Creditea EndToEnd + CREDITEA HM CORTO
		// lea_aux1 (dni)=> dninie
		// lea_aux2 (cantidadsolicitada)=> observaciones (DNI: $dninie Cantidad solicitada: $cantidadSolicitada)
		leontel.Dninie = l.LeaAux1
		// leontel.Observaciones = fmt.Sprintf("DNI:%s Cantidad solicitada:%s", *l.LeaAux1, *l.LeaAux2)
	case 10:
		// Creditea FB (no activo)
	case 11:
		// Creditea Rastreator (no activo)
		// lea_aux1 (cantidad solicitada)
		// lea_aux2 (Tipo contrato)
		// lea_aux3 (ingresos netos)
		// lea_aux4 (dni) => dninie
		// __ => observaciones (DNI: lea_aux4 Ingresos netos: lea_aux3
		// Tipo contrato lea_aux2 Cantidad solicitada lea_aux1)
		// lean_name => Nombre
		leontel.Nombre = l.LeaName
		// leontel.Observaciones = fmt.Sprintf(`
		// "DNI:%s Ingresos netos:%s Tipo Contrato:%s Cantidad solicitada:%s"`,
		// 	*l.LeaAux4, *l.LeaAux3, *l.LeaAux2, *l.LeaAux1)
	case 2:
		//Creditea Stand
		// lea_surname => apellido1
		// lea_aux2 => dninie
		// lea_aux3 => asnef
	case 17, 18, 19, 20:
		// Yoigo
		// lea_aux2 => observaciones
		// lea_aux3 => observaciones2
		leontel.Observaciones = l.LeaAux2
		leontel.Observaciones2 = l.LeaAux3
	case 52:
		// Incidencia Microsoft
		souidor, _ := strconv.Atoi(*l.LeaAux3)
		switch bb := souidor; bb {
		case 46:
		case 49:
			// Antes incidencia 26, 27, 28, 29, 30
			// Microsoft
			// lea_aux4 => tipoordenador
			// lea_aux5 => sector
			// lea_aux6 => presupuesto
			// lea_aux7 => rendimiento
			// lea_aux8 => movilidad
			// lea_aux9 + sou_id => tipouso
			// lea_aux10 => Office365
			// observations => observaciones2
			leontel.Tipoordenador = l.LeaAux4
			leontel.Sector = l.LeaAux5
			leontel.Presupuesto = l.LeaAux6
			leontel.Rendimiento = l.LeaAux7
			leontel.Movilidad = l.LeaAux8
			// leontel.Tipouso = fmt.Sprintf("%s %d", l.LeaAux9, l.SouID)
			leontel.Office365 = l.LeaAux10
			leontel.Observaciones2 = l.Observations
		case 50:
			// Antes incidencia 31, 32, 33, 34, 35
			// Microsoft
			// lea_aux10 => observaciones2
			leontel.Observaciones2 = l.LeaAux10
		case 51:
			// Antes incidencia 36, 37, 38, 39, 40
			// Microsoft
			// lea_aux10 => observaciones2
			leontel.Observaciones2 = l.LeaAux10
		case 48:
			// Microsoft
			// lea_aux10 => observaciones2
			leontel.Observaciones2 = l.LeaAux10
		default:
		}
	case 25:
		// Incidencia Microsoft

	case 53:
		// IPF
		// lea_aux1 => dninie
		// lea_aux2 => cantidadofrecida
		// lea_aux4 => ncliente
		// observations => observaciones
		leontel.Dninie = l.LeaAux1
		leontel.Cantidadofrecida = l.LeaAux2
		leontel.Ncliente = l.LeaAux4
		leontel.Observaciones = l.Observations
	default:
	}

	return leontel
}

func test(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
