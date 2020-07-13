package bot
import(
	"crypto_Trading/model"
)
type AutoTradeData struct{
	w model.TransactionService
}

func NewAutoTrade() *AutoTradeData{
	buyOrSellChan := make(chan model.Transaction)
	transControllerChan := make(chan model.Transaction)
	//go d.w.AccountManager(transControllerChan)
	go d.w.MakeProfit(buyOrSellChan, transControllerChan)
			
	return &AutoTradeData{
		w: ws,
	}
}

func  AutoTrade(tdChan chan model.Transaction, tsChan chan model.TransactionService){
	for{
		select{
		case td := <-tdChan:
			
			buyOrSellChan <- td
			td.Side = "sell"
			buyOrSellChan <- td
		case ts := <-tsChan:
		}
	}
			
}