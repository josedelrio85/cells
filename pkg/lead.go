package apic2c

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"net"
	"fmt"
	"io"

	"github.com/jinzhu/gorm"
)

// LeadLeontel represents the data structure of lea_leads Leontel table
type LeadLeontel struct {
	// lea_ fields
	LeaID            int64     `json:"-"`
	LeaType          int64     `json:"lea_type,omitempty"`
	LeaTs            *time.Time `json:"-"`
	LeaSource        int64     `json:"lea_source,omitempty"`
	LeaLot           int64     `json:"-"`
	LeaAssigned      int64     `json:"-"`
	LeaScheduled     *time.Time `json:"lea_scheduled,omitempty"`
	LeaScheduledAuto int64     `json:"lea_scheduled_auto,omitempty"`
	LeaCost          float64   `json:"-"`
	LeaSeen          int64     `json:"-"`
	LeaClosed        int64     `json:"lea_closed,omitempty"`
	LeaNew           int64     `json:"-"`
	LeaClient        int64     `json:"-"`
	LeaDateControl   *time.Time `json:"-"`

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
	Fechasolicitud           *time.Time `json:"fechasolicitud,omitempty"`
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
	Moviactuallineaprincipal *string   `json:"moviactuallineaprincipal,omitempty"`
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
	LeatypeDescLeontel string `sql:"-" json:"lea_type_desc_leontel"`

	Gclid *string `sql:"-" json:"glcid"`
	Domain *string `sql:"-" json:"domain"`

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
		Wsid:      lead.LeaID,
	}

	// TODO avoid use of huge switch. How??
	switch souid := lead.SouID; souid {
	case 1, 21, 22:
		// Creditea Abandonos
		// lea_aux1 => 1 | 0 ??
		// lea_surname => apellido1
		// lea_aux2 => dninie
		// lea_aux3 => asnef
		leontel.Apellido1 = lead.LeaSurname
		leontel.Dninie = lead.LeaAux2
		leontel.Asnef = lead.LeaAux3
	case 9, 58:
		// Creditea EndToEnd + CREDITEA HM CORTO
		// lea_aux1 (dni)=> dninie
		// lea_aux2 (cantidadsolicitada)=> observaciones (DNI: $dninie Cantidad solicitada: $cantidadSolicitada)
		leontel.Dninie = lead.LeaAux1
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
		leontel.Nombre = lead.LeaName
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
		leontel.Observaciones = lead.LeaAux2
		leontel.Observaciones2 = lead.LeaAux3
	case 52:
		// Incidencia Microsoft
		souidor, _ := strconv.Atoi(*lead.LeaAux3)
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
			leontel.Tipoordenador = lead.LeaAux4
			leontel.Sector = lead.LeaAux5
			leontel.Presupuesto = lead.LeaAux6
			leontel.Rendimiento = lead.LeaAux7
			leontel.Movilidad = lead.LeaAux8
			// leontel.Tipouso = fmt.Sprintf("%s %d", lead.LeaAux9, lead.SouID)
			leontel.Office365 = lead.LeaAux10
			leontel.Observaciones2 = lead.Observations
		case 50:
			// Antes incidencia 31, 32, 33, 34, 35
			// Microsoft
			// lea_aux10 => observaciones2
			leontel.Observaciones2 = lead.LeaAux10
		case 51:
			// Antes incidencia 36, 37, 38, 39, 40
			// Microsoft
			// lea_aux10 => observaciones2
			leontel.Observaciones2 = lead.LeaAux10
		case 48:
			// Microsoft
			// lea_aux10 => observaciones2
			leontel.Observaciones2 = lead.LeaAux10
		default:
		}
	case 25:
		// Incidencia Microsoft
		leontel.Observaciones = lead.Observations
		leontel.Observaciones2 = &lead.LeatypeDescLeontel
	case 53:
		// IPF
		// lea_aux1 => dninie
		// lea_aux2 => cantidadofrecida
		// lea_aux4 => ncliente
		// observations => observaciones
		leontel.Dninie = lead.LeaAux1
		leontel.Cantidadofrecida = lead.LeaAux2
		leontel.Ncliente = lead.LeaAux4
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

	endpoint := "http://localhost:8888/leads/store"
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytevalues))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	leontelresp := []LeontelResp{}
	json.Unmarshal(data, &leontelresp)

	return &leontelresp[0], nil
}

// UpdatePostLeontel updates lead row after the insertion in Leontel is made and its result is succesful
func (lead *Lead) UpdatePostLeontel(db *gorm.DB, leontelID int64) error {
	status := "SENT"
	now := time.Now()
	crmid := strconv.FormatInt(leontelID, 10)
	if result := db.Model(lead).Where("lea_id = ?", lead.LeaID).Update(Lead{LeaExtracted: &now, LeaStatus: &status, LeaCrmid: &crmid}); result.Error != nil {
		return fmt.Errorf("Error updating Lead row after inserting Leontel: %#v", result.Error)
	}
	return nil
}

// GetLeontelValues queries for the Leontel equivalences
// of sou_id and lea_type values
func (lead *Lead) GetLeontelValues(db *gorm.DB) error {
	source := Source{}
	leatype := 	Leatype{}
	
	if result := db.Where("sou_id = ?", lead.SouID).First(&source); result.Error != nil {
		return fmt.Errorf("Error retrieving SouIDLeontel value: %#v", result.Error)
	}
	if result := db.Where("leatype_id = ?", lead.LeatypeID).First(&leatype); result.Error != nil {
		return fmt.Errorf("error retrieving LeatypeIDLeontel value: %#v", result.Error)
	}
	lead.SouIDLeontel = source.SouIdcrm
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

