package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/PouyaT/redis-weather-service-api/internal/client"
)

const REDIS_DEFAULT_ADDRESS = "localhost:6379"
const REDIS_DEFAULT_PWD = ""

var tmpl = template.Must(template.ParseFiles("web/templates/index.html"))

type PageData struct {
	Lat, Lon     float64
	TempC, TempF float64
	Error        string
}

// handles HTTP requests for the weather page
// Supports GET and POST requests
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// Handles GET requests by displaying the default index template
	case "GET":
		tmpl.Execute(w, nil)
	// Handles POST requests using data from the lat and lon form values,
	case "POST":
		ctx := r.Context()
		// Reads form values
		lat := r.FormValue("lat")
		lon := r.FormValue("lon")
		var data PageData

		// Parses lat and long to make sure they're floats
		latFloat, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			fmt.Println("failed to parse lattiude", err)
			data.Error = fmt.Sprintf("Failed to parse lattiude: %v", err)
			tmpl.Execute(w, data)
			return
		}
		lonFloat, err := strconv.ParseFloat(lon, 64)
		if err != nil {
			fmt.Println("failed to parse longitude", err)
			data.Error = fmt.Sprintf("Failed to parse longitude: %v", err)
			tmpl.Execute(w, data)
			return
		}
		data.Lat = latFloat
		data.Lon = lonFloat

		// Creates redis cache key based off the user inputs
		coordinatesKey := fmt.Sprintf("%.2f,%.2f", data.Lat, data.Lon)
		// Checks if the cache key exists and if so returns early
		cachedTemp, err := client.RedisClient.Get(ctx, coordinatesKey).Float64()
		if err == nil {
			data.TempC = cachedTemp
			data.TempF = cachedTemp*(9.0/5.0) + 32
			fmt.Println("Using redis cached value")
			tmpl.Execute(w, data)
			return
		}
		// If cache key doesn't exists then get the data from the nws
		tempC, err := client.GetTemperature(data.Lat, data.Lon)
		if err != nil {
			fmt.Println("failed to get temperature", err)
			data.Error = fmt.Sprintf("Failed to get temperature: %v", err)
		}
		// Store the weather data
		client.RedisClient.Set(ctx, coordinatesKey, tempC, 0)
		// Set the weather data and submit it
		tempF := (tempC)*(9.0/5.0) + 32
		data.TempC = tempC
		data.TempF = tempF
		tmpl.Execute(w, data)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	if err := client.InitRedis(REDIS_DEFAULT_ADDRESS, REDIS_DEFAULT_PWD); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handler)
	fmt.Println("Server started on localhost:8080")
	http.ListenAndServe(":8080", nil)

}
