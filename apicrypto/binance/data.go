package binance

type Rate struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	Limit         int64  `json:"limit"`
}

type Filter struct {
	FilterType  string `json:"filterType"`
	MinPrice    string `json:"minPrice"`
	MaxPrice    string `json:"maxPrice"`
	TickSize    string `json:"tickSize"`
	MinQty      string `json:"minQty"`
	MaxQty      string `json:"maxQty"`
	StepSize    string `json:"stepSize"`
	MinNotional string `json:"minNotional"`
}

type Symbol struct {
	SymbolID           string   `json:"symbol"`
	Status             string   `json:"status"`
	BaseAsset          string   `json:"baseAsset"`
	BaseAssetPrecision int64    `json:"baseAssetPrecision"`
	QuoteAsset         string   `json:"quoteAsset"`
	QuotePrecision     int64    `json:"quotePrecision"`
	OrderTypes         []string `json:"orderTypes"`
	IcebergAllowed     bool     `json:"icebergAllowed"`
	Filters            []Filter `json:"filters"`
}

type Exchange struct {
	Timezone        string   `json:"timezone"`
	ServerTime      int64    `json:"serverTime"`
	RateLimits      []Rate   `json:"rateLimits"`
	ExchangeFilters []string `json:"exchangeFilters"`
	Symbols         []Symbol `json:"symbols"`
}
type ServerTime struct {
	Servertime int64 `json:"serverTime"`
}

type Order struct {
	Symbol           string `json:"symbol"`
	Side             string `json:"side"`
	Type             string `json:"type"`
	TimeInForce      string `json:"timeInForce"`
	Quantity         string `json:"quantity"`
	Price            string `json:"price"`
	RecvWindow       string `json:"recvWindow"`
	Timestamp        string `json:"timestamp"`
	NewClientOrderId string `json:"newClientOrderId"`
	StopPrice        string `json:"stopPrice"`
	IcebergQty       string `json:"icebergQty"`
	NewOrderRespType string `json:"newOrderRespType"`
}

type PriceCommission struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}

type OrderACK struct {
	Symbol        string            `json:"symbol"`
	OrderId       string            `json:"orderId"`
	ClientOrderId string            `json:"clientOrderId"`
	TransactTime  string            `json:"transactTime"`
	Price         string            `json:"price"`
	OrigQty       string            `json:"origQty"`
	ExecutedQty   string            `json:"executedQty"`
	Status        string            `json:"status"`
	TimeInForce   string            `json:"timeInForce"`
	Type          string            `json:"type"`
	Side          string            `json:"side"`
	Fills         []PriceCommission `json:"fills"`
}

type AppError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

type OrderBook struct {
	LastUpdateId int64   `json:"lastUpdateId"`
	Bids         []Label `json:"bids"`
	Asks         []Label `json:"asks"`
}

type OrderBookWithInterface struct {
	LastUpdateId int64           `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}
type Label struct {
	Price    string
	Quantity string
}

type Trades []Trade
type Trade struct {
	ID           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}
