package hitbtc

import (
	"fmt"
	"log"
	// "math"
	"myhitbtc/model"
	"myhitbtc/twilio"
	// "strconv"
	// "strings"
	"time"

	// "github.com/apourchet/investment/lib/ema"
	// hatchuid "github.com/nu7hatch/gouuid"
	// "github.com/pkg/errors"
	// "github.com/rs/xid"
)

type Worker struct {
	//uUIDGetChan   UUIDChan
	twilio        *twilio.TwilioAPI
	API           *Client
	Symbol        *model.Symbol
	candleLimit   string
	candlePeriod  string
	placeOrderApprovalChan chan bool
}

func newWorker(c *Client, sm *model.Symbol, limit, period string, approvalChan chan bool, tw *twilio.TwilioAPI) *Worker {
	w := &Worker{
		twilio:        tw,
		API:           c,
		Symbol:        sm,
		candleLimit:   limit,
		candlePeriod:  period,
		placeOrderApprovalChan: approvalChan,
	}
	return w
}

func (w *Worker)MakeProfit(buyOrSellChan, transControllerChan chan model.Transaction, countGen, countRecover chan float64){
	for{
		select{
		case td := <-buyOrSellChan:
			td.Count = <-countGen
			go func (td model.Transaction){
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
						log.Printf("makeProfit: For %s: type: \"%s\": in %s operation: Goroutine No. %.2f: Ended... |reasonForEnding = %s |sOrbOrderID = %s |\n", w.symbol.ID, td.tYPE, td.operation, td.count, reasonForEnding, td.soldOrBoughtOrderID)
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
					orderID := placeOrder(td.orderQty, bestPrice, w.symbol.SymbolID, td.Side)
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