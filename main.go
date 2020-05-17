package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task/app/models"
	"task/conf"
	"task/pkg/cron"
	"task/routers"
	"time"
)

func main() {
	s := flag.Bool("start", false, "start server")
	i := flag.Bool("init", false, "Init Table")
	_ = flag.String("env", "local", "env")
	flag.Parse()
	if *s == true {
		start()
	}
	if *i == true {
		initTable()
	}
}

func start() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.Config.Server.HTTPPort),
		Handler:        router,
		ReadTimeout:    conf.Config.Server.ReadTimeout,
		WriteTimeout:   conf.Config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := cron.C.AddFunc(conf.Config.App.Cron, func() {
		models.Remind(conf.Config.App.RemindDays)
	})
	if err != nil {
		log.Fatal("add cron func fail:", err)
	}
	cron.C.Start()

	go func() {
		// service connections
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalf("start Serve fails : %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func initTable() {
	models.InitTable()
}
