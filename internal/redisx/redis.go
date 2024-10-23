package redisx

import (
	"github.com/retail-ai-inc/beanq/v3"
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisOnce sync.Once
	client    redis.UniversalClient
	BqConfig  beanq.BeanqConfig
)

func Client() redis.UniversalClient {

	redisOnce.Do(func() {

		if err := viper.Unmarshal(&BqConfig); err != nil {
			log.Fatalf("viper unmarshal err:%+v \n", err)
		}
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
			PoolSize:     BqConfig.Redis.PoolSize,
			MinIdleConns: BqConfig.Redis.MinIdleConnections,
			PoolTimeout:  BqConfig.Redis.PoolTimeout,
			PoolFIFO:     false,
		})
	})

	return client
}
