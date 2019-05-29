package apic2c

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
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

type firstobject struct {
	ID int64
}

// Storer is an interface used to force the handler to implement
// the described methods
type Storer interface {
	Open() error
	Close()
	CreateTable() error
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
func (env *Database) CreateTable() error {
	env.db.AutoMigrate(&firstobject{})

	if !env.db.HasTable(&firstobject{}) {
		env.db.CreateTable(&firstobject{})
	}
	return nil
}
