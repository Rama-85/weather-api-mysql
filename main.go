package main

import (
	"database/sql"
	"encoding/json"
	f "fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type OpenMetoResponse struct {
	Latitude              float64      `json:"latitude"`
	Longitude             float64      `json:"longitude"`
	Generationtime_ms     float64      `json:"generationtime_ms"`
	Utc_offset_seconds    int          `json:"utc_offset_seconds"`
	Timezone              string       `json:"timezone"`
	Timezone_abbreviation string       `json:"timezone_abbrevation"`
	Elevation             float64      `json:"elevation"`
	Hourly_units          Hourly_units `json:"hourly_units"`
	Hourly                Hourly       `json:"hourly"`
}
type Hourly_units struct {
	Time                 string `json:"time"`
	Temperature_3m       string `json:"temperature_3m"`
	Relativehumidity_2m  string `json:"relativehumidity_2m"`
	Windspeed_10m        string `json:"windspeed_10m"`
	Winddirection_10m    string `json:"winddirection_10m"`
	Pressure_msl         string `json:"pressure_msl"`
	Soil_temperature_6cm string `json:"soil_temperature_6cm"`
	Visibility           string `json:"visibility"`
	Rain                 string `json:"rain"`
}
type Hourly struct {
	Time                 []string  `json:"time"`
	Temperature_3m       []float64 `json:"temperature_3m"`
	Relativehumidity_2m  []int     `json:"relativehumidity_2m"`
	Windspeed_10m        []float64 `json:"windspeed_10m"`
	Winddirection_10m    []int     `json:"winddirection_10m"`
	Pressure_msl         []float64 `json:"pressure_msl"`
	Soil_temperature_6cm []float64 `json:"soil_temperature_6cm"`
	Visibility           []float64 `json:"visibility"`
	Rain                 []float64 `json:"rain"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from GoLang!\n"))
}

func weather_data() (OpenMetoResponse, error) {
	var r OpenMetoResponse
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=11.7117117&longitude=79.3271609&timezone=IST&hourly=temperature_2m&hourly=relativehumidity_2m&hourly=windspeed_10m&hourly=winddirection_10m&hourly=pressure_msl&hourly=soil_temperature_6cm&hourly=visibility&hourly=rain")
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		return r, err
	}
	f.Println(r)
	return r, nil
}

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/weather")
	if err != nil {
		f.Println("error validating sql.open arguments")
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		f.Println("error verifying connection with db.ping")
		panic(err.Error())
	}
	insert, err := db.Query("INSERT INTO `weather`.`weather_data`(`temperature_3m`)VALUES(`temperature_3m`);")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	f.Println("successful connection to Database!")

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			loc, err := weather_data()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(loc)
		})
	http.ListenAndServe(":8080", nil)
}
