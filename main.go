package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	Version = "0.1"
)

var (
	port  int
	token string
)

var flags = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

func init() {
	flags.IntVar(&port, "p", 8888, "port to run on.")
	flags.StringVar(&token, "t", "changeme", "token that the github webhook should pass")

	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "buildbot version %s\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] path\n\n", os.Args[0])
		flags.PrintDefaults()
	}
}

func handleDockerHubWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		http.Error(w, "405 Method not allow", http.StatusMethodNotAllowed)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

        fmt.Println(body)
}

func main() {
	flags.Parse(os.Args[1:])
	log.Println("Running on port", port)
	log.Println("Using token", token)

	http.HandleFunc("/", handleDockerHubWebhook)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
