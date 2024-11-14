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
	router.File("/", HeaderRule(func(ctx *BeanContext) error {
		fd, err := fs.Sub(folder, "ui")
		if err != nil {
			log.Fatalf("static files error:%+v \n", err)
		}
		http.FileServer(http.FS(fd)).ServeHTTP(ctx.Writer, ctx.Request)
		return nil
	}))

	router.Get("/ping", HeaderRule(ping))
	router.Get("/schedule", MigrateMiddleWare(NewSchedule().List))
	router.Get("/queue/list", MigrateMiddleWare(NewQueue().List))
	router.Get("/queue/detail", MigrateMiddleWare(NewQueue().Detail))
	router.Get("/logs", MigrateMiddleWare(NewLogs().List))
	router.Get("/log", MigrateMiddleWare(NewLog().List))

	router.Get("/redis", MigrateMiddleWare(NewRedisInfo().Info))
	router.Get("/redis/monitor", MigrateMiddleWare(NewRedisInfo().Monitor))

	router.Post("/login", HeaderRule(NewLogin().Login))
	router.Get("/clients", MigrateMiddleWare(NewClient().List))
	router.Get("/dashboard", MigrateMiddleWare(NewDashboard().Info))
	router.Get("/event_log/list", MigrateMiddleWare(NewEventLog().List))
	router.Get("/event_log/detail", MigrateMiddleWare(NewEventLog().Detail))
	router.Post("/event_log/delete", MigrateMiddleWare(NewEventLog().Delete))
	router.Post("/event_log/edit", MigrateMiddleWare(NewEventLog().Edit))
	router.Post("/event_log/retry", MigrateMiddleWare(NewEventLog().Retry))

	router.Get("/user/list", MigrateMiddleWare(NewUser().List))
	router.Post("/user/add", MigrateMiddleWare(NewUser().Add))
	router.Post("/user/del", MigrateMiddleWare(NewUser().Delete))
	router.Post("/user/edit", MigrateMiddleWare(NewUser().Edit))

	router.Get("/googleLogin", NewLogin().GoogleLogin)
	router.Get("/callback", NewLogin().GoogleCallBack)

	router.Get("/dlq/list", MigrateMiddleWare(NewDlq().List))

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
