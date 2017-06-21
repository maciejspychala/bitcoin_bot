package main

import (
    . "bitcoin/bittrex"
    "fmt"
    "sync"
    "strings"
    "io/ioutil"
)


func loadCredentials() (apiKey, secret string) {
    credentials, _ := ioutil.ReadFile("credentials")
    cred := strings.Split(string(credentials), "\n")
    return cred[0], cred[1]
}

func getCoinInfo(w Wallet, client *Client, ch chan<- float64) {
    var btcValue float64
    if w.Currency == "BTC" {
        btcValue = w.Balance
    } else {
        t, _ := client.GetTick("BTC-" + w.Currency)
        btcValue = t.Last * w.Balance
    }
    printMutex.Lock()
    fmt.Printf("%s\t%12.8f\n", w.ToString(), btcValue)
    printMutex.Unlock()
    ch<- btcValue
}

var printMutex sync.Mutex

func main() {
    key, secret := loadCredentials()
    client := NewClient(key, secret)
    wallets, _ := client.GetWallets()
    ch := make(chan float64)
    var wholeWalletValue float64
    fmt.Printf("%s\t%12s\t%12s\n", "name", "balance", "btc value")
    for _, w := range wallets {
        go getCoinInfo(w, client, ch)
    }
    for i := 0; i < len(wallets); i++ {
        value := <-ch
        wholeWalletValue += value
    }
    fmt.Printf("\nwallet value : %12.8f btc\n", wholeWalletValue)
}
