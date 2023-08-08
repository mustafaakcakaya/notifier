package main

import (
	"encoding/json"
	"fmt"
	"github.com/markcheno/go-talib"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

const (
	cryptoCurrency      = "BTC"
	currency            = "USDT"
	requestPeriodSecond = 1 * 60 * 15 // request every 15 min
	rsiPeriod           = 14
	bbPeriod            = 20
	maxDataCount        = 80
)

var lastTenAnalyses []bool

type AppData struct {
	ClosePrices          []float64 `json:"closePrices"`
	PreviousRsi          float64   `json:"previousRsi"`
	LastPrice            float64   `json:"lastPrice"`
	CurrentRsi           float64   `json:"currentRsi"`
	LastRequestTimestamp time.Time `json:"timestamp"`
}

var (
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUsername = "your.email@gmail.com"
	smtpPassword = "your_password"
)

func main() {
	//DONE: ensure data is getting from futures/binance, not spot/binance market
	//TODO: delete sound alert.
	//TODO: send signal alert when buy signal occurs.
	//TODO: test the algorithm with futures
	//TODO: connect to binance to auto buy and sell algorithm.
	//TODO: write test codes.
	PrintASCII()

	data, err := loadDataFromJSON()
	if err != nil {
		// If the file doesn't exist or there's an error, create a new data instance
		data = &AppData{
			ClosePrices: []float64{},
			PreviousRsi: 0,
		}
	}

	err, closeFunc := InitAudioPlayers()
	if err != nil {
		return
	}
	defer closeFunc()

	apiURL := fmt.Sprintf("https://fapi.binance.com/fapi/v1/ticker/24hr?symbol=%s%s", cryptoCurrency, currency)

	for {
		lastTenAnalyses = append(lastTenAnalyses, checkForBuySignal())

		// limit to 10 to check last 10 signal
		if len(lastTenAnalyses) > 10 {
			lastTenAnalyses = lastTenAnalyses[1:]
		}

		// Eğer son 10 analizde buy sinyali yoksa e-posta gönder
		if !containsBuySignal(lastTenAnalyses) {
			err := sendEmail(smtpHost, smtpPort, smtpUsername, smtpPassword)
			if err != nil {
				fmt.Println("Error sending email:", err)
			}
		}

		// Check if the last request was less than 15 minutes ago
		elapsedTime := time.Since(data.LastRequestTimestamp)
		remainingTime := requestPeriodSecond - int(elapsedTime.Seconds())

		// If less than 15 minutes ago, wait for the remaining time
		if remainingTime > 0 {
			fmt.Printf("Waiting %d seconds before making the next request...\n", remainingTime)
			time.Sleep(time.Duration(remainingTime) * time.Second)
		}

		// Make a GET request to Binance API
		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		// Unmarshal the JSON response
		var tickerData struct {
			LastPrice string `json:"lastPrice"`
		}
		if err := json.Unmarshal(body, &tickerData); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Parse the last price as a float64
		lastPrice, err := strconv.ParseFloat(tickerData.LastPrice, 64)
		if err != nil {
			fmt.Println("Error parsing last price:", err)
			return
		}

		// Perform analysis and print results
		if len(data.ClosePrices) >= rsiPeriod && len(data.ClosePrices) >= bbPeriod {
			// RSI calculation
			rsiData := talib.Rsi(data.ClosePrices, rsiPeriod)

			// RSI Divergence analysis:
			RsiDivergenceAnalysis(rsiData, lastPrice)

			// Update currentRsi and lastPrice in the data struct
			data.CurrentRsi = rsiData[len(rsiData)-1]
			data.LastPrice = lastPrice

		} else {
			fmt.Printf("Collecting Data... btc price: %.2f, time: %s\n", lastPrice, GetFormattedNow())
		}

		// Update the data with the current timestamp
		data.LastRequestTimestamp = time.Now()
		// Update the data with the current values
		data.ClosePrices = append(data.ClosePrices, lastPrice)
		data.PreviousRsi = previousRsi

		// Limit closePrices to closePriceLimit elements
		if len(data.ClosePrices) > maxDataCount {
			data.ClosePrices = data.ClosePrices[len(data.ClosePrices)-maxDataCount:]
		}

		// Save the updated data to the JSON file
		e := saveDataToJSON(data)
		if e != nil {
			fmt.Println("Error saving data to JSON:", err)
		}

		// Wait for a few seconds before making the next request
		time.Sleep(requestPeriodSecond * time.Second)
	}
}

var previousRsi float64 // Variable to store the last printed RSI value

func RsiDivergenceAnalysis(rsiData []float64, price float64) {
	// Get the current RSI value
	currentRsi := rsiData[len(rsiData)-1]

	// Get the current time
	now := GetFormattedNow()

	// Check for RSI divergence
	if currentRsi < 30 && currentRsi > previousRsi {
		fmt.Println(fmt.Sprintf("Buy signal, btc price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			price, currentRsi, previousRsi, now))
		ResumeBuyPlayer()
	} else if currentRsi > 70 && currentRsi < previousRsi {
		fmt.Println(fmt.Sprintf("Sell signal, btc price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			price, currentRsi, previousRsi, now))
	} else {
		fmt.Println(fmt.Sprintf("Just wait, btc price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
			price, currentRsi, previousRsi, now))
		PausePlayers()
	}

	// Update the previousRsi variable with the current RSI value
	previousRsi = currentRsi
}

func GetFormattedNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func InitAudioPlayers() (error, func()) {
	err := InitPlayers()
	if err != nil {
		fmt.Println("Error initializing audio player:", err)
		return err, nil
	}

	return err, ClosePlayers
}
func PrintASCII() {
	fmt.Println("")
	fmt.Println(" _______                                      ______           __ __           __    __            __     __  ______  __                   \n|       \\                                    /      \\         |  \\  \\         |  \\  |  \\          |  \\   |  \\/      \\|  \\                  \n| ▓▓▓▓▓▓▓\\__    __ __    __                 |  ▓▓▓▓▓▓\\ ______ | ▓▓ ▓▓         | ▓▓\\ | ▓▓ ______  _| ▓▓_   \\▓▓  ▓▓▓▓▓▓\\\\▓▓ ______   ______  \n| ▓▓__/ ▓▓  \\  |  \\  \\  |  \\     ______     | ▓▓___\\▓▓/      \\| ▓▓ ▓▓         | ▓▓▓\\| ▓▓/      \\|   ▓▓ \\ |  \\ ▓▓_  \\▓▓  \\/      \\ /      \\ \n| ▓▓    ▓▓ ▓▓  | ▓▓ ▓▓  | ▓▓    |      \\     \\▓▓    \\|  ▓▓▓▓▓▓\\ ▓▓ ▓▓         | ▓▓▓▓\\ ▓▓  ▓▓▓▓▓▓\\\\▓▓▓▓▓▓ | ▓▓ ▓▓ \\   | ▓▓  ▓▓▓▓▓▓\\  ▓▓▓▓▓▓\\\n| ▓▓▓▓▓▓▓\\ ▓▓  | ▓▓ ▓▓  | ▓▓     \\▓▓▓▓▓▓     _\\▓▓▓▓▓▓\\ ▓▓    ▓▓ ▓▓ ▓▓         | ▓▓\\▓▓ ▓▓ ▓▓  | ▓▓ | ▓▓ __| ▓▓ ▓▓▓▓   | ▓▓ ▓▓    ▓▓ ▓▓   \\▓▓\n| ▓▓__/ ▓▓ ▓▓__/ ▓▓ ▓▓__/ ▓▓                |  \\__| ▓▓ ▓▓▓▓▓▓▓▓ ▓▓ ▓▓         | ▓▓ \\▓▓▓▓ ▓▓__/ ▓▓ | ▓▓|  \\ ▓▓ ▓▓     | ▓▓ ▓▓▓▓▓▓▓▓ ▓▓      \n| ▓▓    ▓▓\\▓▓    ▓▓\\▓▓    ▓▓                 \\▓▓    ▓▓\\▓▓     \\ ▓▓ ▓▓         | ▓▓  \\▓▓▓\\▓▓    ▓▓  \\▓▓  ▓▓ ▓▓ ▓▓     | ▓▓\\▓▓     \\ ▓▓      \n \\▓▓▓▓▓▓▓  \\▓▓▓▓▓▓ _\\▓▓▓▓▓▓▓                  \\▓▓▓▓▓▓  \\▓▓▓▓▓▓▓\\▓▓\\▓▓          \\▓▓   \\▓▓ \\▓▓▓▓▓▓    \\▓▓▓▓ \\▓▓\\▓▓      \\▓▓ \\▓▓▓▓▓▓▓\\▓▓      \n                  |  \\__| ▓▓                                                                                                               \n                   \\▓▓    ▓▓                                                                                                               \n                    \\▓▓▓▓▓▓                                                                                                                \n\n \n\tBuy - Sell Notifier uygulamasına Hoşgeldiniz!\n\n\tBu program, belirtilen kripto para birimi (bitcoin) için gerçek zamanlı fiyatları alacak,\n\tRSI (Relative Strength Index) ve Bollinger Bands analizini gerçekleştirecek,\n\tve al/sat sinyalleri verecektir.\n\n\tAnaliz sonuçları, her veri güncellemesinde ekrana yazdırılacaktır.\n\n\tAnalizden çıkmak için programı durdurabilirsiniz.\n\n\tAnaliz Başladı...\n\t")
	//fmt.Println("  ____                    __  ____       _ _      _   _       _   _  __ _           \n | __ ) _   _ _   _      / / / ___|  ___| | |    | \\ | | ___ | |_(_)/ _(_) ___ _ __ \n |  _ \\| | | | | | |    / /  \\___ \\ / _ \\ | |    |  \\| |/ _ \\| __| | |_| |/ _ \\ '__|\n | |_) | |_| | |_| |   / /    ___) |  __/ | |    | |\\  | (_) | |_| |  _| |  __/ |   \n |____/ \\__,_|\\__, |  /_/    |____/ \\___|_|_|    |_| \\_|\\___/ \\__|_|_| |_|\\___|_|   \n              |___/                                                                 ")
	fmt.Println("")
}

// Function to load data from the JSON file
func loadDataFromJSON() (*AppData, error) {
	data := &AppData{}
	file, err := os.Open("data.json")
	if err != nil {
		return data, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(byteValue, data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Function to save data to the JSON file and append data if the file already exists
func saveDataToJSON(data *AppData) error {
	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Format the JSON for readability
	err = encoder.Encode(data)
	return err
}

func sendEmail(smtpHost, smtpPort, smtpUsername, smtpPassword string) error {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	to := []string{"alici.mail.adresi@gmail.com"} // E-posta alıcısı adresi
	msg := []byte("Subject: Buy Signal Notification\n\nNo buy signals in the last 10 analyses.")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUsername, to, msg)
	return err
}

func containsBuySignal(lastTenAnalyses []bool) bool {
	for _, buySignal := range lastTenAnalyses {
		if buySignal {
			return true
		}
	}
	return false
}

func checkForBuySignal() bool {
	// Burada buy sinyali kontrolü gerçekleştirilir, ve true/false döner
	// Eğer buy sinyali varsa true, yoksa false döner
	// Bu işlem sizin mevcut analiz mantığınıza göre gerçekleştirilmelidir
	return false
}
