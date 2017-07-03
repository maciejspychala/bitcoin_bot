package main

import (
    . "bitcoin/bittrex"
    "fmt"
    "math"
    "time"
    "strings"
    "io/ioutil"
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
            if s.MarketName == "BTC-" + w.Currency && w.Balance > 0.0 {
                var btcValue float64
                btcValue = w.Balance * s.Last
                fmt.Printf("%-10s%15.8f %15.8f %15.8f %+9.2f %% %15.8f %15.8f %15.8f\n",
                        w.Currency, w.Balance, s.Last, s.PrevDay, ((s.Last / s.PrevDay) - 1.0) * 100.0  ,s.High, s.Low, btcValue)
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


func main() {
    key, secret := loadCredentials()
    client := NewClient(key, secret)
    displayWallets(client)
    market := "BTC-XRP"
    for {
        ticks, _ := client.GetTicks(market, "fiveMin")
        ema24 := GetEMA(ticks, 24)
        ema48 := GetEMA(ticks, 48)
        fmt.Printf("2h: %12.8f\t4h: %12.8f\n", ema24, ema48)
        history, _ := client.GetMarketHistory(market)
        wantBuyDate := history[0].TimeStamp
        wantBuyPrice := math.Min(ema24, ema48) * 0.9987
        fmt.Printf("wanted: date: %v\tprice: %12.8f\n", wantBuyDate, wantBuyPrice)
        for {
            latestTick, _ := client.GetLatestTick(market, "fiveMin")
            if !latestTick.T.Time.Equal(ticks[len(ticks)-1].T.Time) {
                break
            }

            orders, _ := client.GetMarketHistory(market)
            for _, order := range orders {
                if wantBuyPrice >= order.Price && wantBuyDate.Time.Before(order.TimeStamp.Time) {
                    fmt.Printf("[BUY] get: date: %v\tprice: %12.8f\n", order.TimeStamp, order.Price)
                }
            }

            time.Sleep(20 * time.Second)
        }
    }
}
