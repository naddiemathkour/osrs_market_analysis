package models

type ItemData struct {
	AvgHighPrice    int `json:"avgHighPrice"`
	AvgLowPrice     int `json:"avgLowPrice"`
	HighPriceVolume int `json:"highPriceVolume"`
	LowPriceVolume  int `json:"lowPriceVolume"`
}

type ItemObject struct {
	ID   string   `json:"id"`
	Data ItemData `json:"data"`
}
