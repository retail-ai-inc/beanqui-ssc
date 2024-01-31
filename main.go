package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	. "github.com/retail-ai-inc/beanqui/internal/routers"
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

	// init redis
	client := redisx.Client()

	// init http server
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	mux.Handle("/", NewIndex())
	mux.Handle("/schedule", Auth(NewSchedule(client)))
	// queue:list, detail
	mux.Handle("/queue", Auth(NewQueue(client)))
	mux.Handle("/logs", Auth(NewLogs(client)))
	// restful: detail,delete,retry,archive
	mux.Handle("/log", Auth(NewLog(client)))
	mux.Handle("/redis", Auth(NewRedisInfo(client)))
	mux.Handle("/login", NewLogin())
	mux.Handle("/clients", Auth(NewClient(client)))
	mux.Handle("/dashboard", Auth(NewDashboard(client)))

	srv := http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Printf("server start on port %s \n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
