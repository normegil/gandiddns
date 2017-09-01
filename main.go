package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/chyeh/pubip"
	gandi "github.com/normegil/gandidnsclient"
	flag "github.com/spf13/pflag"
)

var API_KEY string
var ZONE_NAME string
var RECORD_NAME string
var TIME_TO_LIVE int
var FORCED_IP string

func init() {
	flag.StringVarP(&API_KEY, "api-key", "k", "", "Key to access Gandi LiveDNS API (Get it in Gandic.net > Your profile > Security).")
	flag.StringVarP(&ZONE_NAME, "zone", "z", "", "Zone to be updated with the Public IP.")
	flag.StringVarP(&RECORD_NAME, "record", "n", "", "Record name to be updated with the Public IP.")
	flag.IntVarP(&TIME_TO_LIVE, "ttl", "t", 0, "New time to live for the changed record (Leave empty to keep original).")

	flag.StringVar(&FORCED_IP, "ip", "", "Use this option to update gandi record with the given IP")
	flag.Parse()
}

func main() {
	var ip net.IP
	if "" != FORCED_IP {
		ip = net.ParseIP(FORCED_IP)
		if nil == ip {
			panic(errors.New("String is not an IP: " + ip.String()))
		}
	} else {
		var err error
		ip, err = pubip.Get()
		if err != nil {
			panic(err)
		}
	}

	cli := gandi.NewClient(API_KEY)
	zone, err := getZone(*cli, ZONE_NAME)
	if err != nil {
		panic(err)
	} else if "" == zone.ID {
		fmt.Printf("Zone not found: %s", ZONE_NAME)
		os.Exit(1)
	}

	recordName := RECORD_NAME
	if RECORD_NAME == "" {
		recordName = "@"
	}

	ttl := time.Duration(TIME_TO_LIVE) * time.Second
	if 0 == ttl {
		record, err := getRecord(*cli, zone.ID, gandi.TypeA, recordName)
		if err != nil {
			panic(err)
		} else if 0 != record.TimeToLive.Seconds() {
			ttl = record.TimeToLive
		}
	}

	recordToSave := gandi.ZoneRecord{
		Type:       gandi.TypeA,
		Name:       recordName,
		Values:     []string{ip.String()},
		TimeToLive: ttl,
	}

	if err := cli.ChangeRecord(zone.ID, recordToSave); err != nil {
		panic(err)
	}
}

func getZone(cli gandi.Client, zoneName string) (gandi.Zone, error) {
	zones, err := cli.ListZones()
	if err != nil {
		return gandi.Zone{}, err
	}
	for _, zone := range zones {
		if zoneName == zone.Name {
			return zone, nil
		}
	}
	return gandi.Zone{}, nil
}

func getRecord(cli gandi.Client, zoneID string, recordType gandi.RecordType, name string) (gandi.ZoneRecord, error) {
	records, err := cli.ListRecords(zoneID)
	if err != nil {
		return gandi.ZoneRecord{}, err
	}
	for _, record := range records {
		if record.Name == name && record.Type == recordType {
			return record, nil
		}
	}
	return gandi.ZoneRecord{}, nil
}
