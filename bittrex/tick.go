package bittrex

import (
    "fmt"
    "time"
    "encoding/json"
)



type Tick struct {
    Bid float64
    Ask float64
    Last float64
}


type TickInterval struct {
    O float64
    H float64
    L float64
    C float64
    V float64
    T TickTime
}

type TickTime struct {
    date time.Time
}

func (t *TickTime) UnmarshalJSON(b []byte) error {
    date, err := time.Parse("\"2006-01-02T15:04:05\"", string(b))
    t.date = date
    return err
}

func (t *TickInterval) ToString() string {
    return fmt.Sprintf("%s\tV:%15.8f\tO:%15.8f\tC:%15.8f", t.T.date.Format("15:04"), t.V, t.O, t.C)
}

func (c *Client) GetTick(market string) (t Tick, e error) {
    response, err := c.get("/public/getticker", map[string]string {"market" : market}, "1.1")
    err = json.Unmarshal(response.Result, &t)
    return t, err
}

func (c *Client) GetIntervalTicks(market string, interval string) (t []TickInterval, e error) {
    response, err := c.get("/pub/market/GetTicks",
        map[string]string {"marketName" : market, "tickInterval" : interval}, "2.0")
    err = json.Unmarshal(response.Result, &t)
    return t, err
}

