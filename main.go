package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	leads "github.com/bysidecar/leads/pkg"
	model "github.com/bysidecar/leads/pkg/model"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	port := GetSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	database := &leads.Database{
		Host:      GetSetting("DB_HOST"),
		Port:      portInt,
		User:      GetSetting("DB_USER"),
		Pass:      GetSetting("DB_PASS"),
		Dbname:    GetSetting("DB_NAME"),
		Charset:   "utf8",
		ParseTime: "True",
		Loc:       "Local",
	}

	ch := leads.Handler{
		Storer: database,
	}

	if err := database.Open(); err != nil {
		log.Fatalf("error opening database connection. err: %s", err)
	}
	defer database.Close()

	if err := database.CreateTable(model.Lead{}); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	if err := database.CreateTable(model.RcableExp{}); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	if err := database.CreateTable(model.Microsoft{}); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	if err := database.CreateTable(model.Creditea{}); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	router := mux.NewRouter()
	router.PathPrefix("/store/leads/").Handler(ch.HandleFunction()).Methods(http.MethodPost)
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
