package arbitrage

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance"
	"os"
	"time"
)

func WsFetchDepth(symbols []string) {
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
		path := []string{route + targetCoin}
		for j := 0; j < len(routeCoins); j++ {
			fullPath := path
			if routeCoins[j] != route {
				fullPath = append(fullPath, routeCoins[j]+route, routeCoins[j]+targetCoin)
				paths = append(paths, fullPath)
			}
		}
	}
	fmt.Println(paths)
	return paths
}

func CompareSymbols(path []string, client *binance.Client) {
	a, _ := client.NewDepthService().Symbol(path[0]).Do(context.Background())
	b, _ := client.NewDepthService().Symbol(path[0]).Do(context.Background())
	c, _ := client.NewDepthService().Symbol(path[0]).Do(context.Background())
}
