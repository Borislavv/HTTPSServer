package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
		fmt.Fprintf(w, "Hello, HTTPS world!")
	})

	flag.Parse()
	domains := flag.Args()
	if len(domains) == 0 {
		log.Fatalf("Specify domains as arguments!")
	}


	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
	}
}