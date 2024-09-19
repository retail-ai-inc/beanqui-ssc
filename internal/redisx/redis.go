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
	client    redis.UniversalClient
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

func Client() redis.UniversalClient {

	redisOnce.Do(func() {
		initCfg()
		hosts := strings.Split(BqConfig.Redis.Host, ",")
		for i, h := range hosts {
			hs := strings.Split(h, ":")
			if len(hs) == 1 {
				hosts[i] = strings.Join([]string{h, BqConfig.Redis.Port}, ":")
			}
		}
		client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        hosts,
			Password:     BqConfig.Redis.Password,
			DB:           BqConfig.Redis.Database,
			MaxRetries:   BqConfig.JobMaxRetries,
			DialTimeout:  BqConfig.Redis.DialTimeout,
			ReadTimeout:  BqConfig.Redis.ReadTimeout,
			WriteTimeout: BqConfig.Redis.WriteTimeout,
			PoolSize:     BqConfig.PoolSize,
			MinIdleConns: BqConfig.Redis.MinIdleConnections,
			PoolTimeout:  BqConfig.Redis.PoolTimeout,
			PoolFIFO:     false,
		})
	})

	return client
}
