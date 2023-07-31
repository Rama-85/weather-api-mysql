package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type celsius float64
type km float64
type mm int

type Airtemp struct {
	maximum, minimum celsius
}

type Roadtemp struct {
	maximum, minimum celsius
}

type Visibility struct {
	maximum, minimum km
}

type Rain struct {
	maximum, minimum mm
}

type WeatherData struct {
	gorm.Model
	Name string `json:"name"`
	Main struct {
		Airtemp       `json:"Airtemp"`
		AirPressure   float64 `json:"AirPressure"`
		Humidity      float64 `json:"Humidity"`
		WindDirection string  `json:"WindDirection"`
		Roadtemp      `json:"Roadtemp"`
		Visibility    `json:"Visibility"`
		Windspeed     float64 `json:"Windspeed"`
		Rain          `json:"Rain"`
	} `json:"main"`
}

func GetWeather(db *gorm.DB, wd *[]WeatherData) (err error) {
	err = db.Find(wd).Error
	if err != nil {
		return err
	}
	return nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from GoLang!\n"))
}

func query(city string) (WeatherData, error) {

	resp, err := http.Get("http://api.open-meteo.com/v1/forecast?" + "&q=" + city)
	if err != nil {
		return WeatherData{}, err
	}

	defer resp.Body.Close()

	var d WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}
	return d, nil

}

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})

	http.ListenAndServe(":8080", nil)
}
