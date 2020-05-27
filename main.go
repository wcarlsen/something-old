package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

type Deprecations []Deprecation

// Tag data from the call
type Deprecation struct {
	Name       string
	Namespace  string
	Kind       string
	ApiVersion string
	RuleSet    string
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func FetchApiDeprecations() {
	for true {
		cmd := exec.Command("./script.sh")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Fetched api deprecations")

		time.Sleep(5 * time.Second)
	}
}

func K8sApiDeprecationsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./data.json")
	if err != nil {
		log.Fatal(err)
	}
	// json data
	var obj Deprecations

	// unmarshall it
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(data)))
}

func main() {
	// Start go routine
	go FetchApiDeprecations()

	// Initialise new router
	r := mux.NewRouter()

	// Initialise routes
	r.HandleFunc("/healthz", HealthCheckHandler)

	r.HandleFunc("/k8s/api/deprecations", K8sApiDeprecationsHandler)
	//r.HandleFunc("/k8s/api/deprecations/{namespace}", K8sApiDeprecationsNamespaceHandler)

	// Bind port and pass in router
	log.Fatal(http.ListenAndServe(":3000", r))
}
