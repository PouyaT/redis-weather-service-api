package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// To get the current temperate of a location three requests are needed
// 1. Based on the lat and long use the /points/{lat},{long} endpoint to get a url with all the near by observationStations
// 2. With the response from step 1, the url given will return a list of urls for each observationStation
// 3. Query one of the observationStation, and in the body of properties.temperature there will be a Celcius value.

const BASE_URL = "https://api.weather.gov"

type latestStationsResponse struct {
	Properties struct {
		Temperature struct {
			Value *float64 `json:"value"`
		} `json:"temperature"`
	} `json:"properties"`
}

type pointsResponse struct {
	Properties struct {
		ObservationStationsURL string `json:"observationStations"`
	} `json:"properties"`
}

type observationStationsResponse struct {
	Stations []string `json:"observationStations"`
}

// Hits the points endpoint: https://www.weather.gov/documentation/services-web-api#/default/point
// returns the observationStations Url
func getPoints(lattitude string, longitude string) (string, error) {
	getPointsEndpoint := fmt.Sprintf("%s/points/%s,%s", BASE_URL, lattitude, longitude)
	resp, err := http.Get(getPointsEndpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var getPointsResponse pointsResponse
	err = json.Unmarshal(body, &getPointsResponse)
	if err != nil {
		return "", err
	}

	return getPointsResponse.Properties.ObservationStationsURL, nil
}

// Hits the observationStations Url from getPoints() which is the gridpoints stations endpoint: https://www.weather.gov/documentation/services-web-api#/default/gridpoint_stations
// returns a list of station endpoints and returns the first one:
func getStations(lattitude string, longitude string) (string, error) {
	observationStationsUrl, err := getPoints(lattitude, longitude)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(observationStationsUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var observationStations observationStationsResponse
	err = json.Unmarshal(body, &observationStations)
	if err != nil {
		return "", err
	}

	return observationStations.Stations[0], nil
}

// Using the station url returned from getStations() it hits the stations/obs/latest endpoint: https://www.weather.gov/documentation/services-web-api#/default/station_observation_latest
// returns a the temperature from the station in C
func GetTemperature(lattitude string, longitude string) (float64, error) {
	stationUrl, err := getStations(lattitude, longitude)
	if err != nil {
		return 0, nil
	}
	latestObsStationUrl := fmt.Sprintf("%s/observations/latest", stationUrl)
	resp, err := http.Get(latestObsStationUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var latestStation latestStationsResponse
	err = json.Unmarshal(body, &latestStation)
	if err != nil {
		return 0, err
	}
	if latestStation.Properties.Temperature.Value == nil {
		return 0, fmt.Errorf("no temperature data available for %s", stationUrl)
	}
	return *latestStation.Properties.Temperature.Value, nil

}
