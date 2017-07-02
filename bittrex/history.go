package bittrex

import (
    "encoding/json"
    "time"
)

type HistoryOrder struct {
    Id int64
    TimeStamp HistoryTime
    Quantity float64
    Price float64
    Total float64
    FillType string
    OrderType string
}

type HistoryTime struct {
    time.Time
}

func (t *HistoryTime) UnmarshalJSON(b []byte) error {
    date, err := time.Parse("\"2006-01-02T15:04:05.999\"", string(b))
    *t = HistoryTime{date}
    return err
}

func (c *Client) GetMarketHistory(market string) (h []HistoryOrder, e error) {
    response, e := c.get("/public/getmarkethistory", map[string]string {"market" : market}, "1.1")
    e = json.Unmarshal(response.Result, &h)
    return h, e
}

func GetTickFromDate(h []HistoryOrder, t TickTime) (tick TickInterval) {
    tick.C = h[0].Price
    tick.L = h[0].Price
    tick.T = t
    for i := 0; i < len(h) && t.Time.Before(h[i].TimeStamp.Time); i++ {
        item := h[i]
        tick.O = item.Price
        if tick.H < item.Price {
            tick.H = item.Price
        }
        if tick.L > item.Price {
            tick.L = item.Price
        }
        tick.V += item.Quantity
        tick.BV += item.Total
    }
    return tick
}
