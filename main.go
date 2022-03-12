package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Listen string `env:"HISTORY_SERVER_LISTEN_ADDR" envDefault:":8080"`
	TTL    uint   `env:"LOCATION_HISTORY_TTL_SECONDS"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	startServer(&cfg)
}

func startServer(cfg *Config) {
	r := gin.Default()
	ctrl := Controller{
		history: NewOrdersHistory(),
	}
	orders := r.Group("/location/:order_id")
	orders.POST("/now", ctrl.PostOrderHistory)
	orders.GET("", ctrl.GetOrderHistory)
	orders.DELETE("", ctrl.DeleteOrderHistory)
	r.Run(cfg.Listen)
}
