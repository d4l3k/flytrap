package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	bind = flag.String("bind", ":8080", "address to bind to")
	out  = flag.String("outdir", "out", "dir to drop output files into")
)

var slugRE = regexp.MustCompile("[^a-zA-Z0-9_\\-.]+")

func slugify(s string) string {
	return slugRE.ReplaceAllString(s, "-")
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	flag.Parse()
	if err := os.MkdirAll(*out, 0700); err != nil {
		return err
	}
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error! %+v", err)
			return
		}
		log.Printf("%s", body)

		now := time.Now().String()
		name := now + "-" + slugify(r.URL.Path)
		path := filepath.Join(*out, name)
		if err := ioutil.WriteFile(path, body, 0700); err != nil {
			log.Printf("error! %+v", err)
			return
		}
	})
	log.Printf("Listening %s...", *bind)
	if err := http.ListenAndServe(*bind, nil); err != nil {
		return err
	}
	return nil
}
