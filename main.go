package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	lead "github.com/josedelrio85/leads/pkg/leads"
	redisclient "github.com/josedelrio85/leads/pkg/leads/redis"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Leads API...")
	port := getSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	dev := false
	devstring := getSetting("DEV")
	if devstring == "true" {
		dev = true
	}
	log.Printf("Are we working on dev? %t", dev)

	database := &lead.Database{
		Host:      getSetting("DB_HOST"),
		Port:      portInt,
		User:      getSetting("DB_USER"),
		Pass:      getSetting("DB_PASS"),
		Dbname:    getSetting("DB_NAME"),
		Charset:   "utf8",
		ParseTime: "True",
		Loc:       "Local",
	}

	reportdb := &lead.Database{
		Host:      getSetting("DB_HOST_REPORT"),
		Port:      portInt,
		User:      getSetting("DB_USER_REPORT"),
		Pass:      getSetting("DB_PASS_REPORT"),
		Dbname:    getSetting("DB_NAME"),
		Charset:   "utf8",
		ParseTime: "True",
		Loc:       "Local",
	}

	ch := lead.Handler{
		Storer:   database,
		Reporter: reportdb,
		ActiveHooks: []lead.Hookable{
			lead.Hibernated{},
			lead.Phone{},
			lead.DuplicatedTime{},
			lead.DuplicatedSmartCenter{},
			lead.Ontime{},
			lead.RejectSC{},
			lead.MapType{},
			// lead.Gclid{},
		},
		ActiveSc: []lead.Scable{
			lead.LeadLeontel{},
			lead.Evolution{},
		},
		Redis: redisclient.Redis{
			Pool: &redis.Pool{
				MaxIdle:     5,
				IdleTimeout: 60 * time.Second,
				Dial: func() (redis.Conn, error) {
					return redis.Dial("tcp", getSetting("CHECK_LEAD_REDIS")+":6379")
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			},
		},
		Dev: dev,
	}

	// database
	if err := database.Open(); err != nil {
		log.Fatalf("error opening database connection. err: %s", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	// report
	if !ch.Dev {
		log.Println("Opening report DB")
		if err := reportdb.Open(); err != nil {
			log.Fatalf("error opening report database connection. err: %s", err)
		}
		defer reportdb.Close()

		if err := reportdb.AutoMigrate(); err != nil {
			log.Fatalf("error creating the table. err: %s", err)
		}
	}

	router := mux.NewRouter()

	router.PathPrefix("/lead/store/").Handler(ch.HandleFunction()).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":4000", cors.Default().Handler(router)))
}

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}
