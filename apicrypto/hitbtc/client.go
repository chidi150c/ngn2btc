package hitbtc

import (
	"myhitbtc/model"
)

type Client struct {
	url     string
	SesSion *Session
}

func NewClient(host, publicKey, secret string) *Client {
	c := &Client{
		url:     host + "/api/2",
		SesSion: NewSession(),
	}
	c.SesSion.auth = []string{publicKey, secret}
	return c
}

/*GetSymbol returns the actual list of currency symbols (currency pairs) traded
on HitBTC exchange. The first listed currency is called the base currency,
and the second currency is called the quote currency. The currency pair indicates
how much of the quote currency is needed to purchase one unit of the base currency.
*/
func (c *Client) GetSymbol(symbolCode string) (*model.Symbol, error) {
	//Get symbol
	symbol, err := c.SesSion.getSymbol(c.url+"/public/symbol/", symbolCode)
	return symbol, err
}
func (c *Client) GetCoinMarketData() (model.Coinmarketcaps, error) {
	coinMktData, err := c.SesSion.getCoinMarketData()
	return coinMktData, err
}
func (c *Client) GetTicker(symbolCode string) (*model.Ticker, error) {
	//Get symbol
	symbol, err := c.SesSion.getTicker(c.url+"/public/ticker/", symbolCode)
	return symbol, err
}

func (c *Client) GetSymbols() (model.Symbols, error) {
	//Get symbol
	symbols, err := c.SesSion.getSymbols(c.url + "/public/symbol/")
	return symbols, err
}

func (c *Client) GetOrderBook(symbolCode string) (*model.OrderBook, error) {
	//Get symbol
	orderbook, err := c.SesSion.getOrderBook(c.url+"/public/orderbook/", symbolCode)
	return orderbook, err
}

func (c *Client) GetOrderBooks(symbolCode string) (model.OrderBooks, error) {
	//Get symbol
	orderbooks, err := c.SesSion.getOrderBooks(c.url+"/public/orderbook/", symbolCode)
	return orderbooks, err
}

func (c *Client) GetAddress(symbolCode string) (*model.Address, error) {
	//Get symbol
	address, err := c.SesSion.getAddress(c.url+"/account/crypto/address/", symbolCode)
	return address, err
}

func (c *Client) GetAccountBalances() (model.Balances, error) {
	//Get symbol
	balances, err := c.SesSion.getBalances(c.url + "/account/balance")
	return balances, err
}

func (c *Client) GetTradingBalances() (model.Balances, error) {
	//Get symbol
	balances, err := c.SesSion.getBalances(c.url + "/trading/balance")
	return balances, err
}

func (c *Client) Transfer(currencyCode, amount, toExchange string) (*model.TransferOk, error) {
	transfer, err := c.SesSion.postTransfer(c.url+"/account/transfer", currencyCode, amount, toExchange)
	return transfer, err
}

func (c *Client) NewOrder(clientOrderId, symbolCode, side, quantity, price string) (*model.Order, error) {
	//"""Place an order."""
	data := model.Order{Symbol: symbolCode, Side: side, Quantity: quantity}
	if price != "" {
		data.Price = price
	}
	newOrder, err := c.SesSion.putNewOrder(c.url+"/order/"+clientOrderId, data)
	return newOrder, err
}

func (c *Client) GetOrder(clientOrderId, wait string) (*model.Order, error) {
	//"""Get order info."""
	return c.SesSion.getOrder(c.url + "/order/" + clientOrderId + "?wait=" + wait)
}

func (c *Client) CancelOrder(clientOrderId string) (*model.Order, error) {
	//"""Cancel order."""
	return c.SesSion.delete(c.url + "/order/" + clientOrderId)
}

func (c *Client) Withdraw(currencyCode, amount, address, networkFee string) (*model.Transaction, error) {
	//"""Withdraw."""
	data := model.Transaction{Currency: currencyCode, Amount: amount, Address: address}

	if networkFee != "" {
		data.NetworkFee = networkFee
	}
	return c.SesSion.postWithdraw(c.url+"/account/crypto/withdraw", data)
}

func (c *Client) GetTransaction(transactionID string) (*model.Transaction, error) {
	//"""Get transaction info."""
	return c.SesSion.getTransaction(c.url + "/account/transactions/" + transactionID)
}

func (c *Client) GetCandles(symbolCode, limit, period string) (*model.Candles, error) {
	return c.SesSion.getCandles(c.url + "/public/candles/" + symbolCode + "?limit=" + limit + "&period=" + period)
}

// func (c *Client) GetSymbol(symbolCode string) (interface{}, error) {
// 	//Get symbol
// 	symbol, err := c.SesSion.get(c.url + "public/symbol/" + symbolCode)
// 	return symbol, err
// }
