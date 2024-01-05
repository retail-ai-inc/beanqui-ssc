package main

import (
	"flag"
	"log"
	"net/http"

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

	mux := http.NewServeMux()
	mux.Handle("/", &Index{})
	mux.Handle("/schedule", &Schedule{})
	mux.Handle("/queue", &Queue{})
	mux.Handle("/logs", &Logs{})
	// restful: detail,delete,retry,archive
	mux.Handle("/log", &Log{})
	mux.Handle("/redis", &RedisInfo{})
	mux.Handle("/login", &Login{})
	mux.Handle("/clients", &Client{})
	mux.Handle("/dashboard", Auth(&Dashboard{}))

	srv := http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Println("server start.....")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
