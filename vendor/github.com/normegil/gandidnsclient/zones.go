package gandidnsclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Zone struct {
	*zoneJSON
}

type zoneJSON struct {
	ID              string  `json:"uuid,omitempty"`
	Name            string  `json:"name,omitempty"`
	Email           string  `json:"email,omitempty"`
	Serial          int     `json:"serial,omitempty"`
	Minimum         int     `json:"minimum,omitempty"`
	Refresh         int     `json:"refresh,omitempty"`
	Retry           int     `json:"retry,omitempty"`
	Expire          int     `json:"expire,omitempty"`
	SharingID       string  `json:"sharing_id,omitempty"`
	PrimaryNS       url.URL `json:"primary_ns,omitempty"`
	ZoneLink        url.URL `json:"zone_href,omitempty"`
	ZoneRecordsLink url.URL `json:"zone_records_href,omitempty"`
}

func (z Zone) Refresh() time.Duration {
	return time.Duration(z.zoneJSON.Refresh) * time.Second
}

func (z Zone) Retry() time.Duration {
	return time.Duration(z.zoneJSON.Retry) * time.Second
}

func (z Zone) Expire() time.Duration {
	return time.Duration(z.zoneJSON.Expire) * time.Second
}

func (z Zone) MarshalJSON() ([]byte, error) {
	return json.Marshal(z.zoneJSON)
}

func (z *Zone) UnmarshalJSON(b []byte) error {
	var jsonObj zoneJSON
	if err := json.Unmarshal(b, &jsonObj); nil != err {
		return err
	}
	z.zoneJSON = &jsonObj
	return nil
}

func (c Client) ListZones() ([]Zone, error) {
	req, err := c.request("GET", "/zone", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var zones []Zone
	if err := json.Unmarshal(body, &zones); err != nil {
		return nil, err
	}
	return zones, nil
}
