package leads

import (
	"fmt"
	"log"

	model "github.com/bysidecar/leads/pkg/model"
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

	*gorm.DB
}

// Storer is an interface used to force the handler to implement
// the described methods
type Storer interface {
	Open() error
	Close()
	Update(element interface{}, wCond string, wFields []string) error
	Insert(element interface{}) error
	Instance() *gorm.DB
}

// Open function opens a database connection using Database struct parameters
// Set the DB property of the struct
// Return error | nil
func (d *Database) Open() error {
	connstring := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		d.User, d.Pass, d.Host, d.Port, d.Dbname, d.Charset, d.ParseTime, d.Loc)

	DB, err := gorm.Open("mysql", connstring)
	if err != nil {
		log.Fatalf("Error opening database connection %s", err)
		return err
	}

	if err = DB.DB().Ping(); err != nil {
		log.Fatalf("Error pinging database %s", err)
		return err
	}

	d.DB = DB

	return nil
}

// Close Database.DB instance
func (d *Database) Close() {
	d.DB.Close()
}

// CreateTable automatically migrate your schema, to keep your schema update to date.
// and create the table if not exists
func (d *Database) AutoMigrate() error {
	if err := d.DB.AutoMigrate(model.Lead{}).Error; err != nil {
		return err
	}

	if err := d.DB.AutoMigrate(model.RcableExp{}).Error; err != nil {
		return err
	}

	if err := d.DB.AutoMigrate(model.Microsoft{}).Error; err != nil {
		return err
	}

	return nil
}

// Insert generates a new row
func (d *Database) Insert(element interface{}) error {
	if result := d.Create(element); result.Error != nil {
		return fmt.Errorf("Error Insert element: %#v", result.Error)
	}
	return nil
}

// Update generates a new row
func (d *Database) Update(element interface{}, wCond string, wFields []string) error {
	// d.Model(&element).Where("active = ?", true).Update("name", "hello")
	wFieldsArr := []interface{}{}
	for _, z := range wFields {
		wFieldsArr = append(wFieldsArr, z)
	}
	d.Model(&element).Where(wCond, true).Update(wFieldsArr...)

	//.Model(&user).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
	return nil
}

func (d *Database) Instance() *gorm.DB {
	return d.DB
}
