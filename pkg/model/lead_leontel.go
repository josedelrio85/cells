package leads

import "time"

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
	Wsid                     int64      `json:"wsid,omitempty"`
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
	Success bool  `json:"success"`
	ID      int64 `json:"id"`
}
