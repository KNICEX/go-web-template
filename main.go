package go_web_template

import (
	"context"
	"errors"
	"fmt"
	"go-web-template/bootstrap"
	"go-web-template/pkg/conf"
	"go-web-template/pkg/logger"
	"go-web-template/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bootstrap.Init()

	r := routes.Init()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", conf.ServerConf.Port),
		Handler: r,
	}

	go func() {
		fmt.Printf("Server started at :%s%s\n", conf.ServerConf.Port, conf.ServerConf.Prefix)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Error("server shutdown error", err)
	}

	if err := bootstrap.Shutdown(); err != nil {
		logger.L().Error("bootstrap shutdown error", err)
	}

	fmt.Println("Server exiting")
}
