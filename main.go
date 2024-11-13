package main

import (
	"embed"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	. "github.com/retail-ai-inc/beanqui/internal/routers"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"net/http"
)

var (
	port string = ":80"
)

func init() {

	viper.AddConfigPath("./")
	viper.SetConfigName("env")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	//Initialize configuration information
	if err := viper.Unmarshal(&redisx.BqConfig); err != nil {
		log.Fatalf("viper unmarshal err:%+v \n", err)
	}
}

//go:embed ui
var folder embed.FS

func main() {

	//flag.StringVar(&port, "port", ":9090", "port")
	//flag.Parse()

	// init http server
	router := NewRouter()
	// FS static files
	router.File("/", func(ctx *BeanContext) error {
		fd, err := fs.Sub(folder, "ui")
		if err != nil {
			log.Fatalf("static files error:%+v \n", err)
		}
		http.FileServer(http.FS(fd)).ServeHTTP(ctx.Writer, ctx.Request)
		return nil
	})

	router.Get("/ping", ping)
	router.Get("/schedule", Auth(NewSchedule().List))
	router.Get("/queue/list", Auth(NewQueue().List))
	router.Get("/queue/detail", Auth(NewQueue().Detail))
	router.Get("/logs", Auth(NewLogs().List))
	router.Get("/log", Auth(NewLog().List))

	router.Get("/redis", Auth(NewRedisInfo().Info))
	router.Get("/redis/monitor", Auth(NewRedisInfo().Monitor))

	router.Post("/login", NewLogin().Login)
	router.Get("/clients", Auth(NewClient().List))
	router.Get("/dashboard", Auth(NewDashboard().Info))
	router.Get("/event_log/list", Auth(NewEventLog().List))
	router.Get("/event_log/detail", Auth(NewEventLog().Detail))
	router.Delete("/event_log/delete", Auth(NewEventLog().Delete))
	router.Put("/event_log/edit", Auth(NewEventLog().Edit))
	router.Post("/event_log/retry", Auth(NewEventLog().Retry))

	router.Get("/user/list", Auth(NewUser().List))
	router.Post("/user/add", Auth(NewUser().Add))
	router.Delete("/user/del", Auth(NewUser().Delete))
	router.Put("/user/edit", Auth(NewUser().Edit))

	router.Get("/googleLogin", NewLogin().GoogleLogin)
	router.Get("/callback", NewLogin().GoogleCallBack)

	router.Get("/dlq/list", Auth(NewDlq().List))

	log.Printf("server start on port %+v", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalln(err)
	}

}

func ping(ctx *BeanContext) error {

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

	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.Write([]byte("pong"))
	return nil

}
