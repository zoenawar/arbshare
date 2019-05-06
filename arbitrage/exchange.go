package arbitrage

import (
	"github.com/adshao/go-binance"
	"github.com/spf13/viper"
)

func InitClient() *binance.Client {
	apiKey := viper.GetString("apiKey")
	secretKey := viper.GetString("secretKey")

	client := binance.NewClient(apiKey, secretKey)

	return client
}