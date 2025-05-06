package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/spl3g/lab2/internal/proxyproto"
	"github.com/spl3g/lab2/services/permissions-service/internal/config"
	"github.com/spl3g/lab2/services/permissions-service/internal/service"

	"google.golang.org/grpc"
)

func run() error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp4", ":"+conf.Port)
	if err != nil {
		return err
	}

	errChan := make(chan error)

	srv := grpc.NewServer()

	svc, err := service.New(conf)
	if err != nil {
		return err
	}

	proxyproto.RegisterCentrifugoProxyServer(srv, svc)

	exitCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

		cancel()

		srv.GracefulStop()

		close(errChan)

		if err := listener.Close(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		errChan <- srv.Serve(listener)
	}()

	select {
	case err := <-errChan:
		return err
	case <-exitCtx.Done():
		log.Println("exit")
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
