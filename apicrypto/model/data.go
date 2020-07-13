package model

type TransactionID int64

type Transaction struct {
	ID TransactionID
	SymbolID string
	SoldOrBoughtPrice float64
	ProfitTarget float64
	ProfitTargetRate float64
	AdvancingRate float64
	QtyMultiplier float64
	Side string
	TransType string
	OrderQty float64
	PreviousOrderID string
	BuyOrSellChan chan Transaction
	StopMeChan chan string
	Count float64
	UserApiKey string
	UserApiSecret string
	Host string
}

type TransactionService interface{
	//AccountManager(chan Transaction, chan float64, chan float64)
	MakeProfit(chan Transaction, chan Transaction, chan float64, chan float64)
	AddTransaction(*Transaction) (TransactionID, error)
	GetTransaction(TransactionID) (*Transaction, error)
	DeleteTransaction(TransactionID) error
	ListTransactions() ([]*Transaction, error)
	UpdateTransaction(*Transaction) error
}

