package leads

// Source represents data structure of webservice sources table
type Source struct {
	SouID          int64
	SouDescription string
	SouIdcrm       int64
	SouIDEvolution string
}

// TableName sets the default table name
func (Source) TableName() string {
	return "sources"
}
