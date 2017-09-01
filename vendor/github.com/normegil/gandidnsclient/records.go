package gandidnsclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

type RecordType string

const (
	TypeA     = RecordType("A")
	TypeCName = RecordType("CName")
	TypeAAAA  = RecordType("AAAA")
	TypeCAA   = RecordType("CAA")
	TypeCDS   = RecordType("CDS")
	TypeDNAME = RecordType("DNAME")
	TypeDS    = RecordType("DS")
	TypeLOC   = RecordType("LOC")
	TypeMX    = RecordType("MX")
	TypeNS    = RecordType("NS")
	TypePTR   = RecordType("PTR")
	TypeSPF   = RecordType("SPF")
	TypeSRV   = RecordType("SRV")
	TypeSSHFP = RecordType("SSHFP")
	TypeTLSA  = RecordType("TLSA")
	TypeTXT   = RecordType("TXT")
	TypeWKS   = RecordType("WKS")
)

type zoneRecordJSON struct {
	Type       string   `json:"rrset_type,omitempty"`
	TimeToLive int      `json:"rrset_ttl,omitempty"`
	Name       string   `json:"rrset_name,omitempty"`
	Values     []string `json:"rrset_values,omitempty"`
}

type ZoneRecord struct {
	Type       RecordType
	TimeToLive time.Duration
	Name       string
	Values     []string
}

func Copy(z ZoneRecord) *ZoneRecord {
	return &ZoneRecord{
		Type:       z.Type,
		Name:       z.Name,
		TimeToLive: z.TimeToLive,
		Values:     z.Values,
	}
}

func (z ZoneRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(zoneRecordJSON{
		Type:       string(z.Type),
		Name:       z.Name,
		TimeToLive: int(math.Floor(z.TimeToLive.Seconds())),
		Values:     z.Values,
	})
}

func (z *ZoneRecord) UnmarshalJSON(b []byte) error {
	var jsonObj zoneRecordJSON
	if err := json.Unmarshal(b, &jsonObj); nil != err {
		return err
	}
	z.Type = RecordType(jsonObj.Type)
	z.Name = jsonObj.Name
	z.TimeToLive = time.Duration(jsonObj.TimeToLive) * time.Second
	z.Values = jsonObj.Values
	return nil
}

func (c Client) ListRecords(zoneID string) ([]ZoneRecord, error) {
	req, err := c.request("GET", defaultBaseUrl+"/zones/"+zoneID+"/records", strings.NewReader(""))
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
	records := make([]ZoneRecord, 0)
	if err = json.Unmarshal(body, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (c Client) ListRecordsAsText(zoneID string) (string, error) {
	req, err := c.request("GET", defaultBaseUrl+"/zones/"+zoneID+"/records", strings.NewReader(""))
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "text/plain")
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c Client) AddRecord(zoneID string, records []ZoneRecord) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	req, err := c.request("POST", defaultBaseUrl+"/zones/"+zoneID+"/records", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	cli := &http.Client{}
	_, err = cli.Do(req)
	return err
}

func (c Client) ChangeRecord(zoneID string, record ZoneRecord) error {
	toMarshal := Copy(record)
	toMarshal.Type = RecordType("")
	toMarshal.Name = ""
	b, err := json.Marshal(toMarshal)
	if err != nil {
		return err
	}
	req, err := c.request("PUT", "/zones/"+zoneID+"/records/"+record.Name+"/"+string(record.Type), bytes.NewReader(b))
	if err != nil {
		return err
	}
	cli := &http.Client{}
	_, err = cli.Do(req)
	return err
}
