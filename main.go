package main

import (
    . "bitcoin_bot/bittrex"
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

func sell(client *Client, boughtAt, sellAt float64, buyDate time.Time, market string) {
    sold := false
    for !sold {
        time.Sleep(120 * time.Second)
        orders, _ := client.GetMarketHistory(market)
        for _, order := range orders {
            if sellAt <= order.Price && buyDate.Before(order.TimeStamp.Time) {
                fmt.Printf("[%s] [SELL] date: %s\tprice: %12.8f\tboughtAt: %12.8f\n", market, formatDate(order.TimeStamp.Time), order.Price, boughtAt)
                sold = true
                break;
            }
        }
    }
}

func formatDate(d time.Time) string {
    return d.Format("2006-01-02 15:04:05")
}

func startBot(client *Client, market string) {
    for {
        ticks, _ := client.GetTicks(market, "fiveMin")
        ema24 := GetEMA(ticks, 24)
        ema48 := GetEMA(ticks, 48)
        fmt.Printf("[%s] 2h: %12.8f\t4h: %12.8f\n", market, ema24, ema48)
        history, _ := client.GetMarketHistory(market)
        wantBuyDate := history[0].TimeStamp
        wantBuyPrice := math.Min(ema24, ema48) * 0.9987
        fmt.Printf("[%s] wanted: date: %s\tprice: %12.8f\n", market, formatDate(wantBuyDate.Time), wantBuyPrice)
        for {
            orders, _ := client.GetMarketHistory(market)
            for _, order := range orders {
                if wantBuyPrice >= order.Price && wantBuyDate.Time.Before(order.TimeStamp.Time) {
                    sellAt := order.Price * 1.01
                    fmt.Printf("[%s] [BUY] date: %s\tprice: %12.8f\tsell at: %12.8f\n", market, formatDate(order.TimeStamp.Time), order.Price, sellAt)
                    go sell(client, order.Price, sellAt, order.TimeStamp.Time, market)
                    wantBuyDate = order.TimeStamp
                    break;
                }
            }

            latestTick, _ := client.GetLatestTick(market, "fiveMin")
            if !latestTick.T.Time.Equal(ticks[len(ticks)-1].T.Time) {
                break
            }

            time.Sleep(60 * time.Second)
        }
    }
}


func main() {
    key, secret := loadCredentials()
    client := NewClient(key, secret)
    displayWallets(client)
    var markets = [...]string {"BTC-XRP", "BTC-ANS", "BTC-SC", "BTC-NMR", "BTC-DASH"}
    for _, m := range markets {
        go startBot(client, m)
    }
    startBot(client, "BTC-MONA")
}
