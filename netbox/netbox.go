package netbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RespData struct {
	Data []byte
	Err  error
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

func fetchJson(token, url string) ([]byte, error) {
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json; indent=4")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Check if return code is not 200
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("http status is not 200")
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

func prepareJson(b []byte) ([]byte, error) {
	// this func make appropriate (plain) structure
	var nbdevices NetboxDevices
	err := json.Unmarshal(b, &nbdevices)
	if err != nil {
		return nil, err
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

	dataJson, err := json.Marshal(devices)
	if err != nil {
		return nil, err
	}

	return dataJson, nil
}

func NetboxJson(token string, url string, c chan<- RespData) {
	defer close(c)
	nbJson, err := fetchJson(token, url)
	if err != nil {
		c <- RespData{Err: errors.New("can't fetch Json from Netbox")}
	}
	dataJson, err := prepareJson(nbJson)
	if err != nil {
		c <- RespData{Err: errors.New("can't create a new Json and marshal it")}
	}
	r := RespData{Data: dataJson, Err: nil}
	c <- r
}
