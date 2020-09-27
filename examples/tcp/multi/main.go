package main

import (
	"context"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/logger"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.GlobalSimpleLogger()
	service := gnet.NewService(
		gcore.WithServiceType(gcore.ServiceTCPServer|gcore.ServiceTCPAsyncClient),
		gcore.WithLogger(log),
	)

	// client
	c := service.Client()
	err := c.Init(
		gcore.WithAddr("127.0.0.1:7777"),
		gcore.WithEventHandler(NewAsyncHandler()),
	)
	if err != nil {
		log.Fatal("client init error:", err)
	}
	defer c.Close()

	go StartClient(ctx, c)

	// server
	s := service.Server()
	err = s.Init(
		gcore.WithAddr("0.0.0.0:8888"),
		gcore.WithEventHandler(NewServerHandler()),
	)
	if err != nil {
		log.Fatal("server init error:", err)
	}
	log.Fatal("Exit:", s.Serve())
}

func StartClient(ctx context.Context, c gnet.Conn) {
	log := logger.GlobalSimpleLogger()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	data := []byte("multi client")
	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.Write(data); err != nil {
				log.Error("multi client write err:", err)
				return
			}
		}
	}
}
