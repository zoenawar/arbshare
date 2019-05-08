package arbitrage

import (
	"encoding/json"
	"fmt"
    "strconv"
	"github.com/adshao/go-binance"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redisAddr"),
		Password: viper.GetString("redisPass"),
		DB:       0,
	})

	return client
}

func HydrateRedis(symbols []string, client *redis.Client) {
	done := make(chan struct{})
	feed := make(chan *binance.WsPartialDepthEvent, 1)

	go func() {
		WsFetchDepth(symbols, feed)
	}()

	go func() {
		for event := range feed {
			askJson, _ := json.Marshal(event.Asks)
			bidJson, _ := json.Marshal(event.Bids)

			err := client.HSet(event.Symbol, "Asks", askJson).Err()
			err = client.HSet(event.Symbol, "Bids", bidJson).Err()
			if err != nil {
				fmt.Println("Failed to update Redis: ", err)
				continue
			}
			fmt.Println("Updated symbol ", event.Symbol)
		}
	}()

	<-done
}

func GetLatestCoinData(symbol string, field string, client *redis.Client) (float64, float64) {
    dataMarshalled, _ := client.HGet(symbol, field).Result()

    //Temporary type for holding Ask/Bid data
    type Data struct {
        Price string
        Quantity string
    }

    dataList := make([]Data, 1)
    err := json.Unmarshal([]byte(dataMarshalled), &dataList)
    if err != nil {
        fmt.Println(err)
    }

    coinPrice, _ := strconv.ParseFloat(dataList[0].Price, 64)
    coinQuantity, _ := strconv.ParseFloat(dataList[0].Quantity, 64)

    return coinPrice, coinQuantity
}
