package bittrex

import (
    "fmt"
    "encoding/json"
)

type Summary struct {
    MarketName string
    High float64
    Low float64
    Volume float64
    Last float64
    BaseVolume float64
    TimeStamp string
    Bid float64
    Ask float64
    OpenBuyOrders int32
    OpenSellOrders int32
    PrevDay float64
    Created string
}

func (s *Summary) ToString() (text string) {
    return fmt.Sprintf("%-12s\t%15.8f\t%15.8f\tB:%15.8f\tA:%15.8f", s.MarketName, s.Last, s.PrevDay, s.Bid, s.Ask)
}

func (c *Client) GetSummary() (s []Summary, e error) {
    response, err := c.get("/public/getmarketsummaries", nil, "1.1")
    err = json.Unmarshal(response.Result, &s)
    return s, err
}

