package apicrypto

type TransactionID int64

type TransDBType map[TransactionID]*Transaction

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
	userApiKey string
	userApiSecret string
	Host string
	userToken string
}

type TransactionService interface{
	//AccountManager(chan Transaction, chan float64, chan float64)
	MakeProfit(chan Transaction, chan Transaction, chan float64, chan float64)
	SetSession(*Session)
	AddTransaction(*User, *Transaction) (TransactionID, error)
	GetTransaction(*User, TransactionID) (*Transaction, error)
	DeleteTransaction(*User, TransactionID) error
	ListTransactions() ([]*Transaction, error)
	UpdateTransaction(*User, *Transaction) error
}

