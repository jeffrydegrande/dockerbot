package main

import (
        "flag"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "os/exec"
        "encoding/json"
)

const (
        Version = "0.1"
)

var (
        port  int
        token string
)

type Event struct {
   PushData PushData `json:"push_data"`
   Repository Repository `json:"repository"`
}

type PushData struct {
  Pusher string `json:"pusher"`
}

type Repository struct {
   RepoName string `json:"repo_name"`
}


var flags = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

func init() {
        flags.IntVar(&port, "p", 8888, "port to run on.")
        flags.StringVar(&token, "t", "changeme", "token provided through webhook url 'token' parameter")
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

        if (r.URL.Query().Get("token") != token) {
                http.Error(w, "403 Forbidden", http.StatusForbidden)
        }

        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
                log.Println(err.Error())
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        var event Event
        json.Unmarshal(body, &event)

        fmt.Println("Pulling", event.Repository.RepoName)
        cmd := exec.Command("sh", "build.sh", event.Repository.RepoName)
        stdout, err := cmd.Output()
        if err != nil {
                fmt.Println(err.Error())
        }
        fmt.Println(string(stdout))
}

func main() {
        flags.Parse(os.Args[1:])
        log.Println("Running on port", port)

        http.HandleFunc("/", handleDockerHubWebhook)
        log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
