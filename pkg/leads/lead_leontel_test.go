package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActiveLeontel(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Index          int
		Description    string
		Lead           Lead
		Leontel        LeadLeontel
		ExpectedResult bool
	}{
		{
			Index:       1,
			Description: "when a lead not belongs to Leontel (Virgin)",
			Lead: Lead{
				SouID: 79,
			},
			Leontel:        LeadLeontel{},
			ExpectedResult: false,
		},
		{
			Index:       2,
			Description: "when a lead belongs to Leontel (ADESLAS END TO END)",
			Lead: Lead{
				SouID: 74,
			},
			Leontel:        LeadLeontel{},
			ExpectedResult: true,
		},
		{
			Index:       3,
			Description: "when a lead belongs to Leontel (R Cable End to End)",
			Lead: Lead{
				SouID: 64,
			},
			Leontel:        LeadLeontel{},
			ExpectedResult: true,
		},
	}

	for _, test := range tests {
		result := test.Leontel.Active(test.Lead)
		assert.Equal(result, test.ExpectedResult)
	}
}

func TestLeadToLeontel(t *testing.T) {
	assert := assert.New(t)

	t1 := "A"
	t2 := "B"
	t3 := "C"
	t4 := "D"
	t5 := "E"
	t6 := "F"
	t7 := "G"
	t9 := "J"
	observations := "H"

	obsTest := "A -- B"
	obsTest2 := "A -- C -- D -- B"
	obsTest3 := "A -- B -- C -- D -- E -- F"
	obsTest4 := "A -- B -- C -- D -- E -- F -- G -- A -- B -- C -- D -- E"
	obsTest5 := "A -- B -- C"

	kinkonAddress9 := fmt.Sprintf("%s -- %s -- %s", t4, t5, t6)
	obsKinkon9 := fmt.Sprintf("Cobertura -- %s -- Producto -- %s", t1, t9)

	fullnameKinkon10 := fmt.Sprintf("%s %s", t2, t3)
	obsKinkon10 := fmt.Sprintf(`%s -- Teléfono fijo portabilidad: -- %s -- Teléfono movil portabilidad: -- %s -- Teléfono movil 2 portabilidad: -- %s -- Operador movil portabilidad: -- %s -- Teléfono contacto -- %s -- Titular cuenta -- %s -- CCC -- %s`,
		obsKinkon9, t2, t4, t6, t7, t6, t1, t2)
	tEndesa := fmt.Sprintf(`%s -- Apellidos -- %s -- ¿Qué tipo de energía tienes en tu hogar? -- %s -- ¿Cuál es el tamaño de tu vivienda? -- %s -- ¿Cuántas personas viven en casa? -- %s -- ¿Qué tipo de energía usas en la calefacción? -- %s -- ¿Qué tipo de energía usas en la en la cocina? -- %s -- ¿Qué tipo de energía usas en el agua caliente? -- ¿Cada cuanto pones la lavadora? -- ¿Cada cuanto pones la secadora? -- ¿Cada cuanto pones el lavavajillas? -- ¿Eres el propietario de la vivienda? -- ¿Cuál es tu compañía actual?? -- Código postal -- Edad -- External ID -- `,
		t9, t6, t1, t2, t3, t4, t5)
	obsMvf := fmt.Sprintf(`Cobertura -- Producto -- Lead Reference Number -- %s -- Distribution ID -- %s -- ¿Ya tiene una centralita telefónica? -- %s -- ¿Cuantas extensiones necesita? -- %s -- Nº exacto de teléfonos -- %s -- ¿Cuántos empleados tiene su empresa? -- %s -- ¿Qué funcionalidad de centralita necesita? -- %s -- Apellidos -- %s -- Código Postal -- %s`, t2, t3, t4, t5, t6, t7, t7, t7, t7)

	tests := []struct {
		Index          int
		Description    string
		Lead           Lead
		ExpectedResult LeadLeontel
	}{
		{
			Index:       1,
			Description: "check data returned for sou_id 9 Creditea EndToEnd",
			Lead: Lead{
				SouID:        9,
				SouIDLeontel: 13,
				LeaDNI:       &t1,
				Creditea: &Creditea{
					RequestedAmount: &t2,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     13,
				Dninie:        &t1,
				Observaciones: &obsTest,
			},
		},
		{
			Index:       2,
			Description: "check data returned for sou_id 11 Creditea Rastreator",
			Lead: Lead{
				SouID:        11,
				SouIDLeontel: 15,
				LeaDNI:       &t1,
				Creditea: &Creditea{
					RequestedAmount: &t2,
					NetIncome:       &t3,
					ContractType:    &t4,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     15,
				Dninie:        &t1,
				Observaciones: &obsTest2,
			},
		},
		{
			Index:       3,
			Description: "check data returned for sou_id 46-49 Microsoft Hazelcambio + Recomendador",
			Lead: Lead{
				SouID:        46,
				SouIDLeontel: 61,
				Observations: &observations,
				Microsoft: &Microsoft{
					ComputerType: &t1,
					Sector:       &t2,
					Usecase:      &t3,
					Budget:       &t4,
					Performance:  &t5,
					Movility:     &t6,
					Office365:    &t7,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:      61,
				Tipoordenador:  &t1,
				Sector:         &t2,
				Tipouso:        &t3,
				Presupuesto:    &t4,
				Rendimiento:    &t5,
				Movilidad:      &t6,
				Office365:      &t7,
				Observaciones2: &observations,
			},
		},
		{
			Index:       4,
			Description: "check data returned for sou_id 48 Microsoft Calculadora",
			Lead: Lead{
				SouID:        48,
				SouIDLeontel: 61,
				Observations: &observations,
				Microsoft: &Microsoft{
					DevicesAverageAge:      &t1,
					DevicesOperatingSystem: &t2,
					DevicesHangFrequency:   &t3,
					DevicesNumber:          &t4,
					DevicesLastYearRepairs: &t5,
					DevicesStartupTime:     &t6,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:      61,
				Observaciones2: &obsTest3,
			},
		},
		{
			Index:       5,
			Description: "check data returned for sou_id 50 Microsoft Ofertas",
			Lead: Lead{
				SouID:         50,
				SouIDLeontel:  61,
				IsSmartCenter: false,
				Microsoft:     &Microsoft{},
			},
			ExpectedResult: LeadLeontel{
				LeaSource: 61,
			},
		},
		{
			Index:       6,
			Description: "check data returned for sou_id 51 Microsoft Ficha Producto",
			Lead: Lead{
				SouID:        51,
				SouIDLeontel: 61,
				Observations: &observations,
				Microsoft: &Microsoft{
					ProductType:        &t1,
					ProductName:        &t2,
					ProductID:          &t3,
					OriginalPrice:      &t4,
					Price:              &t5,
					Brand:              &t6,
					DiscountPercentage: &t7,
					DiscountCode:       &t1,
					ProcessorType:      &t2,
					DiskCapacity:       &t3,
					Graphics:           &t4,
					WirelessInterface:  &t5,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:      61,
				Observaciones2: &obsTest4,
			},
		},
		{
			Index:       7,
			Description: "check data returned for sou_id 54 R Cable Expansion",

			Lead: Lead{
				SouID:        54,
				SouIDLeontel: 63,
				RcableExp: &RcableExp{
					RespValues: &t1,
					Location:   &t2,
					Answer:     &t3,
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     63,
				Observaciones: &obsTest5,
			},
		},
		{
			Index:       8,
			Description: "check data returned for sou_id 64 R Cable End To End (Kinkon) C2C",
			Lead: Lead{
				SouID:         64,
				SouIDLeontel:  73,
				LeaPhone:      &t1,
				IsSmartCenter: false,
				Kinkon: &Kinkon{
					Coverture:   &t1,
					Product:     &t9,
					CovData:     CovData{},
					Portability: Portability{},
					HolderData:  HolderData{},
					BillingInfo: BillingInfo{},
				},
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     73,
				Telefono:      &t1,
				Observaciones: &obsKinkon9,
			},
		},
		{
			Index:       9,
			Description: "check data returned for sou_id 64 R Cable End To End (Kinkon) Coverture check",
			Lead: Lead{
				SouID:            64,
				SouIDLeontel:     73,
				LeatypeIDLeontel: 26,
				Kinkon: &Kinkon{
					Coverture: &t1,
					Product:   &t9,
					CovData: CovData{
						State:    &t2,
						Town:     &t3,
						Street:   &t4,
						Number:   &t5,
						Floor:    &t6,
						CovPhone: &t7,
					},
					Portability: Portability{},
					HolderData:  HolderData{},
					BillingInfo: BillingInfo{},
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     73,
				LeaType:       26,
				Provincia:     &t2,
				Poblacion:     &t3,
				Direccion:     &kinkonAddress9,
				Observaciones: &obsKinkon9,
			},
		},
		{
			Index:       10,
			Description: "check data returned for sou_id 66 Euskaltel End To End (Kinkon) Hiring process",
			Lead: Lead{
				SouID:            66,
				SouIDLeontel:     75,
				LeatypeIDLeontel: 30,
				Kinkon: &Kinkon{
					Coverture: &t1,
					Product:   &t9,
					CovData: CovData{
						State:    &t2,
						Town:     &t3,
						Street:   &t4,
						Number:   &t5,
						Floor:    &t6,
						CovPhone: &t7,
					},
					Portability: Portability{
						Phone:                &t2,
						PhoneProvider:        &t3,
						MobilePhone:          &t4,
						MobilePhoneProvider:  &t5,
						MobilePhone2:         &t6,
						MobilePhoneProvider2: &t7,
					},
					HolderData: HolderData{
						Name:         &t2,
						Surname:      &t3,
						Idnumber:     &t4,
						Mail:         &t5,
						ContactPhone: &t6,
					},
					BillingInfo: BillingInfo{
						AccountHolder: &t1,
						AccountNumber: &t2,
					},
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:              75,
				LeaType:                30,
				Provincia:              &t2,
				Poblacion:              &t3,
				Direccion:              &kinkonAddress9,
				Compaiaactualfibraadsl: &t3,
				Companiaactualmovil:    &t5,
				Dninie:                 &t4,
				Email:                  &t5,
				Nombrecompleto:         &fullnameKinkon10,
				Observaciones:          &obsKinkon10,
			},
		},
		{
			Index:       11,
			Description: "check data returned for sou_id 69 Alterna campaign",
			Lead: Lead{
				SouID:         69,
				SouIDLeontel:  78,
				LeaPhone:      &t7,
				IsSmartCenter: false,
				Alterna: &Alterna{
					PostalCode:  &t1,
					Street:      &t2,
					Number:      &t3,
					InstallType: &t4,
					CPUS:        &t5,
				},
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     78,
				Telefono:      &t7,
				CP:            &t1,
				Calle:         &t2,
				Numero:        &t3,
				Tiposolicitud: &t4,
				Observaciones: &t5,
			},
		},
		{
			Index:       12,
			Description: "check data returned for sou_id 78 Endesa campaign",
			Lead: Lead{
				SouID:         78,
				SouIDLeontel:  87,
				LeaPhone:      &t7,
				LeaName:       &t1,
				IsSmartCenter: false,
				Observations:  &t9,
				Endesa: &Endesa{
					Surname:        &t6,
					TypeEnergy:     &t1,
					HomeSize:       &t2,
					HomePopulation: &t3,
					TypeHeating:    &t4,
					TypeKitchen:    &t5,
				},
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     87,
				Telefono:      &t7,
				Nombre:        &t1,
				Observaciones: &tEndesa,
			},
		},
		{
			Index:       13,
			Description: "check data returned for sou_id 77 Adeslas campaign",
			Lead: Lead{
				SouID:         77,
				SouIDLeontel:  86,
				LeaPhone:      &t7,
				IsSmartCenter: false,
				Observations:  &t2,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     86,
				Telefono:      &t7,
				Observaciones: &t2,
			},
		},
		{
			Index:       14,
			Description: "check data returned for sou_id 64 R Cable End To End (Kinkon) MVF Provider",
			Lead: Lead{
				SouID:            74,
				SouIDLeontel:     83,
				LeatypeIDLeontel: 33,
				Kinkon: &Kinkon{
					Mvf: Mvf{
						LeadReferenceNumber:      &t2,
						DistributionID:           &t3,
						HasSwitchboard:           &t4,
						ExtensionsNumber:         &t5,
						PhoneAmount:              &t6,
						EmployeeNumber:           &t7,
						SwitchboardFunctionality: &t7,
						Surname:                  &t7,
						PostalCode:               &t7,
					},
				},
				IsSmartCenter: false,
			},
			ExpectedResult: LeadLeontel{
				LeaSource:     83,
				LeaType:       33,
				Observaciones: &obsMvf,
			},
		},
	}

	for _, test := range tests {
		leontel := test.Lead.LeadToLeontel()
		assert.Equal(test.ExpectedResult, leontel)
	}
}

func TestSenLeadToLeontel(t *testing.T) {
	assert := assert.New(t)

	phone := HelperRandstring(9)

	tests := []struct {
		Description      string
		TypeRequest      string
		StatusCode       int
		LeadLeontel      []LeadLeontel
		ExpectedResult   bool
		ExpectedResponse LeontelResp
	}{
		{
			Description: "when HandleFunction receive a POST request with no data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusInternalServerError,
			LeadLeontel: []LeadLeontel{},
			ExpectedResponse: LeontelResp{
				Success: false,
				Error:   "Origen del lead erroneo use getSources()",
			},
		},
		{
			Description: "when HandleFunction receive a POST request with data",
			TypeRequest: http.MethodPost,
			StatusCode:  http.StatusOK,
			LeadLeontel: []LeadLeontel{
				{
					LeaSource: 23,
					LeaType:   2,
					Telefono:  &phone,
				},
			},
			ExpectedResponse: LeontelResp{
				Success: true,
				ID:      9999,
			},
		},
	}

	for _, test := range tests {

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			testresp, err := json.Marshal(test.ExpectedResponse)
			assert.NoError(err)

			res.WriteHeader(test.StatusCode)
			res.Write(testresp)
		}))
		defer func() { ts.Close() }()

		body, err := json.Marshal(test.LeadLeontel)
		assert.NoError(err)

		req, err := http.NewRequest(test.TypeRequest, ts.URL, bytes.NewBuffer(body))
		assert.NoError(err)

		http := &http.Client{}
		resp, err := http.Do(req)
		assert.NoError(err)

		assert.Equal(resp.StatusCode, test.StatusCode)

		response := LeontelResp{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(err)

		log.Println(response)

		assert.Equal(response.Success, test.ExpectedResponse.Success)
		assert.Equal(response.ID, test.ExpectedResponse.ID)
	}
}
