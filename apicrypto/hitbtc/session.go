package hitbtc

import (
	"encoding/json"
	"log"
	"myhitbtc/model"
	"net/http"
	netUrl "net/url"
	"strings"

	"github.com/pkg/errors"
)

type Simulate struct {
	PlacedOrder          float64
	LastPrice            float64
	QuoteCurrencyBalance float64
	BaseCurrencyBalance  float64
}
type Session struct {
	auth                 []string
	BuyOrderChan         chan Simulate
	SellOrderChan        chan Simulate
	SellBalanceCheckChan chan Simulate
	BuyBalanceCheckChan  chan Simulate
}

func NewSession() *Session {
	return &Session{
		auth:                 []string{},
		BuyOrderChan:         make(chan Simulate),
		SellOrderChan:        make(chan Simulate),
		SellBalanceCheckChan: make(chan Simulate, 1),
		BuyBalanceCheckChan:  make(chan Simulate, 1),
	}
}

func (s *Session) getSymbol(base, endpoint string) (*model.Symbol, error) {
	// Create Client
	client := &http.Client{}
	url := base + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var symbol model.Symbol
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		if endpoint == "" {
			return nil, errors.New("Get " + url + " Error: symbol empty please use the getSymbols method")
		} else {
			err := Jdata.Decode(&symbol)
			if err != nil {
				return nil, errors.Wrap(err, "Get "+url+" Error")
			}
			return &symbol, nil
		}
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}

func (s *Session) getTicker(base, endpoint string) (*model.Ticker, error) {
	// Create Client
	client := &http.Client{}
	url := base + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var ticker model.Ticker
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		if endpoint == "" {
			return nil, errors.New("Get " + url + " Error: ticker empty please use the getTickers method")
		} else {
			err := Jdata.Decode(&ticker)
			if err != nil {
				return nil, errors.Wrap(err, "Get "+url+" Error")
			}
			return &ticker, nil
		}
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}

func (s *Session) getSymbols(base string) (model.Symbols, error) {
	// Create Client
	client := &http.Client{}
	url := base
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		var symbols model.Symbols
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&symbols)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return symbols, nil
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)
}

func (s *Session) getOrderBook(base, endpoint string) (*model.OrderBook, error) {
	// Create Client
	client := &http.Client{}
	url := base + endpoint + "?limit=1"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var orderBook model.OrderBook
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		if endpoint == "" {
			return nil, errors.New("Get " + url + " Error: symbol empty please use the OrderBooks method")
		} else {
			err := Jdata.Decode(&orderBook)
			if err != nil {
				return nil, errors.Wrap(err, "Get "+url+" Error")
			}
		}
		return &orderBook, nil
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}

func (s *Session) getOrderBooks(base, endpoint string) (model.OrderBooks, error) {
	// Create Client
	client := &http.Client{}
	url := base + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		var orderBooks model.OrderBooks
		var orderBook model.OrderBook
		Jdata := json.NewDecoder(Resp.Body)
		if endpoint == "" {
			err := Jdata.Decode(&orderBooks)
			if err != nil {
				return nil, errors.Wrap(err, "Get "+url+" Error")
			}
		} else {
			err := Jdata.Decode(&orderBook)
			if err != nil {
				return nil, errors.Wrap(err, "Get "+url+" Error")
			}
			orderBooks = append(orderBooks, orderBook)
		}
		return orderBooks, nil
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}

func (s *Session) getAddress(base, endpoint string) (*model.Address, error) {
	// Create Client
	client := &http.Client{}
	url := base + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		//var addresss model.Addresss
		var address model.Address
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&address)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &address, nil
	}
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status)
}

func (s *Session) getBalances(url string) (model.Balances, error) {
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		var data model.Balances
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return data, nil
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}

func (s *Session) getCoinMarketData() (model.Coinmarketcaps, error) {
	bitResp, err := http.Get("https://api.coinmarketcap.com/v1/ticker")
	if err != nil {
		return nil, err
	}
	bitdecoder := json.NewDecoder(bitResp.Body)
	var data model.Coinmarketcaps
	err = bitdecoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *Session) postTransfer(url string, Msg ...string) (*model.TransferOk, error) {
	var MsgData = netUrl.Values{}
	MsgData.Set("currency", Msg[0])
	MsgData.Set("amount", Msg[1])
	MsgData.Set("type", Msg[2])
	msgDataReader := *strings.NewReader(MsgData.Encode())
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, &msgDataReader)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data model.TransferOk
	var edata model.AppError
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &data, nil
	}
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + "App Error Message: " + edata.Error.Message + "App Error Description: " + edata.Error.Description)
}

func (s *Session) putNewOrder(url string, dat model.Order) (*model.Order, error) {
	var MsgData = netUrl.Values{}
	MsgData.Set("symbol", dat.Symbol)
	MsgData.Set("side", dat.Side)
	MsgData.Set("timeInForce", "GTC")
	MsgData.Set("quantity", dat.Quantity)
	MsgData.Set("price", dat.Price)
	msgDataReader := *strings.NewReader(MsgData.Encode())
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url, &msgDataReader)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data model.Order
	var edata model.AppError
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &data, nil
	}
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)
}

func (s *Session) getOrder(url string) (*model.Order, error) {
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var order3 interface{}
	var order2 model.RespAsMap
	order2 = make(model.RespAsMap, 1)
	order2[0] = make(map[string]string)
	var order model.Order
	var edata model.AppError
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&order3)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		switch order3.(type) {
		case model.RespAsMap:
			order2 = order3.(model.RespAsMap)
			order = model.Order{
				Id:            order2[0]["id"],
				ClientOrderId: order2[0]["clientOrderId"],
				Symbol:        order2[0]["symbol"],
				Side:          order2[0]["side"],
				Status:        order2[0]["status"],
				Type:          order2[0]["type"],
				TimeInForce:   order2[0]["timeInForce"],
				Quantity:      order2[0]["quantity"],
				Price:         order2[0]["price"],
				CumQuantity:   order2[0]["cumQuantity"],
				CreatedAt:     order2[0]["createdAt"],
				UpdatedAt:     order2[0]["updatedAt"],
			}
		case model.Order:
			order = order3.(model.Order)
		}
		return &order, nil
	}
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)
}

func (s *Session) delete(url string) (*model.Order, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var order3 interface{}
	var order2 model.RespAsMap
	order2 = make(model.RespAsMap, 1)
	order2[0] = make(map[string]string)
	var order model.Order
	var edata model.AppError
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&order3)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		log.Println(order3)
		switch order3.(type) {
		case model.RespAsMap:
			order2 = order3.(model.RespAsMap)
			order = model.Order{
				Id:            order2[0]["id"],
				ClientOrderId: order2[0]["clientOrderId"],
				Symbol:        order2[0]["symbol"],
				Side:          order2[0]["side"],
				Status:        order2[0]["status"],
				Type:          order2[0]["type"],
				TimeInForce:   order2[0]["timeInForce"],
				Quantity:      order2[0]["quantity"],
				Price:         order2[0]["price"],
				CumQuantity:   order2[0]["cumQuantity"],
				CreatedAt:     order2[0]["createdAt"],
				UpdatedAt:     order2[0]["updatedAt"],
			}
		case model.Order:
			order = order3.(model.Order)
		}

		return &order, nil
	}
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)
}

func (s *Session) postWithdraw(url string, dat model.Transaction) (*model.Transaction, error) {
	var MsgData = netUrl.Values{}
	MsgData.Set("currency", dat.Currency)
	MsgData.Set("amount", dat.Amount)
	MsgData.Set("address", dat.Address)
	MsgData.Set("networkFee", dat.NetworkFee)
	msgDataReader := *strings.NewReader(MsgData.Encode())
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url, &msgDataReader)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data model.Transaction
	var edata model.AppError
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &data, nil
	}
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}

func (s *Session) getTransaction(url string) (*model.Transaction, error) {
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(s.auth[0], s.auth[1])
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data model.Transaction
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &data, nil
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}

func (s *Session) getCandles(url string) (*model.Candles, error) {
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	//req.SetBasicAuth(s.auth[0], s.auth[1])
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data model.Candles
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		Jdata := json.NewDecoder(Resp.Body)
		err := Jdata.Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &data, nil
	}
	var edata model.AppError
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	//fmt.Println(Resp.StatusCode)
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

}
