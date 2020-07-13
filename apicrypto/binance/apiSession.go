package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"time"
	//"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// type Simulate struct {
// 	PlacedOrder          float64
// 	LastPrice            float64
// 	QuoteCurrencyBalance float64
// 	BaseCurrencyBalance  float64
// }
type Session struct {
	auth []string
	// BuyOrderChan         chan Simulate
	// SellOrderChan        chan Simulate
	// SellBalanceCheckChan chan Simulate
	// BuyBalanceCheckChan  chan Simulate
}

func NewSession() *Session {
	return &Session{
		auth: []string{},
		// BuyOrderChan:         make(chan Simulate),
		// SellOrderChan:        make(chan Simulate),
		// SellBalanceCheckChan: make(chan Simulate, 1),
		// BuyBalanceCheckChan:  make(chan Simulate, 1),
	}
}

// func (s *Session) getSymbol(base, endpoint string) (*Symbol, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	url := base + endpoint
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var symbol Symbol
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		Jdata := json.NewDecoder(Resp.Body)
// 		if endpoint == "" {
// 			return nil, errors.New("Get " + url + " Error: symbol empty please use the getSymbols method")
// 		} else {
// 			err := Jdata.Decode(&symbol)
// 			if err != nil {
// 				return nil, errors.Wrap(err, "Get "+url+" Error")
// 			}
// 			return &symbol, nil
// 		}
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) getTicker(base, endpoint string) (*Ticker, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	url := base + endpoint
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var ticker Ticker
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		Jdata := json.NewDecoder(Resp.Body)
// 		if endpoint == "" {
// 			return nil, errors.New("Get " + url + " Error: ticker empty please use the getTickers method")
// 		} else {
// 			err := Jdata.Decode(&ticker)
// 			if err != nil {
// 				return nil, errors.Wrap(err, "Get "+url+" Error")
// 			}
// 			return &ticker, nil
// 		}
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) getSymbols(base, endpoint string) (Symbols, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	url := base + endpoint
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		var symbols Symbols
// 		var symbol Symbol
// 		Jdata := json.NewDecoder(Resp.Body)
// 		if endpoint == "" {
// 			err := Jdata.Decode(&symbols)
// 			if err != nil {
// 				return nil, errors.Wrap(err, "Get "+url+" Error")
// 			}
// 		} else {
// 			err := Jdata.Decode(&symbol)
// 			if err != nil {
// 				return nil, errors.Wrap(err, "Get "+url+" Error")
// 			}
// 			symbols = append(symbols, symbol)
// 		}
// 		return symbols, nil
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) getOrderBook(base, endpoint string) (*OrderBook, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	url := base + endpoint
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var orderBook OrderBook
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		Jdata := json.NewDecoder(Resp.Body)
// 		if endpoint == "" {
// 			return nil, errors.New("Get " + url + " Error: symbol empty please use the OrderBooks method")
// 		} else {
// 			err := Jdata.Decode(&orderBook)
// 			if err != nil {
// 				return nil, errors.Wrap(err, "Get "+url+" Error")
// 			}
// 		}
// 		return &orderBook, nil
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) getOrderBooks(base, endpoint string) (OrderBooks, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	url := base + endpoint
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		var orderBooks OrderBooks
// 		var orderBook OrderBook
// 		Jdata := json.NewDecoder(Resp.Body)
// 		if endpoint == "" {
// 			err := Jdata.Decode(&orderBooks)
// 			if err != nil {
// 				return nil, errors.Wrap(err, "Get "+url+" Error")
// 			}
// 		} else {
// 			err := Jdata.Decode(&orderBook)
// 			if err != nil {
// 				return nil, errors.Wrap(err, "Get "+url+" Error")
// 			}
// 			orderBooks = append(orderBooks, orderBook)
// 		}
// 		return orderBooks, nil
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) getAddress(base, endpoint string) (*Address, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	url := base + endpoint
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	req.SetBasicAuth(s.auth[0], s.auth[1])
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		//var addresss Addresss
// 		var address Address
// 		Jdata := json.NewDecoder(Resp.Body)
// 		err := Jdata.Decode(&address)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Get "+url+" Error")
// 		}
// 		return &address, nil
// 	}
// 	//fmt.Println(Resp.StatusCode)
// 	return nil, errors.New("Get " + url + " Error: " + Resp.Status)
// }

// func (s *Session) getBalances(url string) (Balances, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	req.SetBasicAuth(s.auth[0], s.auth[1])
// 	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		var data Balances
// 		Jdata := json.NewDecoder(Resp.Body)
// 		err := Jdata.Decode(&data)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Get "+url+" Error")
// 		}
// 		return data, nil
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return nil, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) postTransfer(url string, Msg ...string) (*TransferOk, error) {
// 	var MsgData = netUrl.Values{}
// 	MsgData.Set("currency", Msg[0])
// 	MsgData.Set("amount", Msg[1])
// 	MsgData.Set("type", Msg[2])
// 	msgDataReader := *strings.NewReader(MsgData.Encode())
// 	// Create Client
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("POST", url, &msgDataReader)
// 	req.Header.Add("Accept", "application/json")
// 	req.SetBasicAuth(s.auth[0], s.auth[1])
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data TransferOk
// 	var edata AppError
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		Jdata := json.NewDecoder(Resp.Body)
// 		err := Jdata.Decode(&data)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Get "+url+" Error")
// 		}
// 		return &data, nil
// 	}
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return &data, errors.New("Get " + url + " Error: " + Resp.Status + "App Error Message: " + edata.Error.Message + "App Error Description: " + edata.Error.Description)
// }
func (s *Session) exchangeInfo(url string) (*Exchange, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data Exchange
	var edata AppError
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		err := json.NewDecoder(Resp.Body).Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &data, nil
	}
	err = json.NewDecoder(Resp.Body).Decode(&edata)
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + "App Error Message: " + edata.Msg)
}

func (s *Session) orderBookWithInterface(url, symbol, limit string) (*OrderBookWithInterface, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url+"?symbol="+symbol+"&limit="+limit, nil)
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data OrderBookWithInterface
	var edata AppError
	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
		err := json.NewDecoder(Resp.Body).Decode(&data)
		if err != nil {
			return nil, errors.Wrap(err, "Get "+url+" Error")
		}
		return &data, nil
	}
	Edata := json.NewDecoder(Resp.Body)
	err = Edata.Decode(&edata)
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + "App Error Message: " + edata.Msg)
}

func (s *Session) checkServerTime(url string) (*ServerTime, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data ServerTime
	var edata AppError
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
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + "App Error Message: " + edata.Msg)
}

func (s *Session) testNewOrder(host string, dat Order) (*Order, error) {
	tim := strconv.FormatInt(time.Now().Unix()*1000+1000, 10)
	urlPart := "symbol=" + dat.Symbol + "&side=" + dat.Side + "&type=LIMIT&timeInForce=GTC&quantity=" + dat.Quantity + "&price=" + dat.Price + "&recvWindow=5000&timestamp=" + tim
	secret := []byte(s.auth[1])
	message := []byte(urlPart)
	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	Signature := hex.EncodeToString(hash.Sum(nil))
	url := host + "?" + urlPart + "&signature=" + Signature
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-MBX-APIKEY", s.auth[0])
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data Order
	var edata AppError
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
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Msg)
}

func (s *Session) newOrder(host string, dat Order) (*Order, error) {
	tim := strconv.FormatInt(time.Now().Unix()*1000+1000, 10)
	urlPart := "symbol=" + dat.Symbol + "&side=" + dat.Side + "&type=LIMIT&timeInForce=GTC&quantity=" + dat.Quantity + "&price=" + dat.Price + "&newClientOrderId=" + dat.NewClientOrderId + "&recvWindow=5000&timestamp=" + tim
	secret := []byte(s.auth[1])
	message := []byte(urlPart)
	hash := hmac.New(sha256.New, secret)
	hash.Write(message)
	Signature := hex.EncodeToString(hash.Sum(nil))
	url := host + "?" + urlPart + "&signature=" + Signature
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-MBX-APIKEY", s.auth[0])
	Resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var data Order
	var edata AppError
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
	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Msg)
}

// func (s *Session) postWithdraw(url string, dat Transaction) (*Transaction, error) {
// 	var MsgData = netUrl.Values{}
// 	MsgData.Set("currency", dat.Currency)
// 	MsgData.Set("amount", dat.Amount)
// 	MsgData.Set("address", dat.Address)
// 	MsgData.Set("networkFee", dat.NetworkFee)
// 	msgDataReader := *strings.NewReader(MsgData.Encode())
// 	// Create Client
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("PUT", url, &msgDataReader)
// 	req.Header.Add("Accept", "application/json")
// 	req.SetBasicAuth(s.auth[0], s.auth[1])
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data Transaction
// 	var edata AppError
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		Jdata := json.NewDecoder(Resp.Body)
// 		err := Jdata.Decode(&data)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Get "+url+" Error")
// 		}
// 		return &data, nil
// 	}
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) getTransaction(url string) (*Transaction, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	req.SetBasicAuth(s.auth[0], s.auth[1])
// 	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data Transaction
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		Jdata := json.NewDecoder(Resp.Body)
// 		err := Jdata.Decode(&data)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Get "+url+" Error")
// 		}
// 		return &data, nil
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }

// func (s *Session) getCandles(url string) (*Candles, error) {
// 	// Create Client
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("Accept", "application/json")
// 	//req.SetBasicAuth(s.auth[0], s.auth[1])
// 	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	Resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data Candles
// 	if Resp.StatusCode >= 200 && Resp.StatusCode < 300 {
// 		Jdata := json.NewDecoder(Resp.Body)
// 		err := Jdata.Decode(&data)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Get "+url+" Error")
// 		}
// 		return &data, nil
// 	}
// 	var edata AppError
// 	Edata := json.NewDecoder(Resp.Body)
// 	err = Edata.Decode(&edata)
// 	//fmt.Println(Resp.StatusCode)
// 	return &data, errors.New("Get " + url + " Error: " + Resp.Status + " App Error Message: " + edata.Error.Message + " App Error Description: " + edata.Error.Description)

// }
