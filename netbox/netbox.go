package netbox

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
