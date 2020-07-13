package apiexch

import "context"

// General
const apiCode = "fecf93e1-56a0-4eef-90de-a388f84c7e05"

type ExchID int64

type ExDBType map[ExchID]*Exch

type Exch struct {
	ID                   ExchID
	Sellers              []User
	Buyers               []User
	Markets              []Market
	Status               string
	Message              string
	Name                 string
	MinDeposit           string
	MinWithdrawal        string
	DepositNetworkFee    string
	WithdrawalNetworkFee string
	OwnerName            string
	OwnerToken           string
}

type ExchServicer interface {
	BuyBitcoin(context.Context)
	SellBitcoin(context.Context)
	WithdrawBitcoin(context.Context)
	DepositBitcoin(context.Context)
	AddExch(context.Context, *Exch) (ExchID, error)
	GetExch(context.Context, ExchID) (*Exch, error)
	DeleteExch(context.Context, ExchID) error
	UpdateExch(context.Context, *Exch) error
}

type WalletAddressAPI struct {
	GUID    string `json:"guid"`
	Address string `json:"address"`
}

type Market struct {
	Bid          string `json:"bid"`
	Lastprice    string `json:"last_price"`
	Volume24h    string `json:"volume24h"`
	Currency     string `json:"currency"`
	Marketname   string `json:"marketname"`
	Ask          string `json:ask"`
	Low24h       string `json:"Low24h"`
	Change24h    string `json:"change24h"`
	High24h      string `json:"high24h"`
	Basecurrency string `json:"basecurrency"`
}
