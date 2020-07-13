package apicrypto
 
import(
	"time"
	"strconv"
	"fmt"
	"log"
	"user-apiv2/apicrypto/binance"
)

type Worker struct{
	session *Session
	Symbol binance.Symbol
	API binance.Client
	minAmount float64
	minQty float64
	minPrice float64
}

func NewWorker (smID string, cl binance.Client) *Worker{
	s := getSymbol(smID, cl)
	w := &Worker{
		API: cl,
		Symbol: s,
	}
	w.minAmount, _ = strconv.ParseFloat(w.Symbol.Filters[2].MinNotional, 64)
	w.minQty, _ = strconv.ParseFloat(w.Symbol.Filters[1].MinQty, 64)
	w.minPrice, _ = strconv.ParseFloat(w.Symbol.Filters[0].MinPrice, 64)
	return w
}

var _ TransactionService = &Worker{}

func (w *Worker)MakeProfit(buyOrSellChan, transControllerChan chan Transaction, countGen, countRecover chan float64){
	for{
		select{
		case td := <-buyOrSellChan:
			td.Count = <-countGen
			go func (td Transaction){
				var bestPrice, costOrSellPrice, quoteBal float64 
				//Manager init Sell or buy
				if td.SoldOrBoughtPrice == 0.0 {
					bestPrice = w.API.getBestPrice()
					td.SoldOrBoughtPrice = bestPrice
					td.StopMeChan = make(chan string, 2)
				}	
				if td.TransType == "A" {	
					//Advancing buy or sell starts
					if td.Side == "buy" {
						td.ProfitTarget = td.ProfitTarget * (1-td.AdvancingRate)			
					}else if td.Side == "sell" {
						td.ProfitTarget = td.ProfitTarget * (1+td.AdvancingRate)
					}else{
						panic("td.Side not \"sell\" or \"buy\": line 45")
					}	
					//Advancing buy or sell ends			
				}else if td.TransType == "B" {
					//sold Or Bought Now buy Or Sell resp. starts
					if td.Side == "buy" {
						td.ProfitTarget = td.SoldOrBoughtPrice * (1-td.ProfitTargetRate)					
					}else if td.Side == "sell" {
						td.ProfitTarget = td.SoldOrBoughtPrice * (1+td.ProfitTargetRate)
					}else{
						panic("td.Side not \"sell\" or \"buy\": line 55")
					}
					//sold Or Bought Now buy Or Sell resp. ends
				}else{
					panic("td.TransType not \"A\" or \"B\": line 64")
				}
								
				for {
					select{
					case stopMsg := <- td.StopMeChan:
						countRecover <- td.Count
						reasonForEnding := fmt.Sprintf("Stopped via td.stopMeChan: %s", stopMsg)
						log.Printf("makeProfit: For %s: type: \"%s\": in %s operation: Goroutine No. %.2f: Ended... |reasonForEnding = %s |sOrbOrderID = %s |\n", w.Symbol.ID, td.tYPE, td.operation, td.count, reasonForEnding, td.soldOrBoughtOrderID)
						return
					}
					bestPrice = w.API.getBestPrice()
					if td.Side == "buy" {
						if bestPrice > td.ProfitTarget{
							continue
						}
					}else if td.Side == "sell" {
						if bestPrice < td.ProfitTarget{
							continue
						}
					}
					orderID := placeOrder(td.orderQty, bestPrice, w.Symbol.SymbolID, td.Side)
					if orderID != "insufficient fund" && orderID != "" && orderID != "below MinAmount" {
						if td.transType == "B" {
							countRecover <- td.Count
							return
						}else if td.transType == "A"{
							td.orderID = orderID
							td.transType = "B"
							td.SoldOrBoughtPrice = bestPrice
							buyOrSellChan <- td
							countRecover <- td.Count
							return
						}
					}
					if orderID == "insufficient fund" {
						log.Printf("Order result is \"%s\"", orderID)
						time.Sleep(time.Minute * 5)
						continue
					}
					if orderID == "below MinAmount" {
						log.Printf("Order result is \"%s\"", orderID)
						time.Sleep(time.Minute * 5)
						continue
					}
				}
			}(td)
		}
	}
}

// func (w *Worker)AccountManager(chan Transaction){
 // 	for{
 // 		select {
 // 		case buyOrSell = <-emaStatusPriceSell:
 // 			//log.Printf("accountManager: for %s: for %s controller: Goroutine No. %.2f: Got \"%s\" from emaStatusPriceSell\n", w.Symbol.ID, bd.operation, bd.count, buyOrSell.status)
 // 			if bd.operation == "buy" {
 // 				//log.Printf("accountManager: for %s: for %s controller: going to kill a %s operator of Gourotine No. %.2f\n", w.Symbol.ID, bd.operation, bd.operation, bd.count)
 // 				stopReason = "trading moved to \"sell\" direction so I was killed to maintain one thread"
 // 				killOperator = true
 // 			} else if bd.operation == "sell" {
 // 				if buyOrSell.lastPrice > bd.highestPrice+tickSize*2.0 {
 // 					//log.Printf("accountManager: for %s: for %s controller: Re-adjusting sell operator to sell immediately as sell operation matached sell status and buyOrSell.lastPrice %.8f > bd.highestPrice %.8f \n", w.Symbol.ID, bd.operation, buyOrSell.lastPrice, bd.highestPrice)
 // 					initData.boughtOrSoldPrice = buyOrSell.lastPrice
 // 					initData.updateEmaDiffChan = buyOrSell.updateChan
 // 					initData.influencedBy = "price"
 // 					stopReason = fmt.Sprintf("controller got a better \"sell\" point \"%.8f\" as informed by checkPrice where lastPrice was higher than highestPrice", buyOrSell.lastPrice)
 // 					amsc.realBuyOrSellSChan <- initData
 // 				} else {
 // 					//log.Printf("accountManager: for %s: for %s controller: did not re-adjust sell operator for immediate selling because lastPrice %.8f is less than highestPrice %.8f + tickSize %.8f * 2.0 \n", w.Symbol.ID, bd.operation, buyOrSell.lastPrice, bd.highestPrice, tickSize)
 // 					continue
 // 				}
 // 			}
 // 		case buyOrSell = <-emaStatusPriceBuy:
 // 			//log.Printf("accountManager: for %s: for %s controller: Goroutine No. %.2f: Got \"%s\" from emaStatusPriceBuy\n", w.Symbol.ID, bd.operation, bd.count, buyOrSell.status)
 // 			if bd.operation == "buy" {
 // 				if buyOrSell.lastPrice < bd.lowestPrice-tickSize*2.0 {
 // 					//log.Printf("accountManager: for %s: for %s controller: Re-adjusting buy operator to buy immediately as buy operation matached buy status and buyOrSell.lastPrice %.8f < bd.lowestPrice %.8f \n", w.Symbol.ID, bd.operation, buyOrSell.lastPrice, bd.lowestPrice)
 // 					initData.boughtOrSoldPrice = buyOrSell.lastPrice
 // 					initData.updateEmaDiffChan = buyOrSell.updateChan
 // 					initData.influencedBy = "price"
 // 					stopReason = fmt.Sprintf("controller got a better \"buy\" point \"%.8f\" as informed by checkPrice where lastPrice is lower than the LowestPrice", buyOrSell.lastPrice)
 // 					amsc.realBuyOrSellSChan <- initData
 // 				} else {
 // 					//log.Printf("accountManager: for %s: for %s controller: did not re-adjust buy operator for immediate buying because lastPrice %.8f is higher than lowestPrice %.8f \n", w.Symbol.ID, bd.operation, buyOrSell.lastPrice, bd.lowestPrice)
 // 					continue
 // 				}
 // 			} else if bd.operation == "sell" {
 // 				//log.Printf("accountManager: for %s: for %s controller: going to kill a %s operator of Gourotine No. %.2f\n", w.Symbol.ID, bd.operation, bd.operation, bd.count)
 // 				stopReason = "trading moved to \"buy\" direction so I was killed to maintain one thread"
 // 				killOperator = true
 // 			}
 // 		case bd.highestPrice = <-bd.updateHighestPriceChan:
 // 			//log.Printf("accountManager: for %s: for %s controller: Goroutine No. %.2f: Got bd.highestPrice updated to %.8f via bd.updateHighestPriceChan\n", w.Symbol.ID, bd.operation, bd.count, bd.highestPrice)
 // 			continue
 // 		case bd.lowestPrice = <-bd.updateLowestPriceChan:
 // 			//log.Printf("accountManager: for %s: for %s controller: Goroutine No. %.2f: Got bd.lowestPrice updated to %.8f via bd.updateLowestPriceChan\n", w.Symbol.ID, bd.operation, bd.count, bd.lowestPrice)
 // 			continue
 // 		}
 // 	}
// }
func (w *Worker)SetSession(sess *Session){
	w.session = sess
}
func (w *Worker)AddTransaction(caller *User, td *Transaction) (TransactionID, error){
	if caller == nil {
		return 0, model.ErrUserRequired
	}
	if caller.Token != td.UserToken{
		return 0, model.ErrUnauthorized
	}
	w.session.addTdDBChan <- td
	return td.ID, nil
}
func (w *Worker)GetTransaction(caller *User, id TransactionID) (*Transaction, error){
	if caller == nil {
		return 0, model.ErrUserRequired
	}
	takeTdChan := make(chan Transaction)
	w.session.getTdDBChan <- getDBData{id, takeTdChan}
	td := <-takeTdChan
	if td == nil {
		return 0, model.ErrTransactionNotFound
	}
	if caller.Token != td.UserToken{
		return 0, model.ErrUnauthorized
	}
	return td, nil
}
func (w *Worker)DeleteTransaction(caller *User, id TransactionID) error{
	if caller == nil {
		return model.ErrUserRequired
	}
	takeTdChan := make(chan Transaction)
	w.session.getTdDBChan <- getDBData{id, takeTdChan}
	td := <-takeTdChan
	if td == nil {
		return model.ErrTransactionNotFound
	}
	if caller.Token != td.UserToken{
		return model.ErrUnauthorized
	}
	w.session.deleteTdDBChan <- id
	return nil
}
func (w *Worker)ListTransactions() ([]*Transaction, error){
return nil, nil
}
func (w *Worker)UpdateTransaction(caller *User, nTd *Transaction) error{
	if caller == nil {
		return model.ErrUserRequired
	}
	takeTdChan := make(chan Transaction)
	w.session.getTdDBChan <- getDBData{nTd.ID, takeTdChan}
	td := <-takeTdChan
	if td == nil {
		return model.ErrTransactionNotFound
	}
	if caller.Token != td.UserToken{
		return model.ErrUnauthorized
	}
	w.session.addTdDBChan <-nTd
	return nil
}

func  getSymbol(smID string, cl binance.Client) binance.Symbol{
	var sym Symbol

	info, _ := cl.ExchangeInfo()
	in := *info
	//fmt.Println(in)

	for _, symbol := range in.Symbols {
		if symbol.SymbolID == smID {
			sym = symbol
		}
	}
	return sym
}
		