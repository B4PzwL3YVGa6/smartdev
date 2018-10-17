package monitoring

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anyandrea/smartdev/lib/database/weatherdb"
	"github.com/anyandrea/smartdev/lib/env"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/anyandrea/smartdev/lib/forecasts"
	"strings"
	"github.com/anyandrea/smartdev/lib/util"
)

func SpawnMonitoring(wdb weatherdb.WeatherDB) {
	go func(wdb weatherdb.WeatherDB) {
		//var alarmTime string
		var clientId string

		// check for VCAP_SERVICES first
		vcap, err := cfenv.Current()
		if err != nil {
			log.Println("Could not parse VCAP environment variables")
			log.Println(err)
		} else {
			service, err := vcap.Services.WithName("weathersms")
			if err != nil {
				log.Println("Could not find weathersms service in VCAP_SERVICES")
				log.Fatal(err)
			}
			clientId = fmt.Sprintf("%v", service.Credentials["client_id"])
		}

		// if WEATHERSMS_CLIENT_ID is not yet set then try to read it from ENV
		if len(clientId) == 0 {
			clientId = env.MustGet("WEATHERSMS_CLIENT_ID")
		}

		ticker := time.NewTicker(1 * time.Minute)
		for {
			<-ticker.C
			if isWindowOpen(wdb) {
				now := time.Now().In(util.Location)
				if now.Minute() == 30 && now.Hour() == 23 {
					sendSms(clientId, "Fenster-Alarm: Achtung, mind. ein Fenster ist noch geöffnet!")
				}

				forecast, _ := forecasts.Get(util.GetDefaultLocation("", ""))
				weather := strings.ToLower(forecast.Forecast.Tabular.Time[0].Symbol.Name)

				if strings.Contains(weather, "rain") ||
					strings.Contains(weather, "shower") ||
					strings.Contains(weather, "hail") || strings.Contains(weather, "sleet") {
					sendSms(clientId, "Regen-Alarm: Achtung, mind. ein Fenster ist noch geöffnet und es wird regnen!")
				}
			}
		}
	}(wdb)
}

func isWindowOpen(wdb weatherdb.WeatherDB) bool {
	windows, err := wdb.GetWindowStates()
	if err != nil {
		return true
	}
	for _, window := range windows {
		if window.Open {
			return true
		}
	}
	return false
}

func sendSms(clientId, text string) error {
	log.Println("Send window state alarm SMS: ", text)
	phoneNum := env.MustGet("WEATHERSMS_PHONENUM")
	var jsonStr = []byte(`{"to": "{` + phoneNum + `}", "text": "` + text + `"}`)
	req, err := http.NewRequest("POST", "https://api.swisscom.com/messaging/sms", bytes.NewBuffer(jsonStr))
	req.Header.Set("SCS-Version", "2")
	req.Header.Set("client_id", clientId)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	return nil
}

