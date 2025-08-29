package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"os"
)

type weatherResponse struct {
	MinutelyWeather []minuteWeather `json:"minutely"`
	Timezone        string          `json:"timezone"`
}

type minuteWeather struct {
	DT            int64 `json:"dt"`
	Precipitation int   `json:"precipitation"`
}

const (
	lat  = 52.24922902394236
	long = 0.14061779650530742
)

func getForecast() (*weatherResponse, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY not set")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&appid=%s&exclude=current,hourly,daily,alerts", lat, long, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var weatherResponse weatherResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	fmt.Println(weatherResponse)

	return &weatherResponse, nil
}


func isHalfHourFromNow(targetTime time.Time) bool {
	return time.Until(targetTime) > 30 * time.Minute
}

func isWithinHalfHourFromNow(targetTime time.Time) bool {
	return time.Until(targetTime) < 30 * time.Minute
}

func ShouldGoHome() bool {
	weatherResponse, err := getForecast()
	if err != nil {
		fmt.Println(err)
	}
	minutelyForecast := weatherResponse.MinutelyWeather
	nextHalfHourDry := true
	for _, minuteForecast := range minutelyForecast {
		localTime := unixToLocal(minuteForecast.DT, weatherResponse.Timezone)
		fmt.Printf("Time: %v\n", localTime)
		if isWithinHalfHourFromNow(localTime) && minuteForecast.Precipitation == 0 {
			nextHalfHourDry = false
			fmt.Printf("Its dry at: %v\n", localTime)

		} else if isHalfHourFromNow(localTime) && nextHalfHourDry && minuteForecast.Precipitation != 0 {
			fmt.Printf("Go home! Gonna rain at: %v\n", localTime)
			return true
		}
	}
	fmt.Println("Get back to work you lazy git!")
	return false
}

func unixToLocal(unixTimestamp int64, timezone string) time.Time {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	utcTime := time.Unix(unixTimestamp, 0).UTC()
	return utcTime.In(location)
}
