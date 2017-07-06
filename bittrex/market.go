package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
)

type MarketOrder struct {
	uuid string
}

func (c *Client) Buy(market string, quantity, rate float64) (o MarketOrder, e error) {
	paramMap := map[string]string{"market": market, "quantity": fmt.Sprintf("%.8f", quantity),
		"rate": fmt.Sprintf("%.8f", rate)}
	response, e := c.get("/market/buylimit", paramMap, "1.1")
	e = json.Unmarshal(response.Result, &o)
	return
}

func (c *Client) Sell(market string, quantity, rate float64) (o MarketOrder, e error) {
	paramMap := map[string]string{"market": market, "quantity": fmt.Sprintf("%.8f", quantity),
		"rate": fmt.Sprintf("%.8f", rate)}
	response, e := c.get("/market/selllimit", paramMap, "1.1")
	e = json.Unmarshal(response.Result, &o)
	return
}

func (c *Client) Cancel(uuid string) error {
	response, e := c.get("/market/cancel", map[string]string{"uuid": uuid}, "1.1")
	if !response.Success {
		return errors.New(response.Message)
	}
	return e
}
