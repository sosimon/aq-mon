package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var apiURL = "https://api.airvisual.com/v2/city"

type airVisualResponse struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}

type data struct {
	Current currentData `json:"current"`
}

type currentData struct {
	Weather   weatherData   `json:"weather"`
	Pollution pollutionData `json:"pollution"`
}

type weatherData struct {
	Timestamp     string  `json:"ts"`
	Temperature   int     `json:"tp"`
	Humidity      int     `json:"hu"`
	Pressure      int     `json:"pr"`
	Windspeed     float32 `json:"ws"`
	WindDirection int     `json:"wd"`
}

type pollutionData struct {
	Timestamp string `json:"ts"`
	AQI       int    `json:"aqius"`
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	fmt.Println("Air Quality Monitor - www.iqair.com")
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("Missing environment variable API_KEY")
	}

	params := url.Values{}
	params.Add("city", getenv("CITY", "San-Francisco"))
	params.Add("state", getenv("STATE", "California"))
	params.Add("country", getenv("COUNTRY", "USA"))
	params.Add("key", apiKey)

	reqURL := apiURL + "?" + params.Encode()
	fmt.Println("Request URL: " + reqURL)
	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var response airVisualResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %s\n", err)
	}
	fmt.Printf("Unmarshalled response: %+v\n", response)
	fmt.Printf("Air quality index: %d\n", response.Data.Current.Pollution.AQI)
}
