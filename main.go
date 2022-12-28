package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/lst123/nb-extractor/netbox"
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

func makeRoot(ch chan []byte) func(http.ResponseWriter, *http.Request) {
	var nb []byte
	var s sync.Mutex
	go func() {
		for {
			nbLA := <-ch
			s.Lock()
			nb = nbLA
			s.Unlock()
		}
	}()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		s.Lock()
		defer s.Unlock()
		if len(nb) > 0 {
			w.Write(nb)
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

	ch := make(chan []byte)
	// Start a http server
	handleRoot := makeRoot(ch)
	http.HandleFunc("/", handleRoot)
	go func() {
		log.Println("Starting server on port 8080.")
		http.ListenAndServe(":8080", nil)
	}()

	// Try to fetch new Json data and publish
	for {
		fmt.Println("Cycle is started...")
		c := make(chan netbox.RespData)
		go netbox.NetboxJson(prCfg.Netbox.Token, fullUrl.URL, c)
		fmt.Println("Got the Nebox data")
		r := <-c

		if r.Err != nil {
			ch <- []byte{}
			log.Println(r.Err)
		} else {
			fmt.Println("Sending data")
			ch <- r.Data
		}

		time.Sleep(10 * time.Minute)
	}
}
