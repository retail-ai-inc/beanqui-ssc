package main

import (
	"flag"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	. "github.com/retail-ai-inc/beanqui/internal/routers"
	"github.com/spf13/viper"
	"log"
	"net/http"
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
	router := NewRouter()
	router.File("/", NewIndex().File)

	router.Get("/ping", ping)
	router.Get("/schedule", Auth(NewSchedule(client).List))
	router.Get("/queue/list", Auth(NewQueue(client).List))
	router.Get("/queue/detail", Auth(NewQueue(client).Detail))
	router.Get("/logs", Auth(NewLogs(client).List))
	router.Get("/log", Auth(NewLog(client).List))
	router.Get("/redis", Auth(NewRedisInfo(client).Info))
	router.Post("/login", NewLogin().Login)
	router.Get("/clients", Auth(NewClient(client).List))
	router.Get("/dashboard", Auth(NewDashboard(client).Info))

	log.Printf("server start on port %+v", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalln(err)
	}

}

func ping(w http.ResponseWriter, r *http.Request) {

	// clientId := r.Header.Get("Client-Id")
	// clientSecret := r.Header.Get("Client-Secret")
	//
	// if clientId == "" || clientSecret == "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// if clientId != viper.GetString("auth.clientId") || clientSecret != viper.GetString("auth.clientSecret") {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	_, _ = w.Write([]byte("No permission"))
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
	return

}
