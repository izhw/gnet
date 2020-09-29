package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for i := 0; i < 10000; i++ {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go Client(i)
			}
		}
	}()
	time.Sleep(10 * time.Millisecond)

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for i := 10001; i < 20000; i++ {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go AsyncClient(i)
			}
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-c:
		fmt.Println("Exit with signal:", sig)
	}
}

func Client(id int) {

	log := logger.GlobalSimpleLogger()
	svc := gnet.NewService(
		gcore.WithServiceType(gcore.ServiceTCPClient),
		gcore.WithAddr("127.0.0.1:7777"),
	)
	c := svc.Client()
	if err := c.Init(); err != nil {
		log.Error("client init error:", err)
		return
	}
	defer c.Close()

	data := []byte("Hello world " + strconv.Itoa(id))
	for i := 0; i < 10; i++ {
		resp, err := c.WriteRead(data)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("recv resp:", string(resp), i)
		time.Sleep(1 * time.Second)
	}
	log.Info("Client Done:", id)
}

func AsyncClient(id int) {
	log := logger.GlobalSimpleLogger()
	svc := gnet.NewService(
		gcore.WithServiceType(gcore.ServiceTCPAsyncClient),
		gcore.WithAddr("127.0.0.1:7777"),
		gcore.WithEventHandler(NewAsyncHandler()),
		gcore.WithHeartbeat([]byte{0}, 5*time.Second),
		gcore.WithLogger(log),
	)
	c := svc.Client()
	if err := c.Init(); err != nil {
		log.Error("async client init error:", err)
		return
	}
	c.SetTag(strconv.Itoa(id))
	defer c.Close()

	data := []byte("Hello world " + strconv.Itoa(id))
	for i := 0; i < 3; i++ {
		if err := c.Write(data); err != nil {
			log.Error(id, "write err:", err)
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Info("AsyncClient Done:", id)
}
