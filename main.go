package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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

type Devices struct {
	Devices []Device
}

func (d *Devices) AddDevice(dev Device) {
	d.Devices = append(d.Devices, dev)
}

type Device struct {
	Id         int
	Name       string
	Model      string
	Vendor     string
	DeviceRole string
	Serial     string
	Site       string
	Location   string
	Rack       string
	Status     string
	PrimaryIP4 string
	PrimaryIP6 string
	Rancid     bool
}

type NetboxDevices struct {
	Nb []NetboxDevice `json:"results"`
}

type NetboxDevice struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	DeviceType struct {
		Model        string `json:"model"`
		Manufacturer struct {
			Vendor string `json:"name"`
		} `json:"manufacturer"`
	} `json:"device_type"`
	DeviceRole struct {
		DeviceRole string `json:"name"`
	} `json:"device_role"`
	Serial string `json:"serial"`
	Site   struct {
		Site string `json:"name"`
	} `json:"site"`
	Location struct {
		Location string `json:"name"`
	} `json:"location"`
	Rack struct {
		Rack string `json:"display"`
	} `json:"rack"`
	Status struct {
		Status string `json:"value"`
	} `json:"status"`
	PrimaryIP4 struct {
		PrimaryIP4 string `json:"address"`
	} `json:"primary_ip4"`
	PrimaryIP6 struct {
		PrimaryIP6 string `json:"address"`
	} `json:"primary_ip6"`
	CustomFields struct {
		Rancid bool `json:"rancid"`
	} `json:"custom_fields"`
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

// func AskNetbox(partUrl string) ([]byte, error) {
// 	fullUrl := n.url + partUrl
// 	client := http.Client{Timeout: 10 * time.Second}
// 	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("Authorization", n.token)
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Accept", "application/json; indent=4")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Check if return code is not 200
// 	if res.StatusCode != http.StatusOK {
// 		return nil, errors.New("Http status is not 200")
// 	}

// 	defer res.Body.Close()

// 	resBody, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resBody, nil
// }

func prepareUrl(m map[string]int) (string, error) {
	if len(m) == 0 {
		return "", errors.New("can't construct URL with params")
	}
	var buffer bytes.Buffer
	buffer.WriteString("?role_id=")
	ids := make([]string, 0)
	for _, value := range m {
		strId := strconv.Itoa(value)
		ids = append(ids, strId)
	}
	buffer.WriteString(strings.Join(ids[:], "&role_id="))
	buffer.WriteString("&limit=0")
	urlString := buffer.String()

	return urlString, nil
}

func parseDevicesJson() {
	jsonFile, err := os.Open("devices.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var nbdevices NetboxDevices
	err = json.Unmarshal(byteValue, &nbdevices)
	if err != nil {
		fmt.Print(err)
		return
	}
	devices := Devices{}
	for i := 0; i < len(nbdevices.Nb); i++ {
		d := Device{}
		d.Id = nbdevices.Nb[i].Id
		d.Name = nbdevices.Nb[i].Name
		d.Model = nbdevices.Nb[i].DeviceType.Model
		d.Vendor = nbdevices.Nb[i].DeviceType.Manufacturer.Vendor
		d.DeviceRole = nbdevices.Nb[i].DeviceRole.DeviceRole
		d.Serial = nbdevices.Nb[i].Serial
		d.Site = nbdevices.Nb[i].Site.Site
		d.Location = nbdevices.Nb[i].Location.Location
		d.Rack = nbdevices.Nb[i].Rack.Rack
		d.Status = nbdevices.Nb[i].Status.Status
		d.PrimaryIP4 = nbdevices.Nb[i].PrimaryIP4.PrimaryIP4
		d.PrimaryIP6 = nbdevices.Nb[i].PrimaryIP6.PrimaryIP6
		d.Rancid = nbdevices.Nb[i].CustomFields.Rancid
		devices.AddDevice(d)
	}
	jData, err := json.Marshal(devices)
	if err != nil {
		fmt.Print(err)
		return
	}
	log.Println(string(jData))
	fmt.Printf("%+v", devices)
}

func main() {
	var prCfg privCfg
	var devRoles = make(map[string]int)

	err := parseYaml("private.yaml", &prCfg)
	if err != nil {
		log.Print(err)
		return
	}

	err = parseYaml("deviceroles.yaml", devRoles)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("%+v\n", prCfg)
	fmt.Printf("%+v\n", devRoles)
	paramsUrl, err := prepareUrl(devRoles)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println(paramsUrl)
	parseDevicesJson()
}
