// SPDX-License-Identifier: Apache-2.0

package timedate

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

// TimeDate JSON message
type TimeDate struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

var timeDateMethods *share.Set
var timeInfo = map[string]string{}

// SetTimeDate set timedate property
func (t *TimeDate) SetTimeDate() error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	b := timeDateMethods.Contains(t.Property)
	if !b {
		return fmt.Errorf("Failed to set timedate:  %s not found", t.Property)
	}

	err = conn.SetTimeDate(t.Property, t.Value)
	if err != nil {
		log.Errorf("Failed to set timedate property: %s", err)
		return err
	}

	return nil
}

// GetTimeDate gets property from timedated
func GetTimeDate(rw http.ResponseWriter, property string) error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get dbus connection: %v", err)
		return err
	}
	defer conn.Close()

	for k := range timeInfo {
		p, err := conn.GetTimeDate(k)
		if err != nil {
			log.Errorf("Failed to get %s", k)
			continue
		}

		switch k {
		case "Timezone":
			v, b := p.Value().(string)
			if !b {
				continue
			}

			timeInfo[k] = v
			break
		case "LocalRTC":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "CanNTP":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "NTP":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "NTPSynchronized":
			v, b := p.Value().(bool)
			if !b {
				continue
			}

			timeInfo[k] = strconv.FormatBool(v)

			break
		case "TimeUSec":
			v, b := p.Value().(uint64)
			if !b {
				continue
			}

			t := time.Unix(0, int64(v))
			timeInfo[k] = t.String()

		case "RTCTimeUSec":
			v, b := p.Value().(uint64)
			if !b {
				continue
			}

			t := time.Unix(0, int64(v))

			timeInfo[k] = t.String()
			break
		}
	}

	if property == "" {
		return share.JSONResponse(timeInfo, rw)
	}

	t := TimeDate{
		Property: property,
		Value:    timeInfo[property],
	}

	return share.JSONResponse(t, rw)
}

// timeDateMethods init timedate package
func InitTimeDate() error {
	timeDateMethods = share.NewSet()

	timeDateMethods.Add("SetTime")
	timeDateMethods.Add("SetTimezone")
	timeDateMethods.Add("SetLocalRTC")
	timeDateMethods.Add("SetNTP")
	timeDateMethods.Add("ListTimezones")

	timeInfo["Timezone"] = ""
	timeInfo["LocalRTC"] = ""
	timeInfo["CanNTP"] = ""
	timeInfo["NTP"] = ""
	timeInfo["NTPSynchronized"] = ""
	timeInfo["TimeUSec"] = ""
	timeInfo["RTCTimeUSec"] = ""

	return nil
}
