package main

import (
	"errors"
	"fmt"
	. "github.com/maciejspychala/bitcoin_bot/bittrex"
	"io/ioutil"
	"math"
	"strings"
	"time"
)

func loadCredentials() (apiKey, secret string) {
	credentials, _ := ioutil.ReadFile("credentials")
	cred := strings.Split(string(credentials), "\n")
	return cred[0], cred[1]
}

func displayWallets(c *Client) {
	wallets, _ := c.GetWallets()
	fmt.Printf("\n%-10s%15s %15s %15s %11s %15s %15s %15s\n",
		"name", "balance", "price", "prev day", "change", "24 high", "24 low", "btc value")
	summaries, _ := c.GetSummary()
	var wholeWalletValue float64
	for _, s := range summaries {
		for _, w := range wallets {
			if s.MarketName == "BTC-"+w.Currency && w.Balance > 0.0 {
				var btcValue float64
				btcValue = w.Balance * s.Last
				fmt.Printf("%-10s%15.8f %15.8f %15.8f %+9.2f %% %15.8f %15.8f %15.8f\n",
					w.Currency, w.Balance, s.Last, s.PrevDay, ((s.Last/s.PrevDay)-1.0)*100.0, s.High, s.Low, btcValue)
				wholeWalletValue += btcValue
			}
		}
	}
	for _, w := range wallets {
		if w.Currency == "BTC" {
			fmt.Printf("%-10s%15.8f\n", "BTC", w.Balance)
			wholeWalletValue += w.Balance
			break
		}
	}
	fmt.Printf("\nwallet value : %12.8f btc\n", wholeWalletValue)
}

func fakeBuy(client *Client, price float64, date time.Time, market string) (time.Time, error) {
	orders, _ := client.GetMarketHistory(market)
	for _, order := range orders {
		if price >= order.Price && date.Before(order.TimeStamp.Time) {
			fmt.Printf("[%s] [BUY] date: %s\tprice: %12.8f\n",
				market, formatDate(order.TimeStamp.Time), order.Price)
			return order.TimeStamp.Time, nil
		}
	}
	return time.Time{}, errors.New("No offers below price")
}

func fakeSell(client *Client, boughtAt, sellAt float64, buyDate time.Time, market string) {
	sold := false
	for !sold {
		time.Sleep(120 * time.Second)
		orders, _ := client.GetMarketHistory(market)
		for _, order := range orders {
			if sellAt <= order.Price && buyDate.Before(order.TimeStamp.Time) {
				fmt.Printf("[%s] [SELL] date: %s\tprice: %12.8f\tboughtAt: %12.8f\tp: %8.6f\n",
					market, formatDate(order.TimeStamp.Time), order.Price, boughtAt, order.Price/boughtAt)
				sold = true
				break
			}
		}
	}
}

func formatDate(d time.Time) string {
	return d.Format("2006-01-02 15:04:05")
}

func startFakeBot(client *Client, market string) {
	earnPercent := 1.015
	for {
		ticks, _ := client.GetTicks(market, "fiveMin")
		ema24 := GetEMA(ticks, 24)
		ema48 := GetEMA(ticks, 48)
		history, _ := client.GetMarketHistory(market)
		wantBuyDate := history[0].TimeStamp
		wantBuyPrice := math.Min(ema24, ema48) * 0.995
		minPrice := GetMinPriceFromLastOrders(history, 15)
		wantBuyPrice = math.Min(wantBuyPrice, minPrice)
		fmt.Printf("[%s] [WANT] date: %s\tprice: %12.8f\n", market, formatDate(wantBuyDate.Time), wantBuyPrice)
		for {
			time.Sleep(60 * time.Second)
			sellDate, err := fakeBuy(client, wantBuyPrice, wantBuyDate.Time, market)
			if err != nil {
				latestTick, _ := client.GetLatestTick(market, "fiveMin")
				if !latestTick.T.Time.Equal(ticks[len(ticks)-1].T.Time) {
					break
				}
				continue
			}
			wantSellPrice := wantBuyPrice * earnPercent
			fakeSell(client, wantBuyPrice, wantSellPrice, sellDate, market)
		}
	}
}

func isOrderCompleted(client *Client, market, uuid string) bool {
	orders, _ := client.GetOpenOrders(market)
	for _, o := range orders {
		if o.OrderUUID == uuid {
			return false
		}
	}
	return true
}

func waitForOrderCompleted(client *Client, market, uuid string) {
	for {
		time.Sleep(30 * time.Second)
		if isOrderCompleted(client, market, uuid) {
			break
		}
	}
}

func getBuyPrice(client *Client, market string) (float64, error) {
	ticks, _ := client.GetTicks(market, "fiveMin")
	ema24 := GetEMA(ticks, 24)
	ema48 := GetEMA(ticks, 48)
	history, err := client.GetMarketHistory(market)
	buyPrice := math.Min(ema24, ema48) * 0.995
	minPrice := GetMinPriceFromLastOrders(history, 15)
	buyPrice = math.Min(buyPrice, minPrice)
	return buyPrice, err
}

func startBot(client *Client, market string) {
	earnPercent := 1.015
	var buyUUID, sellUUID string
	for {
		if buyUUID != "" {
			err := client.Cancel(buyUUID)
			if err != nil {
				fmt.Printf("[%s] cannot cancel order\n", market)
				if sellUUID != "" {
					waitForOrderCompleted(client, market, sellUUID)
				}
			}
		}
		buyPrice, _ := getBuyPrice(client, market)
		tick, _ := client.GetLatestTick(market, "fiveMin")
		quantity := 0.008 / buyPrice
		buyUUID, _ = client.Buy(market, quantity, buyPrice)
		fmt.Printf("[%s] buyUUID: %s\n", market, buyUUID)

		for {
			time.Sleep(30 * time.Second)
			if isOrderCompleted(client, market, buyUUID) {
				sellPrice := buyPrice * earnPercent
				sellUUID, _ = client.Sell(market, quantity, sellPrice)
				waitForOrderCompleted(client, market, sellUUID)
				break
			}
			latestTick, _ := client.GetLatestTick(market, "fiveMin")
			if !latestTick.T.Time.Equal(tick.T.Time) {
				break
			}
		}
	}
}

func main() {
	key, secret := loadCredentials()
	client := NewClient(key, secret)
	var markets = [...]string{"BTC-NMR", "BTC-MEME"}
	for _, m := range markets {
		go startBot(client, m)
	}
	startBot(client, "BTC-OMNI")

}
