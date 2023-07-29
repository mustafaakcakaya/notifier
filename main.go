package main

import (
	"fmt"
	"github.com/binance/binance-connector-go"
	"github.com/markcheno/go-talib"
	"strconv"
	"time"
)

const (
	cryptoCurrency  = "BTC"
	currency        = "USDT"
	rsiPeriod       = 14
	bbPeriod        = 20
	bbStdDevUp      = 2.0
	bbStdDevDn      = 2.0
	closePriceLimit = 100
)

var closePrices []float64

func main() {
	PrintASCII()
	//TODO: ensure data is getting from futures, not spot market
	//TODO: test algorithm with futures
	//TODO: connect to binance to auto buy and sell algorithm.
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(false, "wss://testnet.binance.vision")

	errHandler := func(err error) {
		fmt.Println(err)
	}

	// Depth stream subscription
	doneCh, stopCh, err := websocketStreamClient.WsDepthServe(cryptoCurrency+currency, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		time.Sleep(30 * time.Second)
		stopCh <- struct{}{} // use stopCh to stop streaming
	}()

	<-doneCh
}

func wsDepthHandler(event *binance_connector.WsDepthEvent) {
	if len(event.Bids) == 0 {
		return
	}

	lastPriceStr := event.Bids[0].Price
	lastPrice, _ := strconv.ParseFloat(lastPriceStr, 64)

	// Append closePrices (this is not actually closePrice, just data in a moment)
	closePrices = append(closePrices, lastPrice)

	// Limit closePrices to closePriceLimit elements
	if len(closePrices) > closePriceLimit {
		closePrices = closePrices[len(closePrices)-closePriceLimit:]
	}

	// Perform analysis and print results
	if len(closePrices) >= rsiPeriod && len(closePrices) >= bbPeriod {
		// RSI calculation
		rsiData := talib.Rsi(closePrices, rsiPeriod)

		/*// Bollinger Bands calculation
		bbUpper, bbMiddle, bbLower := talib.BBands(closePrices, bbPeriod, bbStdDevUp, bbStdDevDn, talib.SMA)

		// Log results for the latest entry
		i := len(closePrices) - 1
		currentRsi := rsiData[i]
		fmt.Printf("Close: %.2f, RSI: %.2f, BB Upper: %.2f, BB Middle: %.2f, BB Lower: %.2f\n", closePrices[i], currentRsi, bbUpper[i], bbMiddle[i], bbLower[i])*/

		// RSI Divergence analysis:
		RsiDivergenceAnalysis(rsiData, lastPrice)
	} else {
		fmt.Printf("Collecting Data... btc price: %.2f\n", lastPrice)
	}
}

var lastRsi float64 // Variable to store the last printed RSI value

func RsiDivergenceAnalysis(rsiData []float64, price float64) {
	// Get the current RSI value
	currentRsi := rsiData[len(rsiData)-1]

	// Get the current time
	now := GetFormattedNow()

	// Check for RSI divergence
	if currentRsi < 30 && currentRsi > lastRsi {
		logMessage := fmt.Sprintf("Buy signal, btc price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			price, currentRsi, lastRsi, now)
		fmt.Println(logMessage)
	} else if currentRsi > 70 && currentRsi < lastRsi {
		logMessage := fmt.Sprintf("Sell signal, btc price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			price, currentRsi, lastRsi, now)
		fmt.Println(logMessage)
	} else {
		logMessage := fmt.Sprintf("Just wait, btc price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			price, currentRsi, lastRsi, now)
		fmt.Println(logMessage)
	}

	// Update the lastRsi variable with the current RSI value
	lastRsi = currentRsi
}

func GetFormattedNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func PrintASCII() {
	fmt.Println("")
	fmt.Println(" _______                                      ______           __ __           __    __            __     __  ______  __                   \n|       \\                                    /      \\         |  \\  \\         |  \\  |  \\          |  \\   |  \\/      \\|  \\                  \n| ▓▓▓▓▓▓▓\\__    __ __    __                 |  ▓▓▓▓▓▓\\ ______ | ▓▓ ▓▓         | ▓▓\\ | ▓▓ ______  _| ▓▓_   \\▓▓  ▓▓▓▓▓▓\\\\▓▓ ______   ______  \n| ▓▓__/ ▓▓  \\  |  \\  \\  |  \\     ______     | ▓▓___\\▓▓/      \\| ▓▓ ▓▓         | ▓▓▓\\| ▓▓/      \\|   ▓▓ \\ |  \\ ▓▓_  \\▓▓  \\/      \\ /      \\ \n| ▓▓    ▓▓ ▓▓  | ▓▓ ▓▓  | ▓▓    |      \\     \\▓▓    \\|  ▓▓▓▓▓▓\\ ▓▓ ▓▓         | ▓▓▓▓\\ ▓▓  ▓▓▓▓▓▓\\\\▓▓▓▓▓▓ | ▓▓ ▓▓ \\   | ▓▓  ▓▓▓▓▓▓\\  ▓▓▓▓▓▓\\\n| ▓▓▓▓▓▓▓\\ ▓▓  | ▓▓ ▓▓  | ▓▓     \\▓▓▓▓▓▓     _\\▓▓▓▓▓▓\\ ▓▓    ▓▓ ▓▓ ▓▓         | ▓▓\\▓▓ ▓▓ ▓▓  | ▓▓ | ▓▓ __| ▓▓ ▓▓▓▓   | ▓▓ ▓▓    ▓▓ ▓▓   \\▓▓\n| ▓▓__/ ▓▓ ▓▓__/ ▓▓ ▓▓__/ ▓▓                |  \\__| ▓▓ ▓▓▓▓▓▓▓▓ ▓▓ ▓▓         | ▓▓ \\▓▓▓▓ ▓▓__/ ▓▓ | ▓▓|  \\ ▓▓ ▓▓     | ▓▓ ▓▓▓▓▓▓▓▓ ▓▓      \n| ▓▓    ▓▓\\▓▓    ▓▓\\▓▓    ▓▓                 \\▓▓    ▓▓\\▓▓     \\ ▓▓ ▓▓         | ▓▓  \\▓▓▓\\▓▓    ▓▓  \\▓▓  ▓▓ ▓▓ ▓▓     | ▓▓\\▓▓     \\ ▓▓      \n \\▓▓▓▓▓▓▓  \\▓▓▓▓▓▓ _\\▓▓▓▓▓▓▓                  \\▓▓▓▓▓▓  \\▓▓▓▓▓▓▓\\▓▓\\▓▓          \\▓▓   \\▓▓ \\▓▓▓▓▓▓    \\▓▓▓▓ \\▓▓\\▓▓      \\▓▓ \\▓▓▓▓▓▓▓\\▓▓      \n                  |  \\__| ▓▓                                                                                                               \n                   \\▓▓    ▓▓                                                                                                               \n                    \\▓▓▓▓▓▓                                                                                                                \n\n \n\tBuy - Sell Notifier uygulamasına Hoşgeldiniz!\n\n\tBu program, belirtilen kripto para birimi (bitcoin) için gerçek zamanlı fiyatları alacak,\n\tRSI (Relative Strength Index) ve Bollinger Bands analizini gerçekleştirecek,\n\tve al/sat sinyalleri verecektir.\n\n\tAnaliz sonuçları, her veri güncellemesinde ekrana yazdırılacaktır.\n\n\tAnalizden çıkmak için programı durdurabilirsiniz.\n\n\tAnaliz Başladı...\n\t")
	//fmt.Println("  ____                    __  ____       _ _      _   _       _   _  __ _           \n | __ ) _   _ _   _      / / / ___|  ___| | |    | \\ | | ___ | |_(_)/ _(_) ___ _ __ \n |  _ \\| | | | | | |    / /  \\___ \\ / _ \\ | |    |  \\| |/ _ \\| __| | |_| |/ _ \\ '__|\n | |_) | |_| | |_| |   / /    ___) |  __/ | |    | |\\  | (_) | |_| |  _| |  __/ |   \n |____/ \\__,_|\\__, |  /_/    |____/ \\___|_|_|    |_| \\_|\\___/ \\__|_|_| |_|\\___|_|   \n              |___/                                                                 ")
	fmt.Println("")
}
