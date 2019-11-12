package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	lead "github.com/bysidecar/leads/pkg/leads"
	redisclient "github.com/bysidecar/leads/pkg/leads/redis"
	"github.com/gomodule/redigo/redis"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	port := getSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

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

	ch := lead.Handler{
		Storer: database,
		ActiveHooks: []lead.Hookable{
			lead.Hibernated{},
			lead.Phone{},
			lead.DuplicatedTime{},
			lead.DuplicatedSmartCenter{},
			lead.Ontime{},
			lead.Gclid{},
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
	}

	if err := database.Open(); err != nil {
		log.Fatalf("error opening database connection. err: %s", err)
	}
	defer database.Close()

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
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
