package main

import (
	"fmt"
	"io"
	"time"
	"net/http"
	"encoding/json"
	"os"
)

const (
	lat = 52.24922902394236
	long = 0.14061779650530742
)

type WeatherResponse struct {
	Forecast ForecastData `json:"forecast"`
}

type ForecastData struct {
	ForecastDay []DayData `json:"forecastday"`
}

type DayData struct {
	Hour []HourlyData `json:"hour"`
}

type HourlyData struct {
	Time         string  `json:"time"`
	ChanceOfRain int     `json:"chance_of_rain"`
	PrecipMM     float64 `json:"precip_mm"`
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04", timeStr)
}

func main() {
	apiKey := os.Getenv("WEATHERAPI_KEY")
	if apiKey == "" {
		fmt.Println("WEATHERAPI_KEY not set")
		return
	}
	resp, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=" + fmt.Sprintf("%f,%f", lat, long) + "&hours=3")
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("API error: %d\n", resp.StatusCode)
	}

	// TODO(human): Process the HTTP response to extract weather data
	// Hint: You need to read the response body and convert JSON to Go structs
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Raw response:", string(body))
	var weatherResponse WeatherResponse
	json.Unmarshal(body, &weatherResponse)
	fmt.Print(weatherResponse)
	todaysForecast := weatherResponse.Forecast.ForecastDay[0].Hour
	for _, hourData := range todaysForecast {
		hourTime, err := parseTime(hourData.Time)
		if err != nil {
			println(err)
			return
		}
		hour := hourTime.Hour()
		if hour >= 18 && hour <= 19 && hourData.ChanceOfRain >= 50 {
			fmt.Printf("Go home, might rain at %d\n", hour)
		}
	}
}


