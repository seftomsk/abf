package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seftomsk/abf/cmd"
	"github.com/seftomsk/abf/internal/access"
	"github.com/seftomsk/abf/internal/limiter"
	"github.com/seftomsk/abf/internal/server/web"
)

func main() {
	cfg := cmd.Execute()

	storage, err := access.GetStorage(cfg.Storage)
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}

	loginLimiter := limiter.NewLimiter(
		cfg.LoginLimiter.Capacity,
		time.Duration(cfg.LoginLimiter.CountSeconds)*time.Second)
	passwordLimiter := limiter.NewLimiter(
		cfg.PasswordLimiter.Capacity,
		time.Duration(cfg.PasswordLimiter.CountSeconds)*time.Second)
	ipLimiter := limiter.NewLimiter(
		cfg.IPLimiter.Capacity,
		time.Duration(cfg.IPLimiter.CountSeconds)*time.Second)
	l := limiter.NewMultiLimiter(
		loginLimiter,
		passwordLimiter,
		ipLimiter)
	a := access.NewIPAccess(storage)

	server := web.NewServer(l, a, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Println("failed to stop web server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		log.Println("failed to start web server: " + err.Error())
		cancel()
	}
}
