package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"path"
)

var script string
var httpHost string
var env []string

func handler(w http.ResponseWriter, r *http.Request) {
	sp := script
	if sp == "" {
		sp = fmt.Sprintf(".%s", r.URL.Path)
		if s, err := os.Stat(sp); os.IsNotExist(err) {
			w.WriteHeader(404)
			w.Write([]byte("404 not found"))
			return
		} else if s.IsDir() {
			fs := http.FileServer(http.Dir("."))
			fs.ServeHTTP(w, r)
			return
		}
	}
	if httpHost != "" {
		r.Header.Set("Host", httpHost)
		r.Host = httpHost
		r.URL.Host = httpHost
	}
	cgi := cgi.Handler{
		Path: "./" + path.Base(sp),
		Dir:  path.Dir(sp),
	}
	cgi.ServeHTTP(w, r)
}

func main() {
	var err error
	var dir string
	var addr string
	flag.StringVar(
		&script,
		"f", "",
		"path to the CGI script (relative to -d), if not set, request path is used")
	flag.StringVar(
		&dir,
		"d", "",
		"root directory, if not set, current work dir is used")
	flag.StringVar(&addr, "l", ":8080", "listen address")
	flag.StringVar(
		&httpHost,
		"http-host", "",
		"overrides host header")
	flag.Parse()
	env = flag.Args()
	if dir != "" {
		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	if script != "" {
		if _, err := os.Stat(script); err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
