package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type TickerData struct {
	LastPrice string `json:"lastPrice"`
}

func GetPrice(apiUrl string) (float64, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, fmt.Errorf("Error fetching data: %s", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Error reading response: %s", err)
	}

	// Unmarshal the JSON response
	var tickerData TickerData
	if err := json.Unmarshal(body, &tickerData); err != nil {
		return 0, fmt.Errorf("Error decoding JSON: %s", err)
	}

	// Parse the last price as a float64
	lastPrice, err := strconv.ParseFloat(tickerData.LastPrice, 64)
	if err != nil {
		return 0, fmt.Errorf("Error parsing last price: %s", err)
	}

	return lastPrice, nil
}
