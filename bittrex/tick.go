package bittrex

import (
    "time"
    "encoding/json"
)

type TickInterval struct {
    O float64
    H float64
    L float64
    C float64
    V float64
    T TickTime
    BV float64
}

type TickTime struct {
    time.Time
}

func (t *TickTime) UnmarshalJSON(b []byte) error {
    date, err := time.Parse("\"2006-01-02T15:04:05\"", string(b))
    *t = TickTime{date}
    return err
}

func (c *Client) GetTicks(market string, interval string) (t []TickInterval, e error) {
    response, err := c.get("/pub/market/GetTicks",
        map[string]string {"marketName" : market, "tickInterval" : interval}, "2.0")
    err = json.Unmarshal(response.Result, &t)
    return t, err
}

func (c *Client) GetLatestTick(market string, interval string) (t TickInterval, e error) {
    response, err := c.get("/pub/market/GetLatestTick",
        map[string]string {"marketName" : market, "tickInterval" : interval}, "2.0")
    var ticks []TickInterval
    err = json.Unmarshal(response.Result, &ticks)
    return ticks[0], err
}

func GetEMA(t []TickInterval, intervals int) (ema float64) {
    for i, item := range t {
        n := intervals
        if n > (i + 1) {
            n = i + 1
        }
        alpha := 2.0 / (1.0 + float64(n))
        ema = ((item.BV / item.V) * alpha) + (ema * (1.0 - alpha))
    }
    return
}

