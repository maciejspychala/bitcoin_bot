package bittrex

import (
    "fmt"
    "encoding/json"
)

type wallet struct {
    Currency string
    Balance float64
    Available float64
    Pending float64
    CryptoAddress string
}

func (w *wallet) ToString() (text string) {
    return fmt.Sprintf("%s\t%12.8f", w.Currency, w.Available)
}

func (c *Client) GetWallets() (wallets []wallet, e error) {
    response, err := c.get("/account/getbalances", nil, "1.1")
    err = json.Unmarshal(response.Result, &wallets)
    return wallets, err
}
