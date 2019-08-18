package main

import (
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"
	"github.com/yryz/ds18b20"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfg Config

// Config Struct that holds all the variables needed for this simple app
type Config struct {
	MyDB      string
	Username  string
	Password  string
	OwmAPIKey string
	DbPort    string
	debug     bool
}

func getRpiTemp() float64 {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}
	tf := 0.0

	//    log.Infof("sensor IDs: %v\n", sensors)

	for _, sensor := range sensors {
		t, err := ds18b20.Temperature(sensor)
		if err == nil {
			tf = t*9.0/5.0 + 32.0
			//           log.Infof("Time: %s sensor: %s temperature: %.2f°F\n", time.Now(), sensor, tf)
		}
	}
	return tf
}

func getWeather() float64 {
	w, err := owm.NewCurrent("F", "EN", cfg.OwmAPIKey)
	if err != nil {
		//      log.Infof("Error configuring weather data")
		log.Fatal().Msgf("Fatal Error: %s", err)
	}
	if err := w.CurrentByZip(22304, "US"); err != nil {
		log.Panic().Msgf("Error retrieving weather: %s", err)

	}
	//log.Infof(w.Main.Temp)
	return w.Main.Temp
}

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: cfg.MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func writeData(t map[string]float64) {
	// Create a new HTTPClient
	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.DbPort,
		Username: cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		log.Fatal().Msgf("Fatal Error: %s", err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  cfg.MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal().Msgf("Fatal Error: %s", err)
	}

	for k, v := range t {
		tags := map[string]string{
			"key": k,
		}

		fields := map[string]interface{}{
			"value": v,
		}

		pt, err := client.NewPoint(
			"temperature",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Fatal().Msgf("Fatal Error: %s", err)
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := clnt.Write(bp); err != nil {
		log.Fatal().Msgf("Fatal Error: %s", err)
	}
	return
}

func readConfigs() {
	v := viper.New()
	v.SetConfigName("settings")
	v.AddConfigPath(".")
	v.AddConfigPath("/app")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal().Msgf("Fatal error config file: %s \n", err)
	}

	// Environment variables can be used to override the settings.toml file, but
	// must be prefixed with CS_
	v.SetEnvPrefix("CS")
	v.AutomaticEnv()

	debug := v.GetBool("debug")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	e := v.Unmarshal(&cfg)
	if e != nil {
		log.Warn().Msgf("couldn't read config: %s", e)
	}

	log.Debug().Msgf("%+v\n", cfg)
}

func main() {
	readConfigs()
	var temps map[string]float64
	temps = make(map[string]float64)
	temps["rpi"] = getRpiTemp()
	temps["wx"] = getWeather()
	writeData(temps)
}
