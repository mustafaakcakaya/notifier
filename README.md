# Buy-Sell Notifier

Welcome to the Buy-Sell Notifier application! This program fetches real-time cryptocurrency prices, performs RSI (Relative Strength Index) and Bollinger Bands analysis, and provides buy/sell signals for a specified cryptocurrency (Bitcoin in this case).

## Introduction

The Buy-Sell Notifier is a Go program designed to assist traders and investors by analyzing cryptocurrency price trends and providing alerts when potential buy or sell opportunities arise based on RSI and Bollinger Bands indicators.

## Features

- Fetches real-time cryptocurrency prices from Binance Futures.
- Calculates the Relative Strength Index (RSI) and Bollinger Bands (BB) indicators.
- Provides buy and sell signals based on RSI divergence analysis.
- Sends Telegram notifications for buy and sell signals.
- Saves historical data and analysis results for future reference.

## Getting Started

1. Clone the repository to your local machine.
2. Make sure you have Go installed on your system.
3. Install the required third-party packages using the following command:

   ```sh
   go get -u github.com/go-telegram-bot-api/telegram-bot-api
   go get -u github.com/markcheno/go-talib
    ```
4. Rename data_sample.json to data.json and fill in your Telegram bot token and chat ID.
5. Run the program using the following command:
   ```sh
   go run main.go
   ```

## Usage
The program will continuously fetch data, analyze it, and provide buy/sell signals through Telegram notifications.

## Contributing
Contributions to this project are welcome! Feel free to submit issues and pull requests.

## License
This project is licensed under the MIT License.

## Disclaimer
Please note that cryptocurrency trading involves risks, and the signals provided by this program are not financial advice. Always do your own research and consider consulting with a financial advisor before making any trading decisions.
