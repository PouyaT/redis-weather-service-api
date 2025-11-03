package main

import (
	"context"
	"fmt"

	"github.com/PouyaT/redis-weather-service-api/internal/client"
)

const REDIS_DEFAULT_ADDRESS = "localhost:6379"
const REDIS_DEFAULT_PWD = ""

func main() {
	lat := 30.36
	long := -97.65
	coordinatesKey := fmt.Sprintf("%.2f,%.2f", lat, long)
	tempC, err := client.GetTemperature(lat, long)
	if err != nil {
		fmt.Println("failed to get temperature", err)
		return
	}

	tempF := (tempC)*(9.0/5.0) + 32
	fmt.Printf("Temperature in F: %v\n", tempF)
	fmt.Printf("Temperature in C: %v\n", tempC)

	ctx := context.Background()
	redisClient := client.NewRedisClient(REDIS_DEFAULT_ADDRESS, REDIS_DEFAULT_PWD)
	defer redisClient.Close()

	status, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("failed to connect to redis", err)
	}
	fmt.Println(status)

	redisClient.Set(ctx, coordinatesKey, tempC, 0)
	val, _ := redisClient.Get(ctx, coordinatesKey).Float64()
	if err != nil {
		fmt.Printf("failed to get coordinates: %s, %v\n", coordinatesKey, err)
	}
	fmt.Println(val)
}
