package main

import (
	"fmt"
	"github.com/markcheno/go-talib"
	gecko "github.com/superoo7/go-gecko/v3"
	"time"
)

func main() {
	PrintASCII()
	//TODO: client will be changed to binance
	api := gecko.NewClient(nil)
	cryptoCurrency := "bitcoin"
	currency := "usd"

	closePrices := make([]float64, 0)

	rsiPeriod := 14 // RSI periyodunu değiştirdik (standart olarak 14 kullanılır)
	bbPeriod := 20  // Bollinger Bands period
	// Bollinger Bands hesaplamak için kullanacağımız dönem ve sapma değeri
	bbStdDevUp := 2.0
	bbStdDevDn := 2.0

	for {
		// get real time data
		price, err := api.SimpleSinglePrice(cryptoCurrency, currency)
		if err != nil {
			fmt.Println("API err:", err)
			return
		}

		// log prices
		fmt.Println(fmt.Sprintf("btc price: %v, time:%s", price.MarketPrice, GetFormattedNow()))

		// append closePrices (this is not actually closePrice, just data in a moment)
		closePrices = append(closePrices, float64(price.MarketPrice))

		// En az gerekli veri miktarında ise analizi yapalım
		if len(closePrices) >= rsiPeriod && len(closePrices) >= bbPeriod {
			// RSI calculation
			rsiData := talib.Rsi(closePrices, rsiPeriod)

			// Bollinger Bands calculation
			bbUpper, bbMiddle, bbLower := talib.BBands(closePrices, bbPeriod, bbStdDevUp, bbStdDevDn, talib.SMA)

			// log results
			for i := 0; i < len(closePrices); i++ {
				fmt.Println(fmt.Sprintf("Close: %.2f, RSI: %.2f, BB Upper: %.2f, BB Middle: %.2f, BB Lower: %.2f\n",
					closePrices[i], rsiData[i], bbUpper[i], bbMiddle[i], bbLower[i]))
				fmt.Println(fmt.Sprintf("time:%s", GetFormattedNow()))
			}

			// RSI Divergence analysis:
			RsiDivergenceAnalysis(rsiData)
		}

		// wait for next
		time.Sleep(time.Second * 10)
	}
}

func RsiDivergenceAnalysis(rsiData []float64) {
	for i := 1; i < len(rsiData); i++ {
		// RSI'nin önceki değeri
		previousRsi := rsiData[i-1]

		// RSI'nin şu anki değeri
		currentRsi := rsiData[i]

		// RSI Divergence analiz sonucu:
		if currentRsi > 70 && previousRsi <= 70 {
			// aşırı alım yapıldı.
			fmt.Println(fmt.Sprintf("sell signal, currentRsi: %.2f, previousRsi: %.2f, time:%s", currentRsi, previousRsi, GetFormattedNow()))
		} else if currentRsi < 30 && previousRsi >= 30 {
			// aşırı satım yapıldı.
			fmt.Println(fmt.Sprintf("buy signal, currentRsi: %.2f, previousRsi: %.2f, time:%s", currentRsi, previousRsi, GetFormattedNow()))
		} else {
			//nötr aralıkta
			fmt.Println(fmt.Sprintf("just wait, currentRsi: %.2f, previousRsi: %.2f, time:%s", currentRsi, previousRsi, GetFormattedNow()))
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
