package main

import (
	"fmt"

	"github.com/PouyaT/redis-weather-service-api/internal/client"
)

func main() {
	fmt.Println("Hello world")
	tempC, err := client.GetTemperature("30.3693824", "-97.6551936")
	if err != nil {
		fmt.Println("failed to get temperature", err)
		return
	}

	tempF := (tempC)*(9.0/5.0) + 32
	fmt.Printf("Temperature in F: %v\n", tempF)
	fmt.Printf("Temperature in C: %v\n", tempC)
}
