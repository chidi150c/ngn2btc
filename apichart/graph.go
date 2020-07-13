package apichart

//ADBType is the chat user database
//type ADBType map[int]*GraphPoint

// type GraphPointV2 struct {
// 	Time string
// 	Data float64
// }

type CGraphPoint struct {
	Time string
	Data []Coinmarketcap
}

type GraphPoint struct {
	Time string
	Data BlockchainInfo
}

// type GraphWithTimeID struct {
// 	ItsTime int64
// 	Graph   string
// }

// type Coinmarketcaps []Coinmarketcap
type Coinmarketcap struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Volume24hUsd     string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}
type Coin struct {
	M15    float64 `json:"15m"`
	Last   float64 `json:"last"`
	Buy    float64 `json:"buy"`
	Sell   float64 `json:"sell"`
	Symbol string  `json:"symbol"`
}

type BlockchainInfo struct {
	USD Coin
	AUD Coin
	BRL Coin
	CAD Coin
	CHF Coin
	CLP Coin
	CNY Coin
	DKK Coin
	EUR Coin
	GBP Coin
	HKD Coin
	INR Coin
	ISK Coin
	JPY Coin
	KRW Coin
	NZD Coin
	PLN Coin
	RUB Coin
	SEK Coin
	SGD Coin
	THB Coin
	TWD Coin
}
type Currencylayer struct {
	Success   bool   `json:"success"`
	Terms     string `json:"terms"`
	Privacy   string `json:"privacy"`
	Timestamp int    `json:"timestamp"`
	Source    string `json:"source"`
	Quotes    Quotes
}
type Quotes struct {
	USDNGN float64 `json:"USDNGN"`
}
