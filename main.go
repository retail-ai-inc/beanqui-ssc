package main

import (
	"flag"
	"log"

	"github.com/retail-ai-inc/beanqui/internal/routers"
	"github.com/retail-ai-inc/beanqui/internal/simple_router"
	"github.com/spf13/viper"
)

var port string

func init() {

	viper.AddConfigPath("./")
	viper.SetConfigName("env")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
}
func main() {

	flag.StringVar(&port, "port", ":9090", "port")
	flag.Parse()

	rt := simple_router.New()

	rt.Get("/", routers.IndexHandler)
	rt.Get("/schedule", routers.ScheduleHandler)
	rt.Get("/queue", routers.QueueHandler)
	rt.Get("/log", routers.Auth(routers.LogHandler))
	rt.Get("/redis", routers.RedisHandler)
	rt.Post("/login", routers.LoginHandler)
	rt.Delete("/log/del", routers.Auth(routers.LogDelHandler))
	rt.Post("/log/retry", routers.Auth(routers.LogRetryHandler))
	rt.Post("/log/archive", routers.Auth(routers.LogArchiveHandler))
	rt.Get("/detail", routers.DetailHandler)
	rt.Get("/clients", routers.ClientListHandler)
	rt.Get("/dashboard", routers.DashboardHandler)

	if err := rt.Run(port); err != nil {
		log.Fatalln(err)
	}

}
