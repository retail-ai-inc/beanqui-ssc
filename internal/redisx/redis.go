package redisx

import (
	"log"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanq"
	"github.com/spf13/viper"
)

var (
	redisOnce sync.Once
	client    *redis.Client
	BqConfig  beanq.BeanqConfig
)

func initCfg() {

	vp := viper.New()
	vp.AddConfigPath("./")
	vp.SetConfigName("env")
	vp.SetConfigType("json")
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	if err := vp.Unmarshal(&BqConfig); err != nil {
		log.Fatalln(err)
	}
}

func Client() *redis.Client {

	redisOnce.Do(func() {
		initCfg()
		client = redis.NewClient(&redis.Options{
			Network:  "",
			Addr:     strings.Join([]string{BqConfig.Redis.Host, BqConfig.Redis.Port}, ":"),
			Username: "",
			Password: BqConfig.Redis.Password,
			DB:       BqConfig.Redis.Database,
		})
	})

	return client
}
