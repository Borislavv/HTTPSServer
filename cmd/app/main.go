package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		/* HTTPS Header */
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

	/* Create tmp. cache dir */
	dir := cacheDir()
	if dir != "" {
		certManager.Cache = autocert.DirCache(dir)
	}

	/* Create instance of server */
	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	log.Printf("Serving http/https for domains: %+v", domains)
	go func() {
		/* Serve HTTP, which will redirect automatically to HTTPS */
		h := certManager.HTTPHandler(nil)
		log.Fatal(http.ListenAndServe(":http", h))
	}()

	/* Serve HTTPS */
	log.Fatal(server.ListenAndServeTLS("", ""))
}

/* Try to create a tmp. cache dir. */
func cacheDir() (dir string) {
	if u, _ := user.Current(); u != nil {
		dir = filepath.Join(os.TempDir(), "cache-golang-autocert-"+u.Username)
		if err := os.MkdirAll(dir, 0700); err == nil {
			return dir
		}
	}
	return ""
}