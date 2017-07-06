package bittrex

import (
	"encoding/json"
	"fmt"
)

type Wallet struct {
	Currency      string
	Balance       float64
	Available     float64
	Pending       float64
	CryptoAddress string
}

func (w *Wallet) ToString() (text string) {
	return fmt.Sprintf("%s\t%12.8f", w.Currency, w.Balance)
}

func (c *Client) GetWallets() (wallets []Wallet, e error) {
	response, err := c.get("/account/getbalances", nil, "1.1")
	err = json.Unmarshal(response.Result, &wallets)
	return wallets, err
}
