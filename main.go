package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	leads "github.com/bysidecar/leads/pkg"
	hooks "github.com/bysidecar/leads/pkg/hooks"
	model "github.com/bysidecar/leads/pkg/model"
	"github.com/gomodule/redigo/redis"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	port := GetSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	database := &model.Database{
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
		ActiveHooks: []hooks.Hookable{
			hooks.Asnef{},
			hooks.Ontime{},
			hooks.Hibernated{},
			hooks.Allowed{},
		},
		Redis: model.Redis{
			Pool: &redis.Pool{
				MaxIdle:     5,
				IdleTimeout: 60 * time.Second,
				Dial: func() (redis.Conn, error) {
					return redis.Dial("tcp", GetSetting("CHECK_LEAD_REDIS")+":6379")
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
