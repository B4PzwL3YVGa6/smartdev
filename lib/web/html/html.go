package html

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/anyandrea/smartdev/lib/config"
	"github.com/anyandrea/smartdev/lib/database/weatherdb"
	"github.com/anyandrea/smartdev/lib/forecasts"
	"github.com/anyandrea/smartdev/lib/web"
)

func Index(wdb weatherdb.WeatherDB, sm *scs.Manager) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Weather App",
			Active: "home",
		}
		web.Render().HTML(rw, http.StatusOK, "index", page)
	}
}

func NotFound(rw http.ResponseWriter, req *http.Request) {
	page := &Page{
		Title: "Weather App - Not Found",
	}
	web.Render().HTML(rw, http.StatusNotFound, "not_found", page)
}

func ErrorHandler(sm *scs.Manager) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		Error(sm, rw, req, fmt.Errorf("Internal Server Error"))
	}
}
func Error(sm *scs.Manager, rw http.ResponseWriter, req *http.Request, err error) {
	session := sm.Load(req)
	userName, _ := session.GetString("user_name")
	page := &Page{
		Title:   "Smartdev - Error",
		Content: err,
		User:    userName,
	}
	web.Render().HTML(rw, http.StatusInternalServerError, "error", page)
}

func Graphs(wdb weatherdb.WeatherDB, sm *scs.Manager) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Weather App - Graphs",
			Active: "graphs",
		}

		sensors, err := wdb.GetSensors()
		if err != nil {
			Error(sm, rw, req, err)
			return
		}

		var weeklyLabels, hourlyLabels []string
		weeklyTemperature := make(map[weatherdb.Sensor][]*weatherdb.SensorValue)
		hourlyTemperature := make(map[weatherdb.Sensor][]*weatherdb.SensorValue)
		weeklyHumidity := make(map[weatherdb.Sensor][]*weatherdb.SensorValue)
		hourlyHumidity := make(map[weatherdb.Sensor][]*weatherdb.SensorValue)
		for _, sensor := range sensors {
			switch sensor.Type {
			case "temperature":
				values, err := wdb.GetHourlyAverages(sensor.Id, 48)
				if err != nil {
					Error(sm, rw, req, err)
					return
				}
				hourlyTemperature[*sensor] = values

				if len(hourlyLabels) == 0 && sensor.Id == config.Get().Room.TemperatureSensorID {
					// collect labels
					for _, value := range values {
						hourlyLabels = append(hourlyLabels, value.Timestamp.Format("02.01. - 15:04"))
					}
				}

				values, err = wdb.GetDailyAverages(sensor.Id, 28)
				if err != nil {
					Error(sm, rw, req, err)
					return
				}
				weeklyTemperature[*sensor] = values

				if len(weeklyLabels) == 0 && sensor.Id == config.Get().Room.TemperatureSensorID {
					// collect labels
					for _, value := range values {
						weeklyLabels = append(weeklyLabels, value.Timestamp.Format("02.01.2006"))
					}
				}
			case "humidity":
				values, err := wdb.GetHourlyAverages(sensor.Id, 48)
				if err != nil {
					Error(sm, rw, req, err)
					return
				}
				hourlyHumidity[*sensor] = values

				values, err = wdb.GetDailyAverages(sensor.Id, 48)
				if err != nil {
					Error(sm, rw, req, err)
					return
				}
				weeklyHumidity[*sensor] = values
			}
		}

		page.Content = struct {
			HourlyTemperature map[weatherdb.Sensor][]*weatherdb.SensorValue
			HourlyHumidity    map[weatherdb.Sensor][]*weatherdb.SensorValue
			HourlyLabels      []string
			WeeklyTemperature map[weatherdb.Sensor][]*weatherdb.SensorValue
			WeeklyHumidity    map[weatherdb.Sensor][]*weatherdb.SensorValue
			WeeklyLabels      []string
			Config			  *config.Configuration
		}{
			hourlyTemperature,
			hourlyHumidity,
			hourlyLabels,
			weeklyTemperature,
			weeklyHumidity,
			weeklyLabels,
			config.Get(),
		}

		web.Render().HTML(rw, http.StatusOK, "graphs", page)
	}
}

func Sensors(wdb weatherdb.WeatherDB, sm *scs.Manager) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Weather App - Sensors",
			Active: "sensors",
		}

		session := sm.Load(req)
		userId, _ := session.GetInt("user_id")
		if userId == 0 {
			Unauthorized(rw)
			return
		}

		sensors, err := wdb.GetSensors()
		if err != nil {
			Error(sm, rw, req, err)
			return
		}

		data := make(map[int][]*weatherdb.SensorData, 0)
		for _, sensor := range sensors {
			d, err := wdb.GetSensorData(sensor.Id, 10)
			if err != nil {
				Error(sm, rw, req, err)
				return
			}
			data[sensor.Id] = d
		}

		page.Content = struct {
			Sensors    []*weatherdb.Sensor
			SensorData map[int][]*weatherdb.SensorData
		}{
			sensors,
			data,
		}
		web.Render().HTML(rw, http.StatusOK, "sensors", page)
	}
}

func Forecasts(rw http.ResponseWriter, req *http.Request) {
	canton, city := web.GetLocation(req)
	page := &Page{
		Title:  fmt.Sprintf("Weather App - Forecasts - %s", city),
		Active: "forecasts",
	}

	forecast, err := forecasts.Get(canton, city)
	if err != nil {
		page.Content = struct {
			Canton string
			City   string
			Error  error
		}{
			canton,
			city,
			err,
		}
		web.Render().HTML(rw, http.StatusNotFound, "forecast_error", page)
		return
	}

	page.Content = struct {
		Canton           string
		City             string
		Forecast         forecasts.WeatherForecast
		Today            time.Time
		Tomorrow         time.Time
		DayAfterTomorrow time.Time
	}{
		canton,
		city,
		forecast,
		time.Now(),
		time.Now().AddDate(0, 0, 1),
		time.Now().AddDate(0, 0, 2),
	}
	web.Render().HTML(rw, http.StatusOK, "forecasts", page)
}

