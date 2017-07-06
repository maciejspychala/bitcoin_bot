package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
)

type MarketOrder struct {
	UUID string
}

type OpenOrder struct {
	OrderUUID string
	UUID      string
	Quantity  float64
	Price     float64
}

func (c *Client) Buy(market string, quantity, rate float64) (string, error) {
	paramMap := map[string]string{"market": market, "quantity": fmt.Sprintf("%.8f", quantity),
		"rate": fmt.Sprintf("%.8f", rate)}
	response, e := c.get("/market/buylimit", paramMap, "1.1")
	var o MarketOrder
	e = json.Unmarshal(response.Result, &o)
	return o.UUID, e
}

func (c *Client) Sell(market string, quantity, rate float64) (string, error) {
	paramMap := map[string]string{"market": market, "quantity": fmt.Sprintf("%.8f", quantity),
		"rate": fmt.Sprintf("%.8f", rate)}
	response, e := c.get("/market/selllimit", paramMap, "1.1")
	var o MarketOrder
	e = json.Unmarshal(response.Result, &o)
	return o.UUID, e
}

func (c *Client) Cancel(uuid string) error {
	response, e := c.get("/market/cancel", map[string]string{"uuid": uuid}, "1.1")
	if !response.Success {
		return errors.New(response.Message)
	}
	return e
}

func (c *Client) GetOpenOrders(market string) (o []OpenOrder, e error) {
	response, e := c.get("/market/getopenorders", map[string]string{"market": market}, "1.1")
	e = json.Unmarshal(response.Result, &o)
	return
}
