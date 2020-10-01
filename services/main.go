package main

import (
	// "bufio"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// Complete url example http://api.weatherapi.com/v1/current.json?key=<YOUR_API_KEY>&q=London
	baseWeatherURL = "http://api.weatherapi.com/v1/current.json?"
)

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float32 `json:"lat"`
	Lon            float32 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Condition struct {
	Text string
	Icon string
	Code int
}

type Current struct {
	LastUpdatedEpoch int       `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float32   `json:"temp_c"`
	TempF            float32   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMph          float32   `json:"wind_mph"`
	WindKph          float32   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMb       float32   `json:"pressure_mb"`
	PressureIn       float32   `json:"pressure_in"`
	PrecipMm         float32   `json:"precip_mm"`
	PrecipIn         float32   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float32   `json:"feelslike_c"`
	FeelslikeF       float32   `json:"feelslike_f"`
	VisKm            float32   `json:"vis_km"`
	VisMiles         float32   `json:"vis_miles"`
	Uv               float32   `json:"uv"`
	GustMph          float32   `json:"gust_mph"`
	GustKph          float32   `json:"gust_kph"`
}

type Weather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Success  bool
}

func buildWeatherURL(key string, location string) (url string, err error) {
	if key == "" || location == "" {
		err = errors.New("args: args cannot be empty string.")
	}
	url = baseWeatherURL + "key=" + key + "&q=" + location
	return
}

func readAPIKey() (data string, err error) {
	read, err := ioutil.ReadFile("./api-key")
	if err != nil {
		return
	}
	r := string(read)
	data = strings.TrimSuffix(r, "\n")
	return
}

type LocationDetails struct {
	ZipCode string
}

func getWeatherByZip(zip string) (weather Weather, err error) {
	fmt.Println("weather")
	apiKey, err := readAPIKey()
	if err != nil {
		fmt.Println("error while reading API key", err)
		return
	}
	// fmt.Scanln(&zip)
	weatherURL, err := buildWeatherURL(apiKey, zip)
	fmt.Println("weatherURL", weatherURL)
	res, err := http.Get(weatherURL)
	if err != nil {
		fmt.Println("err", err)
	}
	defer res.Body.Close()
	fmt.Println("status", res.Status)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error while reading body", err)
	}
	weather = Weather{
		Success: true,
	}
	json.Unmarshal(body, &weather)
	fmt.Println("weather", weather)
	return
}

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		// details := LocationDetails{
		// 	ZipCode: r.FormValue("zipcode"),
		// }
		zip := r.FormValue("zipcode")
		weather, err := getWeatherByZip(zip)
		if err != nil {
			fmt.Println("Error while getting weather", err)
		}

		// do something with details

		tmpl.Execute(w, struct{ Weather Weather }{weather})
	})

	fmt.Println("listening on 8080")
	http.ListenAndServe(":8080", nil)

	// scanner := bufio.NewScanner(resp.Body)
	// for i := 0; scanner.Scan(); i++ {
	// 	fmt.Println(scanner.Text())
	// }

}
