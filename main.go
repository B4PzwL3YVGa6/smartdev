package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	"github.com/alexedwards/scs"
	"github.com/anyandrea/smartdev/lib/config"
	"github.com/anyandrea/smartdev/lib/database"
	"github.com/anyandrea/smartdev/lib/database/weatherdb"
	"github.com/anyandrea/smartdev/lib/env"
	"github.com/anyandrea/smartdev/lib/forecasts"
	"github.com/anyandrea/smartdev/lib/util"
	"github.com/anyandrea/smartdev/lib/web/router"
	"github.com/urfave/negroni"
	"github.com/anyandrea/smartdev/lib/monitoring"
)

func main() {
	env.MustGet("WEATHERAPI_PASSWORD")
	wdb := setupDatabase()
	sm := setupSessionManager()

	// setup SIGINT catcher for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start a http server with negroni
	server := startHTTPServer(wdb, sm)

	// wait for SIGINT
	<-stop
	log.Println("Shutting down server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}




func setupDatabase() weatherdb.WeatherDB {
	// setup weather database
	adapter := database.NewAdapter()
	if err := adapter.RunMigrations("lib/database/migrations"); err != nil {
		if !strings.Contains(err.Error(), "no change") {
			log.Println("Could not run database migrations")
			log.Fatal(err)
		}
	}
	wdb := weatherdb.NewWeatherDB(adapter)

	// background jobs
	monitoring.SpawnMonitoring(wdb)
	spawnHousekeeping(wdb)
	spawnForecastCollection(wdb)

	return wdb
}

func spawnHousekeeping(wdb weatherdb.WeatherDB) {
	go func(wdb weatherdb.WeatherDB) {
		for {
			// retention policy of 33 days and minimum 50'000 values
			if err := wdb.Housekeeping(33, 50000); err != nil {
				log.Println("Database housekeeping failed")
				log.Fatal(err)
			}
			time.Sleep(12 * time.Hour)
		}
	}(wdb)
}

func spawnForecastCollection(wdb weatherdb.WeatherDB) {
	go func(wdb weatherdb.WeatherDB) {
		sensorId := config.Get().Forecast.TemperatureSensorID

		for {
			canton, city := util.GetDefaultLocation("", "")
			forecast, err := forecasts.Get(canton, city)
			if err != nil {
				log.Println("Weather forecast collection failed")
				log.Fatal(err)
			}

			if len(forecast.Forecast.Tabular.Time) > 0 {
				value, err := strconv.ParseInt(forecast.Forecast.Tabular.Time[0].Temperature.Value, 10, 64)
				if err != nil {
					log.Println("Could not read temperature value for forecast temperature")
					log.Fatal(err)
				}

				if err := wdb.InsertSensorValue(sensorId, int(value), time.Now()); err != nil {
					log.Println("Could not insert temperature value for forecast temperature")
					log.Fatal(err)
				}
				log.Printf("Weather forecast temperature:%v for %s/%s stored to database\n", value, canton, city)
			}

			time.Sleep(1 * time.Hour)
		}
	}(wdb)
}

func setupNegroni(wdb weatherdb.WeatherDB, sm *scs.Manager) *negroni.Negroni {
	n := negroni.Classic()

	r := router.New(wdb, sm)
	n.UseHandler(r)

	return n
}

func setupSessionManager() *scs.Manager {
	sm := scs.NewCookieManager("newsfeed-super-secret-cookie-key")
	sm.Lifetime(2 * time.Hour)
	sm.Persist(true)
	//sm.HttpOnly(false)
	//sm.Secure(true) // needs HTTPS
	return sm
}

func startHTTPServer(wdb weatherdb.WeatherDB, sm *scs.Manager) *http.Server {
	addr := ":" + env.Get("PORT", "8080")
	server := &http.Server{Addr: addr, Handler: setupNegroni(wdb, sm)}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}
