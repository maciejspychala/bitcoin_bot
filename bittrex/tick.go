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
    Avg float64
}

type TickTime struct {
    Date time.Time
}

func (t *TickTime) UnmarshalJSON(b []byte) error {
    date, err := time.Parse("\"2006-01-02T15:04:05\"", string(b))
    t.Date = date
    return err
}

func (t *TickInterval) ToString() string {
    return fmt.Sprintf("%s\tV:%15.8f\tO:%15.8f\tC:%15.8f\tA:%15.8f", t.T.Date.Format("15:04"), t.V, t.O, t.C, t.Avg)
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

func CountAverages(ticks []TickInterval) {
    for i := 0; i < len(ticks); i++ {
        ticks[i].Avg = (ticks[i].O + ticks[i].C) / 2
    }
}

func HighestInLastDay(ticks []TickInterval) float64 {
    now := ticks[len(ticks) - 1].T.Date
    max := ticks[len(ticks) - 3].Avg
    for i:= len(ticks) - 3; i >= 0 && now.Sub(ticks[i].T.Date).Hours() <= 24.0; i-- {
        cur := ticks[i].Avg
        if cur > max {
            max = cur
        }
    }
    return max
}
