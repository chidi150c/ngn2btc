package model

func countFunc(countRecover chan float64) <-chan float64 {
	//log.Printf("countFunc: For %s: called \n", w.symbol.ID)
	countGen := make(chan float64)
	countHolder := countGen
	var pending []float64
	for i := 1.0; i <= 1000.0; i++ {
		pending = append(pending, i)
	}
	go func(pending []float64, countRecover chan float64) {
		//log.Printf("countFunc: For %s: Goroutine started\n", w.symbol.ID)
		for {
			var first float64
			if len(pending) > 0 {
				////log.Printf("countFunc: for makeProfitFunc: got pending = %v\n", pending)
				first = pending[0]
				countGen = countHolder
				//log.Printf("countFunc: For %s: Total of %d makeProfit Goroutines are at work next one will be assigned %.8f\n", w.symbol.ID, 1000-len(pending), first)
			}
			select {
			case count := <-countRecover:
				pending = append(pending, count)
				//log.Printf("countFunc: For %s: Recovered count %.8f via countRecover and appended to pending whose len is now: %d \n", w.symbol.ID, count, len(pending))
			case countGen <- first:
				//log.Printf("countFunc: For %s: processed count request\n", w.symbol.ID)
				pending = pending[1:]
				countGen = nil
			}

		}
	}(pending, countRecover)
	return countGen
}