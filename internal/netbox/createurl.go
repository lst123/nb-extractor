package netbox

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type URL struct {
	URL string
}

func (u *URL) MakeUrl(netboxUrl string, m map[string]int) error {
	if len(m) == 0 {
		return errors.New("can't construct URL with params")
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
	u.URL = netboxUrl + urlString
	return nil
}
