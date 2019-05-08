package arbitrage

import (
	"encoding/json"
	"fmt"
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
