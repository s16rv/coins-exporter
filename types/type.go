package types

import "reflect"

type CurrentPrice struct {
	BTC float64 `json:"btc"`
	ETH float64 `json:"eth"`
	IDR float64 `json:"idr"`
	USD float64 `json:"usd"`
}

type ReturnData struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	MarketData struct {
		CurrentPrice CurrentPrice `json:"current_price"`
	} `json:"market_data"`
}

func (cp *CurrentPrice) GetField(field string) float64 {
	r := reflect.ValueOf(cp)
	f := reflect.Indirect(r).FieldByName(field)
	return float64(f.Float())
}
