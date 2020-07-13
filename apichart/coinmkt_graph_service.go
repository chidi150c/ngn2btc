package apichart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//GraphService performs graph operations and services
type GraphService struct {
	session *Session
}

//ReceiveGraphPointFromAPI reieves graph API from provider and sends it into the app
func (b *GraphService) GraphPointFromCoinmkt(done chan bool, sigs chan os.Signal) {
	log.Println("*********Started GraphService.GraphPointFromCoinmkt********")
	defer log.Println("*********End GraphService.GraphPointFromCoinmkt********")
	var graphpoint CGraphPoint
	graphpoint.Data = make([]Coinmarketcap, 1)
	var ngnToUsd = Currencylayer{}
	//To get naira rate
	ngnrate, err := http.Get("http://apilayer.net/api/live?access_key=c0f3e8485d391286b59f12781f51fcf0&currencies=NGN&format=1")
	if err != nil {
		log.Printf("could not connecct to CurrencyLayer for NGN Rate to USD: Error is %v", err)
	}
	ngnrater := json.NewDecoder(ngnrate.Body)
	err = ngnrater.Decode(&ngnToUsd)
	if err != nil {
		log.Fatalf("could not json decode into NgnRate from ngnrate.Body Error is %v", err)
	}
	for {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println(sig)
			done <- true
			return
		case <-time.After(time.Hour):
			ngnrate, err = http.Get("http://apilayer.net/api/live?access_key=c0f3e8485d391286b59f12781f51fcf0&currencies=NGN&format=1")
			if err != nil {
				log.Printf("could not connecct to CurrencyLayer for NGN Rate to USD: Error is %v", err)
			}
			ngnrater := json.NewDecoder(ngnrate.Body)
			err = ngnrater.Decode(&ngnToUsd)
			if err != nil {
				log.Fatalf("could not json decode into NgnRate from ngnrate.Body Error is %v", err)
			}
		default:
			bitResp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/bitcoin/")
			if err != nil {
				log.Printf("could not connecct to bitconnet: Error is %v", err)
				time.Sleep(time.Minute * b.session.GraphPointTime * 60)
				continue
			} else {
				t := time.Now().Unix() * 1000
				graphpoint.Time = strconv.FormatInt(t, 10)
				//.Format("Mon Jan 2 15:04:05")
			}
			bitdecoder := json.NewDecoder(bitResp.Body)
			err = bitdecoder.Decode(&graphpoint.Data)
			if err != nil {
				log.Fatalf("could no json decode into GraphPoint from bitResp.Body Error is %v", err)
			}
			priceUsd, err := strconv.ParseFloat(graphpoint.Data[0].PriceUsd, 64)
			if err != nil {
				log.Fatal(err)
			}
			priceNGN := priceUsd * ngnToUsd.Quotes.USDNGN
			graphpoint.Data[0].PriceUsd = fmt.Sprintf("%.2f", priceNGN)
			b.session.CGraphPointChan <- graphpoint
			log.Println("**In GraphPointFromCoinmkt3 Successfully sent a received graph will wait here for a moment")
			time.Sleep(b.session.GraphPointTime)

			fmt.Println()
		}
	}
}

//CreateGraph4Coinmkt the graph out of the channel of graph api
func (b *GraphService) PopulateChartData() {
	log.Println("*********Started GraphService.PopulateChartData********")
	defer log.Println("*********End GraphService.PopulateChartData********")
	headfootFs, err := os.OpenFile("tools/asset/graphPoint.json", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer headfootFs.Close()

	var bodydata = []string{}

	//get statistics to know the size of file with .Size() method later
	hfStat, e := headfootFs.Stat()
	if e != nil {
		log.Fatal(e)
	}

	var buf = make([]byte, hfStat.Size())
	n, _ := headfootFs.Read(buf)
	if n != 0 {
		b1 := bytes.Replace(buf, []byte("[["), []byte("["), 1)
		b2 := bytes.Replace(b1, []byte("]]"), []byte("]"), 1)
		b3 := bytes.Split(b2, []byte(",\n"))

		fmt.Println()
		for _, v := range b3 {
			bodydata = append(bodydata, string(v))
		}
	}
	for {
		for point := range b.session.CGraphPointChan {
			pointData := point.Data[0].PriceUsd //already converted to NGN
			bodydata = append(bodydata, fmt.Sprintf("[%s,%s]", point.Time, pointData))
			bodydatamodified := strings.Join(bodydata, ",\n")
			headNfoot := fmt.Sprintf("[%s]", bodydatamodified)
			_, err = headfootFs.WriteAt([]byte(headNfoot), 0)
		}
	}

}
