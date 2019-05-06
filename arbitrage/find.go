package arbitrage

import (
	"fmt"
	"github.com/adshao/go-binance"
	"os"
	"time"
)

func FetchDepth(symbols []string) {
	var latestAsks []binance.Ask
	wsPartialDepthHandler := func(event *binance.WsPartialDepthEvent) {
		fmt.Println(event.Symbol, event.Asks)
		latestAsks = event.Asks
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	symbolMap := make(map[string]string)
	for i := 0; i < len(symbols); i++ {
		symbolMap[symbols[i]] = "5"
	}
	doneC, stopC, err := binance.WsCombinedPartialDepthServe(symbolMap, wsPartialDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// use stopC to exit
	go func() {
		time.Sleep(40 * time.Second)
		stopC <- struct{}{}
	}()

	// remove this if you do not want to be blocked here
	<-doneC
}

func ComposePaths(targetCoin string, routeCoins []string) [][]string {
	var paths [][]string
	for i := 0; i < len(routeCoins); i++ {
		route := routeCoins[i]
		path := []string{targetCoin, route}
		for j := 0; j < len(routeCoins); j++ {
			fullPath := path
			if routeCoins[j] != route {
				fullPath = append(fullPath, routeCoins[j])
				paths = append(paths, fullPath)
			}
		}
	}
	return paths
}

/*
func CompareSymbols(a binance.Ask, b binance.Ask, c binance.Ask) {

}
*/