package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

const maxDataCount = 50

type BuySellNotifier struct {
	data         *AppData
	config       *Config
	telegramBot  *tgbotapi.BotAPI
	apiUrl       string
	pair         string
	previousRsi  float64
	lastRsiValue float64
	message      string
}

type AppData struct {
	ClosePrices          []float64 `json:"closePrices"`
	PreviousRsi          float64   `json:"previousRsi"`
	LastPrice            float64   `json:"lastPrice"`
	CurrentRsi           float64   `json:"currentRsi"`
	LastRequestTimestamp time.Time `json:"timestamp"`
}

func (data *AppData) AppendClosePrice(price float64) {
	data.ClosePrices = append(data.ClosePrices, price)
	if len(data.ClosePrices) > maxDataCount {
		data.ClosePrices = data.ClosePrices[len(data.ClosePrices)-maxDataCount:]
	}
}
