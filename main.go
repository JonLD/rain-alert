package main

import (
	"fmt"
	"rain-alert/notifications"
	"rain-alert/weather"
)

const (
	lat = 52.24922902394236
	long = 0.14061779650530742
)

func main() {
	if weather.ShouldGoHome() {
		err := notifications.ShowNotification()
		if err != nil {
			fmt.Println(err)
		}
	}

}


