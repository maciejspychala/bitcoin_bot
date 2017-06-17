package main

import (
    "fmt"
    "time"
    "crypto/hmac"
    "crypto/sha512"
    "strings"
    "encoding/hex"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

const apiUrl = "https://bittrex.com/api/v1.1"

type jsonResponse struct {
    Success bool
    Message string
    Result json.RawMessage
}

type wallet struct {
    Currency string
    Balance float64
    Available float64
    Pending float64
    CryptoAddress string
}

type tick struct {
    Bid float64
    Ask float64
    Last float64
}

type client struct {
    apiKey string
    apiSecret string
    httpClient *http.Client
}

func (w *wallet) toString() (text string) {
    return fmt.Sprintf("%s\t%12.8f", w.Currency, w.Available)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func loadCredentials() (apiKey, secret string) {
    credentials, err := ioutil.ReadFile("credentials")
    check(err)
    cred := strings.Split(string(credentials), "\n")
    return cred[0], cred[1]
}

func newClient(apiKey, apiSecret string) (c *client) {
    return &client{apiKey, apiSecret, &http.Client{Timeout: 10 * time.Second}}
}

func (c *client) get(method string, params map[string]string) (jsonResp jsonResponse, e error) {
    req, err := http.NewRequest("GET", apiUrl + method, nil)
    check(err)

    req.Header.Add("Accept", "application/json")


    nonce := time.Now().UnixNano()
    q := req.URL.Query()

    for key, value := range params {
        q.Set(key, value)
    }

    q.Set("apikey", c.apiKey)
    q.Set("nonce", fmt.Sprintf("%d", nonce))

    req.URL.RawQuery = q.Encode()

    mac := hmac.New(sha512.New, []byte(c.apiSecret))
    _, err = mac.Write([]byte(req.URL.String()))
    sign := hex.EncodeToString(mac.Sum(nil))

    req.Header.Add("apisign", sign)


    resp, err := c.httpClient.Do(req)
    check(err)
    defer resp.Body.Close()

    respBody, _ := ioutil.ReadAll(resp.Body)
    err = json.Unmarshal(respBody, &jsonResp)
    return jsonResp, err
}

func (c *client) getWallets() (wallets []wallet, e error) {
    response, err := c.get("/account/getbalances", nil)
    check(err)
    err = json.Unmarshal(response.Result, &wallets)
    check(err)
    return wallets, err
}

func (c *client) getTick(market string) (t tick, e error) {
    response, err := c.get("/public/getticker", map[string]string {"market" : market})
    check(err)
    err = json.Unmarshal(response.Result, &t)
    check(err)
    return t, err
}

func main() {
    key, secret := loadCredentials()
    client := newClient(key, secret)
    wallets, err := client.getWallets()
    check(err)
    var wholeWalletValue float64
    fmt.Printf("%s\t%12s\t%12s\n", "name", "balance", "btc value")
    for _, w := range wallets {
        var btcValue float64
        if w.Currency == "BTC" {
            btcValue = w.Balance
        } else {
            t, _ := client.getTick("BTC-" + w.Currency)
            btcValue = t.Last * w.Balance
        }
        fmt.Printf("%s\t%12.8f\n", w.toString(), btcValue)
        wholeWalletValue += btcValue
    }
    fmt.Printf("\nwallet value : %12.8f btc\n", wholeWalletValue)
}
