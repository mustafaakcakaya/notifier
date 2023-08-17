package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/markcheno/go-talib"
	"io/ioutil"
	"os"
	"time"
)

func NewBuySellNotifier(pair string, config Config) (*BuySellNotifier, error) {
	var notifier = &BuySellNotifier{
		data:   &AppData{},
		config: &config,
		pair:   pair,
	}

	err := notifier.loadDataFromJSON()
	if err != nil {
		panic(err)
	}

	notifier.telegramBot, err = tgbotapi.NewBotAPI(notifier.config.BotToken)
	if err != nil {
		return nil, err
	}
	notifier.apiUrl = "https://fapi.binance.com/fapi/v1/ticker/24hr?symbol=" + notifier.pair
	return notifier, nil
}

func (bsn *BuySellNotifier) RsiDivergenceAnalysis(rsiData []float64, price float64) {
	// Get the current RSI value
	currentRsi := rsiData[len(rsiData)-1]

	// Get the current time
	now := bsn.GetFormattedNow()

	// Check for RSI divergence
	if currentRsi < 30 && currentRsi > bsn.data.PreviousRsi {
		bsn.message = fmt.Sprintf("Buy signal, %s price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			bsn.pair, price, currentRsi, bsn.data.PreviousRsi, now)
		msg := tgbotapi.NewMessage(bsn.config.ChatID, bsn.message)
		bsn.telegramBot.Send(msg)
		fmt.Println(bsn.message)
	} else if currentRsi > 70 && currentRsi < bsn.data.PreviousRsi {
		bsn.message = fmt.Sprintf("Sell signal, %s price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			bsn.pair, price, currentRsi, bsn.data.PreviousRsi, now)
		msg := tgbotapi.NewMessage(bsn.config.ChatID, bsn.message)
		bsn.telegramBot.Send(msg)
		fmt.Println(bsn.message)
	} else {
		fmt.Println(fmt.Sprintf("Just wait, %s price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			bsn.pair, price, currentRsi, bsn.data.PreviousRsi, now))
	}

	// Update the previousRsi variable with the current RSI value
	bsn.data.PreviousRsi = currentRsi
}
func (bsn *BuySellNotifier) Start() {
	for {
		bsn.WaitBeforeNextRequest()

		bsn.PerformAnalysis()
		bsn.SaveData()

		time.Sleep(requestPeriodSecond * time.Second)
	}
}

func (bsn *BuySellNotifier) WaitBeforeNextRequest() {
	elapsedTime := time.Since(bsn.data.LastRequestTimestamp)
	remainingTime := requestPeriodSecond - int(elapsedTime.Seconds())

	if remainingTime > 0 {
		fmt.Printf("Waiting %d seconds before making the next request...\n", remainingTime)
		time.Sleep(time.Duration(remainingTime) * time.Second)
	}
}
func (bsn *BuySellNotifier) PerformAnalysis() {
	// Make a GET request to Binance API
	lastPrice, err := GetPrice(bsn.apiUrl)
	if err != nil {
		fmt.Println("Error parsing last price:", err)
		return
	}

	// Perform analysis and print results
	if len(bsn.data.ClosePrices) >= rsiPeriod && len(bsn.data.ClosePrices) >= bbPeriod {
		// RSI calculation
		rsiData := talib.Rsi(bsn.data.ClosePrices, rsiPeriod)

		// RSI Divergence analysis:
		bsn.RsiDivergenceAnalysis(rsiData, lastPrice)

		// Update currentRsi and lastPrice in the data struct
		bsn.data.CurrentRsi = rsiData[len(rsiData)-1]
		bsn.data.LastPrice = lastPrice

	} else {
		fmt.Printf("Collecting Data... btc price: %.2f, time: %s\n", lastPrice, bsn.GetFormattedNow())
	}

	// Update the data with the current timestamp
	bsn.data.LastRequestTimestamp = time.Now()

	bsn.data.AppendClosePrice(lastPrice)

	bsn.data.PreviousRsi = bsn.data.PreviousRsi
}

func (bsn *BuySellNotifier) SaveData() {
	// Save the updated data to the JSON file
	err := bsn.saveDataToJSON(bsn.data)
	if err != nil {
		fmt.Println("Error saving data to JSON:", err)
	}
}

func (bsn *BuySellNotifier) GetFormattedNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (bsn *BuySellNotifier) loadDataFromJSON() error {
	file, err := os.Open("./data/" + bsn.pair + "_" + "data.json")
	if err != nil {
		// Handle the case where the file doesn't exist yet
		fmt.Printf("File doesn't exist for %s, creating...\n", bsn.pair)
		file, err = os.Create(bsn.pair + "_" + "data.json")
		if err != nil {
			return err
		}
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &bsn.data)
	if err != nil {
		return err
	}

	return nil
}

func (bsn *BuySellNotifier) saveDataToJSON(data *AppData) error {
	file, err := os.OpenFile("./data/"+bsn.pair+"_"+"data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Format the JSON for readability
	err = encoder.Encode(data)
	return err
}
