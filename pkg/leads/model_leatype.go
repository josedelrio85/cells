package leads

// Leatype represents data structure of webservice sources table
type Leatype struct {
	LeatypeID          int64 `sql:"primary_key"`
	LeatypeDescription string
	LeatypeIdcrm       int64
}

// TableName sets the default table name
func (Leatype) TableName() string {
	return "leadtypes"
}
