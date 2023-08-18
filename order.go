package main

import "fmt"

func (bsn *BuySellNotifier) openOrder(pair string, orderType string, price float64, currentRsi float64, previousRsi float64) {
	now := bsn.GetFormattedNow()
	bsn.message = fmt.Sprintf("%s signal, %s price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
		bsn.pair, orderType, price, currentRsi, previousRsi, now)
	fmt.Println(bsn.message)
	bsn.sendMessageToTelegram()

	openOrder := OpenOrder{OrderType: orderType, OpeningPrice: price, Leverage: 20}
	bsn.message = fmt.Sprintf("%s Order Opened opening price %.2f", orderType, price)
	fmt.Println(bsn.message)
	bsn.sendMessageToTelegram()
	bsn.OpenOrders[pair] = openOrder
}

func (bsn *BuySellNotifier) closeOrder(pair string, price float64, currentRsi float64, previousRsi float64) {
	openOrder := bsn.OpenOrders[pair]
	reserveOrder := reverseOrderType(bsn.OpenOrders[pair].OrderType)
	now := bsn.GetFormattedNow()
	bsn.message = fmt.Sprintf("%s signal, %s price: %.2f, currentRsi: %.2f, previousRsi: %.2f, time: %s",
		bsn.pair, reserveOrder, price, currentRsi, previousRsi, now)
	fmt.Println(bsn.message)
	bsn.sendMessageToTelegram()

	percentageDiff := float64(openOrder.Leverage) * (openOrder.OpeningPrice*100/(price) - 100)
	bsn.message = fmt.Sprintf("%s Order Closed opening price %.2f, closing price %.2f percentage diff %.2f", openOrder.OrderType, openOrder.OpeningPrice, price, percentageDiff)
	bsn.sendMessageToTelegram()
	fmt.Println(bsn.message)
	delete(bsn.OpenOrders, pair)
}

func reverseOrderType(orderType string) string {
	if orderType == "LONG" {
		return "SHORT"
	}
	return "LONG"
}

func (bsn *BuySellNotifier) isLiq(price float64, pair string) {
	openOrder := bsn.OpenOrders[pair]
	now := bsn.GetFormattedNow()
	percentageDiff := float64(openOrder.Leverage) * (openOrder.OpeningPrice*100/(price) - 100)
	if percentageDiff > 100 || percentageDiff < -100 {
		message := fmt.Sprintf("You Are Liq, order %s opening price: %.2f, current price: %.2f, percentage diff: %.2f, time: %s",
			openOrder.OrderType, openOrder.OpeningPrice, price, percentageDiff, now)
		fmt.Println(message)
		bsn.message = message
		bsn.sendMessageToTelegram()
		delete(bsn.OpenOrders, pair)
		fmt.Println(bsn.OpenOrders)
	}
}
