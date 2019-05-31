package apic2c

// Source represents data structure of webservice sources table
type Source struct {
	SouID          int64 `sql:"primary_key"`
	SouDescription string
	SouIdcrm       int64
}
