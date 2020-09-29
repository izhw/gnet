package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
		SyncPool(ctx)
	}()
	time.Sleep(2 * time.Second)

	go func() {
		AsyncPool(ctx)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-c:
		fmt.Println("Exit with signal:", sig)
	}
}

func SyncPool(ctx context.Context) {
	log := logger.GlobalSimpleLogger()
	svc := gnet.NewService(
		gcore.WithServiceType(gcore.ServiceTCPPool),
		gcore.WithAddr("127.0.0.1:7777"),
		gcore.WithPoolSize(2, 5),
		gcore.WithHeartbeat([]byte{0}, 30*time.Second),
		//gcore.WithPoolGetTimeout(10*time.Second),
		//gcore.WithPoolIdleTimeout(30*time.Minute),
	)
	p := svc.Pool()
	if err := p.Init(); err != nil {
		log.Error("pool init error:", err)
		return
	}
	defer p.Close()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	data := []byte("Hello sync")

	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for j := 0; j < 5; j++ {
				go func() {
					c, err := p.Get()
					if err != nil {
						log.Error("Pool Get error:", err)
						return
					}
					defer p.Put(c)
					resp, err := c.WriteRead(data)
					if err != nil {
						c.Close()
						log.Error("Pool WriteRead error:", err)
						return
					}
					log.Info("Pool WriteRead:", string(resp))
				}()
			}
		}
	}
}

func AsyncPool(ctx context.Context) {
	log := logger.GlobalSimpleLogger()
	svc := gnet.NewService(
		gcore.WithServiceType(gcore.ServiceTCPAsyncPool),
		gcore.WithAddr("127.0.0.1:7777"),
		gcore.WithEventHandler(NewAsyncHandler()),
		gcore.WithPoolSize(0, 5),
		gcore.WithPoolGetTimeout(5*time.Second),
		gcore.WithPoolIdleTimeout(30*time.Minute),
		gcore.WithHeartbeat([]byte{0}, 30*time.Second),
	)
	p := svc.Pool()
	if err := p.Init(); err != nil {
		log.Error("pool init error:", err)
		return
	}
	defer p.Close()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	data := []byte("Hello async")

	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for j := 0; j < 5; j++ {
				go func() {
					c, err := p.Get()
					if err != nil {
						log.Error("AsyncPool Get error:", err)
						return
					}
					defer p.Put(c)
					if err = c.Write(data); err != nil {
						log.Error("AsyncPool Write error:", err)
						return
					}
				}()
			}
		}
	}
}
