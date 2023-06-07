package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lst123/nb-extractor/internal/netbox"
	"gopkg.in/yaml.v3"
)

type privCfg struct {
	Netbox struct {
		Site  string `yaml:"site"`
		Token string `yaml:"token"`
	} `yaml:"netbox"`
	Server struct {
		Site  string `yaml:"site"`
		Token string `yaml:"token"`
	} `yaml:"server"`
}

type Nb struct {
	data *[]byte
	mu   sync.Mutex
}

func parseYaml(file string, i interface{}) error {
	f, err := os.Open("./configs/" + file)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(i)
	if err != nil {
		return err
	}
	return nil
}

func makeRoot(nb *Nb, token string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer")
		if len(splitToken) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reqToken = strings.TrimSpace(splitToken[1])

		if reqToken != token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if nb.data != nil {
			w.Write(*nb.data)
		}
	}
}

func main() {
	// Get server params
	var prCfg privCfg
	err := parseYaml("private.yaml", &prCfg)
	if err != nil {
		log.Print(err)
		return
	}

	// Get roles
	var devRoles = make(map[string]int)
	err = parseYaml("deviceroles.yaml", devRoles)
	if err != nil {
		log.Print(err)
		return
	}

	// Prepare the full url to fetch from Netbox
	var fullUrl netbox.URL
	err = fullUrl.MakeUrl(prCfg.Netbox.Site, devRoles)
	if err != nil {
		log.Print(err)
		return
	}

	nb := Nb{}

	// Start a http server
	handleRoot := makeRoot(&nb, prCfg.Server.Token)
	http.HandleFunc("/", handleRoot)
	go func() {
		log.Println("Starting server on port 8080.")
		http.ListenAndServe(":8080", nil)
	}()

	// Try to fetch new Json data and publish
	for {
		nbLocal, err := netbox.NetboxJson(prCfg.Netbox.Token, fullUrl.URL)
		if err != nil {
			log.Println(err)
		} else {
			if len(nbLocal) > 0 {
				nb.mu.Lock()
				nb.data = &nbLocal
				nb.mu.Unlock()
			}
		}
		log.Println("Got the Netbox data")

		time.Sleep(5 * time.Minute)
	}
}
