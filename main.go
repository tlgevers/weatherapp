package main

import (
	// "bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// Complete url example http://api.weatherapi.com/v1/current.json?key=<YOUR_API_KEY>&q=London
	baseWeatherURL = "http://api.weatherapi.com/v1/current.json?"
)

type Location struct {
	Name            string
	Region          string
	Country         string
	Lat             float32
	Lon             float32
	Tz_id           string
	Localtime_epoch int
	Localtime       string
}

type Condition struct {
	Text string
	Icon string
	Code int
}

type Current struct {
	Last_updated_epoch int
	Last_updated       string
	Temp_c             float32
	Temp_f             float32
	Is_day             int
	Condition          Condition `json:"condition"`
	Wind_mph           float32
	Wind_kph           float32
	Wind_degree        int
	Wind_dir           string
	Pressure_mb        float32
	Pressure_in        float32
	Precip_mm          float32
	Precip_in          float32
	Humidity           int
	Cloud              int
	Feelslike_c        float32
	Feelslike_f        float32
	Vis_km             float32
	Vis_miles          float32
	Uv                 float32
	Gust_mph           float32
	Gust_kph           float32
}

type Weather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
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

func main() {
	fmt.Println("weather")
	apiKey, err := readAPIKey()
	if err != nil {
		fmt.Println("error while reading API key", err)
	}
	fmt.Println("Enter a zipcode")
	var zip string
	fmt.Scanln(&zip)
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
	weather := Weather{}
	json.Unmarshal(body, &weather)
	fmt.Println("weather", weather)

	// scanner := bufio.NewScanner(resp.Body)
	// for i := 0; scanner.Scan(); i++ {
	// 	fmt.Println(scanner.Text())
	// }

}
