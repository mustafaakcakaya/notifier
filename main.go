package main

import (
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/markcheno/go-talib"
	"strconv"
	"time"
)

var (
	periodIntervalAsSecond int64 = 1
	periodCount            int64 = 20
	closePrices            []float64
	rsiPeriod              = 14 // RSI periyodunu değiştirdik (standart olarak 14 kullanılır)
	bbPeriod               = 20 // Bollinger Bands period
	bbStdDevUp             = 2.0
	bbStdDevDn             = 2.0
)

func main() {
	PrintASCII()
	//TODO: connect to binance to auto buy and sell algorithm.
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(false, "wss://testnet.binance.vision")

	cryptoCurrency := "BTC"
	currency := "USDT"

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
	// Get the best bid price as the last traded price
	if len(event.Bids) > 0 {
		lastPriceStr := event.Bids[0].Price
		lastPrice, _ := strconv.ParseFloat(lastPriceStr, 64) //assume always converted

		// Append closePrices (this is not actually closePrice, just data in a moment)
		closePrices = append(closePrices, lastPrice)

		// Let's limit the closePrices to 40 elements
		if len(closePrices) > 40 {
			closePrices = closePrices[len(closePrices)-40:]
		}

		// Perform analysis and print results
		if len(closePrices) >= rsiPeriod && len(closePrices) >= bbPeriod {
			// RSI calculation
			rsiData := talib.Rsi(closePrices, rsiPeriod)

			// Bollinger Bands calculation
			bbUpper, bbMiddle, bbLower := talib.BBands(closePrices, bbPeriod, bbStdDevUp, bbStdDevDn, talib.SMA)

			// Log results for the latest entry
			i := len(closePrices) - 1
			currentRsi := rsiData[i]
			fmt.Println(fmt.Sprintf("Close: %.2f, RSI: %.2f, BB Upper: %.2f, BB Middle: %.2f, BB Lower: %.2f",
				closePrices[i], currentRsi, bbUpper[i], bbMiddle[i], bbLower[i]))

			// RSI Divergence analysis:
			RsiDivergenceAnalysis(rsiData, lastPrice)
		}

		fmt.Println(fmt.Sprintf("Collecting Data... btc price: %.2f", lastPrice))
	}
}

func RsiDivergenceAnalysis(rsiData []float64, lastPrice float64) {
	for i := 1; i < len(rsiData); i++ {
		// RSI'nin önceki değeri
		previousRsi := rsiData[i-1]

		// RSI'nin şu anki değeri
		currentRsi := rsiData[i]

		// RSI Divergence analiz sonucu:
		if currentRsi > 70 && previousRsi <= 70 {
			// aşırı alım yapıldı.
			fmt.Println(fmt.Sprintf("sell signal, btc price:%.2f, currentRsi: %.2f, previousRsi: %.2f, time:%s", lastPrice, currentRsi, previousRsi, GetFormattedNow()))
		} else if currentRsi < 30 && previousRsi >= 30 {
			// aşırı satım yapıldı.
			fmt.Println(fmt.Sprintf("buy signal, btc price:%.2f, currentRsi: %.2f, previousRsi: %.2f, time:%s", lastPrice, currentRsi, previousRsi, GetFormattedNow()))
		} else {
			//nötr aralıkta
			fmt.Println(fmt.Sprintf("just wait, btc price:%.2f, currentRsi: %.2f, previousRsi: %.2f, time:%s", lastPrice, currentRsi, previousRsi, GetFormattedNow()))
		}
	}
}

func GetFormattedNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func PrintASCII() {
	fmt.Println("")
	fmt.Println(" _______                                      ______           __ __           __    __            __     __  ______  __                   \n|       \\                                    /      \\         |  \\  \\         |  \\  |  \\          |  \\   |  \\/      \\|  \\                  \n| ▓▓▓▓▓▓▓\\__    __ __    __                 |  ▓▓▓▓▓▓\\ ______ | ▓▓ ▓▓         | ▓▓\\ | ▓▓ ______  _| ▓▓_   \\▓▓  ▓▓▓▓▓▓\\\\▓▓ ______   ______  \n| ▓▓__/ ▓▓  \\  |  \\  \\  |  \\     ______     | ▓▓___\\▓▓/      \\| ▓▓ ▓▓         | ▓▓▓\\| ▓▓/      \\|   ▓▓ \\ |  \\ ▓▓_  \\▓▓  \\/      \\ /      \\ \n| ▓▓    ▓▓ ▓▓  | ▓▓ ▓▓  | ▓▓    |      \\     \\▓▓    \\|  ▓▓▓▓▓▓\\ ▓▓ ▓▓         | ▓▓▓▓\\ ▓▓  ▓▓▓▓▓▓\\\\▓▓▓▓▓▓ | ▓▓ ▓▓ \\   | ▓▓  ▓▓▓▓▓▓\\  ▓▓▓▓▓▓\\\n| ▓▓▓▓▓▓▓\\ ▓▓  | ▓▓ ▓▓  | ▓▓     \\▓▓▓▓▓▓     _\\▓▓▓▓▓▓\\ ▓▓    ▓▓ ▓▓ ▓▓         | ▓▓\\▓▓ ▓▓ ▓▓  | ▓▓ | ▓▓ __| ▓▓ ▓▓▓▓   | ▓▓ ▓▓    ▓▓ ▓▓   \\▓▓\n| ▓▓__/ ▓▓ ▓▓__/ ▓▓ ▓▓__/ ▓▓                |  \\__| ▓▓ ▓▓▓▓▓▓▓▓ ▓▓ ▓▓         | ▓▓ \\▓▓▓▓ ▓▓__/ ▓▓ | ▓▓|  \\ ▓▓ ▓▓     | ▓▓ ▓▓▓▓▓▓▓▓ ▓▓      \n| ▓▓    ▓▓\\▓▓    ▓▓\\▓▓    ▓▓                 \\▓▓    ▓▓\\▓▓     \\ ▓▓ ▓▓         | ▓▓  \\▓▓▓\\▓▓    ▓▓  \\▓▓  ▓▓ ▓▓ ▓▓     | ▓▓\\▓▓     \\ ▓▓      \n \\▓▓▓▓▓▓▓  \\▓▓▓▓▓▓ _\\▓▓▓▓▓▓▓                  \\▓▓▓▓▓▓  \\▓▓▓▓▓▓▓\\▓▓\\▓▓          \\▓▓   \\▓▓ \\▓▓▓▓▓▓    \\▓▓▓▓ \\▓▓\\▓▓      \\▓▓ \\▓▓▓▓▓▓▓\\▓▓      \n                  |  \\__| ▓▓                                                                                                               \n                   \\▓▓    ▓▓                                                                                                               \n                    \\▓▓▓▓▓▓                                                                                                                \n\n \n\tBuy - Sell Notifier uygulamasına Hoşgeldiniz!\n\n\tBu program, belirtilen kripto para birimi (bitcoin) için gerçek zamanlı fiyatları alacak,\n\tRSI (Relative Strength Index) ve Bollinger Bands analizini gerçekleştirecek,\n\tve al/sat sinyalleri verecektir.\n\n\tAnaliz sonuçları, her veri güncellemesinde ekrana yazdırılacaktır.\n\n\tAnalizden çıkmak için programı durdurabilirsiniz.\n\n\tAnaliz Başladı...\n\t")
	//fmt.Println("  ____                    __  ____       _ _      _   _       _   _  __ _           \n | __ ) _   _ _   _      / / / ___|  ___| | |    | \\ | | ___ | |_(_)/ _(_) ___ _ __ \n |  _ \\| | | | | | |    / /  \\___ \\ / _ \\ | |    |  \\| |/ _ \\| __| | |_| |/ _ \\ '__|\n | |_) | |_| | |_| |   / /    ___) |  __/ | |    | |\\  | (_) | |_| |  _| |  __/ |   \n |____/ \\__,_|\\__, |  /_/    |____/ \\___|_|_|    |_| \\_|\\___/ \\__|_|_| |_|\\___|_|   \n              |___/                                                                 ")
	//fmt.Println("                                                                 ▄▄    ▄▄                                              ▄▄     ▄▄▄▄ ▄▄                  \n▀███▀▀▀██▄                                     ▄█▀▀▀█▄█        ▀███  ▀███                 ▀███▄   ▀███▀          ██    ██   ▄█▀ ▀▀ ██                  \n  ██    ██                                    ▄██    ▀█          ██    ██                   ███▄    █            ██         ██▀                        \n  ██    █████  ▀███ ▀██▀   ▀██▀               ▀███▄     ▄▄█▀██   ██    ██                   █ ███   █   ▄██▀██▄██████▀███  █████ ▀███   ▄▄█▀██▀███▄███ \n  ██▀▀▀█▄▄ ██    ██   ██   ▄█                   ▀█████▄▄█▀   ██  ██    ██                   █  ▀██▄ █  ██▀   ▀██ ██    ██   ██     ██  ▄█▀   ██ ██▀ ▀▀ \n  ██    ▀█ ██    ██    ██ ▄█                  ▄     ▀████▀▀▀▀▀▀  ██    ██                   █   ▀██▄█  ██     ██ ██    ██   ██     ██  ██▀▀▀▀▀▀ ██     \n  ██    ▄█ ██    ██     ███                   ██     ████▄    ▄  ██    ██                   █     ███  ██▄   ▄██ ██    ██   ██     ██  ██▄    ▄ ██     \n▄████████  ▀████▀███▄   ▄█                    █▀█████▀  ▀█████▀▄████▄▄████▄               ▄███▄    ██   ▀█████▀  ▀████████▄████▄ ▄████▄ ▀█████▀████▄   \n                      ▄█                                                                                                                               \n                    ██▀                                                                                                                                \n")
	fmt.Println("")
}
