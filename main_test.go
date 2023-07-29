package main

import (
	binance_connector "github.com/binance/binance-connector-go"
	"strconv"
	"testing"
)

func TestAlgorithm(t *testing.T) {
	// Simulate WebSocket data feed with test data
	testData := generateTestData()
	closePrices = nil // Reset closePrices slice before running the test
	for _, event := range testData {
		wsDepthHandler(event)
	}
}

// Helper function to generate test data for WebSocket events
func generateTestData() []*binance_connector.WsDepthEvent {
	testData := make([]*binance_connector.WsDepthEvent, 0)

	// Simulate WebSocket events
	for i := 0; i < 50; i++ {
		event := &binance_connector.WsDepthEvent{
			Bids: []binance_connector.Bid{
				{Price: strconv.FormatFloat(40+float64(i), 'f', 2, 64)},
			},
		}
		testData = append(testData, event)
	}

	return testData
}
