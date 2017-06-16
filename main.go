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
)

const apiUrl = "https://bittrex.com/api/v1.1"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type client struct {
    apiKey string
    apiSecret string
    httpClient *http.Client
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

func (c *client) get(method string) (response []byte, e error) {
    req, err := http.NewRequest("GET", apiUrl + method, nil)
    check(err)

    req.Header.Add("Accept", "application/json")

    nonce := time.Now().UnixNano()
    q := req.URL.Query()
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
    return ioutil.ReadAll(resp.Body)
}

func main() {
    key, secret := loadCredentials()
    client := newClient(key, secret)
    response, err := client.get("/account/getbalances")
    check(err)
    fmt.Printf("%s\n", response)
}
