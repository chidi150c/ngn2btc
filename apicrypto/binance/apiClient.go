package binance

import (
	"strconv"
)

type Client struct {
	url     string
	SesSion *Session
}

func NewClient(host, publicKey, secret string) *Client {
	c := &Client{
		url:     host,
		SesSion: NewSession(),
	}
	c.SesSion.auth = []string{publicKey, secret}
	return c
}

// /*GetSymbol returns the actual list of currency symbols (currency pairs) traded
// on HitBTC exchange. The first listed currency is called the base currency,
// and the second currency is called the quote currency. The currency pair indicates
// how much of the quote currency is needed to purchase one unit of the base currency.
// */
// func (c *Client) GetSymbol(symbolCode string) (*Symbol, error) {
// 	//Get symbol
// 	symbol, err := c.SesSion.getSymbol(c.url+"/public/symbol/", symbolCode)
// 	return symbol, err
// }

// func (c *Client) GetTicker(symbolCode string) (*Ticker, error) {
// 	//Get symbol
// 	symbol, err := c.SesSion.getTicker(c.url+"/public/ticker/", symbolCode)
// 	return symbol, err
// }

// func (c *Client) GetSymbols(symbolCode string) (Symbols, error) {
// 	//Get symbol
// 	symbols, err := c.SesSion.getSymbols(c.url+"/public/symbol/", symbolCode)
// 	return symbols, err
// }

// func (c *Client) GetOrderBook(symbolCode string) (*OrderBook, error) {
// 	//Get symbol
// 	orderbook, err := c.SesSion.getOrderBook(c.url+"/public/orderbook/", symbolCode)
// 	return orderbook, err
// }

// func (c *Client) GetOrderBooks(symbolCode string) (OrderBooks, error) {
// 	//Get symbol
// 	orderbooks, err := c.SesSion.getOrderBooks(c.url+"/public/orderbook/", symbolCode)
// 	return orderbooks, err
// }

// func (c *Client) GetAddress(symbolCode string) (*Address, error) {
// 	//Get symbol
// 	address, err := c.SesSion.getAddress(c.url+"/account/crypto/address/", symbolCode)
// 	return address, err
// }

// func (c *Client) GetAccountBalances() (Balances, error) {
// 	//Get symbol
// 	balances, err := c.SesSion.getBalances(c.url + "/account/balance")
// 	return balances, err
// }

// func (c *Client) GetTradingBalances() (Balances, error) {
// 	//Get symbol
// 	balances, err := c.SesSion.getBalances(c.url + "/trading/balance")
// 	return balances, err
// }

// func (c *Client) Transfer(currencyCode, amount, toExchange string) (*TransferOk, error) {
// 	transfer, err := c.SesSion.postTransfer(c.url+"/account/transfer", currencyCode, amount, toExchange)
// 	return transfer, err
// }

func (c *Client) OrderBook(symbol, limit string) (*OrderBook, error) {
	//"""Place an order."""
	limitAsInt, _ := strconv.ParseInt(limit, 10, 64)
	orderBookWithInterface, err := c.SesSion.orderBookWithInterface(c.url+"/api/v1/depth", symbol, limit)
	orderBook := OrderBook{
		LastUpdateId: orderBookWithInterface.LastUpdateId,
		Bids:         make([]Label, limitAsInt),
		Asks:         make([]Label, limitAsInt),
	}
	for i := 0; i < int(limitAsInt); i++ {
		orderBook.Bids[i].Price = orderBookWithInterface.Bids[i][0].(string)
		orderBook.Bids[i].Quantity = orderBookWithInterface.Bids[i][1].(string)
		orderBook.Asks[i].Price = orderBookWithInterface.Asks[i][0].(string)
		orderBook.Asks[i].Quantity = orderBookWithInterface.Asks[i][1].(string)
	}
	//_ = Label{Price: orderBook.Asks[0][0].(string), Quantity: orderBook.Bids[0][1].(string)}
	return &orderBook, err
}

func (c *Client) ExchangeInfo() (*Exchange, error) {
	//"""Place an order."""
	exchangeInf, err := c.SesSion.exchangeInfo(c.url + "/api/v1/exchangeInfo")
	return exchangeInf, err
}

func (c *Client) CheckServerTime() (*ServerTime, error) {
	//"""Place an order."""
	serverTime, err := c.SesSion.checkServerTime(c.url + "/api/v1/time")
	return serverTime, err
}

func (c *Client) TestNewOrder(symbolCode, side, quantity, price string) (*Order, error) {
	//"""Place an order."""
	data := Order{Symbol: symbolCode, Side: side, Quantity: quantity, Price: price}
	testNewOrder, err := c.SesSion.testNewOrder(c.url+"/api/v3/order/test", data)
	return testNewOrder, err
}

func (c *Client) NewOrder(clientOrderId, symbolCode, side, quantity, price string) (*Order, error) {
	//"""Place an order."""
	data := Order{NewClientOrderId: clientOrderId, Symbol: symbolCode, Side: side, Quantity: quantity, Price: price}
	newOrder, err := c.SesSion.newOrder(c.url+"/api/v3/order", data)
	return newOrder, err
}

// func (c *Client) GetOrder(clientOrderId, wait string) (*Order, error) {
// 	//"""Get order info."""
// 	return c.SesSion.getOrder(c.url + "/api/v3/order " + clientOrderId + "?wait=" + wait)
// }

// func (c *Client) CancelOrder(clientOrderId string) (*Order, error) {
// 	//"""Cancel order."""
// 	return c.SesSion.delete(c.url + "/order/" + clientOrderId)
// }

// func (c *Client) Withdraw(currencyCode, amount, address, networkFee string) (*Transaction, error) {
// 	//"""Withdraw."""
// 	data := Transaction{Currency: currencyCode, Amount: amount, Address: address}

// 	if networkFee != "" {
// 		data.NetworkFee = networkFee
// 	}
// 	return c.SesSion.postWithdraw(c.url+"/account/crypto/withdraw", data)
// }

// func (c *Client) GetTransaction(transactionID string) (*Transaction, error) {
// 	//"""Get transaction info."""
// 	return c.SesSion.getTransaction(c.url + "/account/transactions/" + transactionID)
// }

// func (c *Client) GetCandles(symbolCode, limit, period string) (*Candles, error) {
// 	return c.SesSion.getCandles(c.url + "/public/candles/" + symbolCode + "?limit=" + limit + "&period=" + period)
// }

// // func (c *Client) GetSymbol(symbolCode string) (interface{}, error) {
// // 	//Get symbol
// // 	symbol, err := c.SesSion.get(c.url + "public/symbol/" + symbolCode)
// // 	return symbol, err
// // }
