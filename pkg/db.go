package apic2c

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // go mysql driver
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql import driver for gorm
)

// Database is a struct to manage DB environment configuration.
type Database struct {
	Host      string
	Port      int64
	User      string
	Pass      string
	Dbname    string
	Charset   string
	ParseTime string
	Loc       string

	db *gorm.DB
}

// Storer is an interface used to force the handler to implement
// the described methods
type Storer interface {
	Open() error
	Close()
	CreateTable(table interface{}) error
	Update(element interface{}, wCond string, wFields []string) error
	Insert(element interface{}) error
	Instance() *gorm.DB
}

// Open function opens a database connection using Database struct parameters
// Set the db property of the struct
// Return error | nil
func (env *Database) Open() error {
	connstring := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		env.User, env.Pass, env.Host, env.Port, env.Dbname, env.Charset, env.ParseTime, env.Loc)

	db, err := gorm.Open("mysql", connstring)
	if err != nil {
		log.Fatalf("Error opening database connection %s", err)
		return err
	}

	if err = db.DB().Ping(); err != nil {
		log.Fatalf("Error pinging database %s", err)
		return err
	}

	env.db = db

	return nil
}

// Close Database.db instance
func (env *Database) Close() {
	env.db.Close()
}

// CreateTable automatically migrate your schema, to keep your schema update to date.
// and create the table if not exists
func (env *Database) CreateTable(table interface{}) error {
	env.db.AutoMigrate(table)

	if !env.db.HasTable(table) {
		env.db.CreateTable(table)
	}
	return nil
}

// Insert generates a new row
func (env *Database) Insert(element interface{}) error {
	if result := env.db.Debug().Create(element); result.Error != nil {
		return fmt.Errorf("Error Insert element: %#v", result.Error)
	}
	return nil
}

// Update generates a new row
func (env *Database) Update(element interface{}, wCond string, wFields []string) error {
	// env.db.Model(&element).Where("active = ?", true).Update("name", "hello")
	wFieldsArr := []interface{}{}
	for _, z := range wFields {
		wFieldsArr = append(wFieldsArr, z)
	}
	env.db.Debug().Model(&element).Where(wCond, true).Update(wFieldsArr...)
	
	// db.Model(&user).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
	return nil
}

// Instance returns an instance of gorm DB
func (env *Database) Instance() *gorm.DB {
	return env.db
}

// TableName sets the default table name
func (Lead) TableName() string {
	return "leads"
}

// TableName sets the default table name
func (LeadTest) TableName() string {
	return "lead_tests"
}

// TableName sets the default table name
func (Source) TableName() string {
	return "sources"
}

// TableName sets the default table name
func (Leatype) TableName() string {
	return "leadtypes"
}
