package apic2c

// Leatype represents data structure of webservice sources table
type Leatype struct {
	LeatypeID          int64 `sql:"primary_key"`
	LeatypeDescription string
	LeatypeIdcrm       int64
}
