package main

import (
    . "bitcoin/bittrex"
    "fmt"
    "strings"
    "io/ioutil"
)


func loadCredentials() (apiKey, secret string) {
    credentials, _ := ioutil.ReadFile("credentials")
    cred := strings.Split(string(credentials), "\n")
    return cred[0], cred[1]
}

func main() {
    key, secret := loadCredentials()
    client := NewClient(key, secret)
    wallets, _ := client.GetWallets()
    var wholeWalletValue float64
    fmt.Printf("%s\t%12s\t%12s\n", "name", "balance", "btc value")
    for _, w := range wallets {
        var btcValue float64
        if w.Currency == "BTC" {
            btcValue = w.Balance
        } else {
            t, _ := client.GetTick("BTC-" + w.Currency)
            btcValue = t.Last * w.Balance
        }
        fmt.Printf("%s\t%12.8f\n", w.ToString(), btcValue)
        wholeWalletValue += btcValue
    }
    fmt.Printf("\nwallet value : %12.8f btc\n", wholeWalletValue)

    ticks := make(map[string][]TickInterval)

    summaries, _ := client.GetSummary()
    for _, s := range summaries {
        market := s.MarketName
        fmt.Printf("%s\n", s.ToString())
        tempTicks, _ := client.GetIntervalTicks(market, "fiveMin")
        ticks[market] = append(ticks[market], tempTicks...)
        fmt.Printf("%s\n\n", ticks[market][len(ticks[market]) - 1].ToString())
    }
}
