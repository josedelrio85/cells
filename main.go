package main

import (
	"log"
	"net/http"

	apic2c "github.com/bysidecar/api_ws/pkg"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	ch := apic2c.Handler{}
	r.PathPrefix("/apic2c/").Handler(ch.HandleFunction())

	log.Fatal(http.ListenAndServe(":5000", cors.Default().Handler(r)))
}
