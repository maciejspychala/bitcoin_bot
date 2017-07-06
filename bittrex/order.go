package bittrex

import (
	"encoding/json"
	"fmt"
)

type Order struct {
	Quantity float64
	Rate     float64
}

func (c *Client) GetOrderBook(market, orderType string) (o []Order, e error) {
	response, e := c.get("/public/getorderbook",
		map[string]string{"market": market, "type": orderType}, "1.1")
	e = json.Unmarshal(response.Result, &o)
	return
}

func (o Order) String() string {
	return fmt.Sprintf("q: %15.8f\t r: %15.8f", o.Quantity, o.Rate)
}
