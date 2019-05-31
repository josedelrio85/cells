package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	apic2c "github.com/bysidecar/api_ws/pkg"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {

	port := GetSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	database := &apic2c.Database{
		Host:      GetSetting("DB_HOST"),
		Port:      portInt,
		User:      GetSetting("DB_USER"),
		Pass:      GetSetting("DB_PASS"),
		Dbname:    GetSetting("DB_NAME"),
		Charset:   "utf8",
		ParseTime: "True",
		Loc:       "Local",
	}
	ch := apic2c.Handler{
		Storer: database,
	}

	if err := database.Open(); err != nil {
		log.Fatalf("error opening database connection. err: %s", err)
	}
	defer database.Close()

	lead := apic2c.Lead{}
	leadtest := apic2c.LeadTest{}
	source := apic2c.Source{}

	if err := database.CreateTable(lead); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	if err := database.CreateTable(leadtest); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	if err := database.CreateTable(source); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	base := router.PathPrefix("/apic2c").Subrouter()

	rcable := base.PathPrefix("/rcable").Subrouter()
	rcable.HandleFunc("/incomingC2C", ch.RcableHandler)

	creditea := base.PathPrefix("/creditea").Subrouter()
	creditea.HandleFunc("/test", ch.TestHandler)

	router.Use(apic2c.Middleware)
	log.Fatal(http.ListenAndServe(":5000", cors.Default().Handler(router)))
}

// GetSetting reads an ENV VAR setting, it does crash the service if with an
// error message if any setting is not found.
//
// - setting: The setting (ENV VAR) to read.
//
// Returns the setting value.
func GetSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}
