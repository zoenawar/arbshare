package arbitrage

import (
	"fmt"
	"github.com/adshao/go-binance"
	"os"
	"time"
)

func WsFetchDepth(symbols []string, out chan<- *binance.WsPartialDepthEvent) {
	var latestAsks []binance.Ask
	wsPartialDepthHandler := func(event *binance.WsPartialDepthEvent) {
		latestAsks = event.Asks
		if len(latestAsks) > 0 {
			fmt.Println("Event processed", event.Symbol)
			out <- event
		}
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

func ComposePaths(targetCoin string, routeCoins []string) ([][]string, []string) {
	var paths [][]string
	var pairs []string
	pairsMap := make(map[string]bool)
	for i := 0; i < len(routeCoins); i++ {
		route := routeCoins[i]
		path := []string{route + targetCoin}
		pairsMap[route+targetCoin] = true
		for j := 0; j < len(routeCoins); j++ {
			fullPath := path
			if routeCoins[j] != route {
				fullPath = append(fullPath, routeCoins[j]+route, routeCoins[j]+targetCoin)
				pairsMap[routeCoins[j]+route] = true
				pairsMap[routeCoins[j]+targetCoin] = true
				paths = append(paths, fullPath)
			}
		}
	}

	for key := range pairsMap {
		pairs = append(pairs, key)
	}

	fmt.Println(paths)
	return paths, pairs
}

func CompareSymbols(path []string) {
	var a, b, c binance.Ask
	results := make(chan *binance.WsPartialDepthEvent, 1)
	done := make(chan struct{})
	eventC := make(chan struct{})
	go func() {
		WsFetchDepth(path, results)
	}()

	go func() {
		for event := range results {
			switch event.Symbol {
			case path[0]:
				a = event.Asks[0]
			case path[1]:
				b = event.Asks[0]
			case path[2]:
				c = event.Asks[0]
			}
			eventC <- struct{}{}
		}
	}()

	go func() {
		for range eventC {
			fmt.Println(path[0], a)
			fmt.Println(path[1], b)
			fmt.Println(path[2], c)
		}
	}()

	<-done
}
