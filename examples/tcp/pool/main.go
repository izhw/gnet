package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/pool"
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
	p, err := pool.NewPool("127.0.0.1:7777",
		gnet.WithPoolSize(2, 5),
		//gnet.WithPoolGetTimeout(10*time.Second),
		//gnet.WithPoolIdleTimeout(30*time.Minute),
		//gnet.WithHeartbeat([]byte{0}, 30*time.Second),
	)
	if err != nil {
		fmt.Println("New pool error:", err)
		return
	}
	defer p.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	data := []byte("Hello sync")

	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for j := 0; j < 10; j++ {
				go func() {
					c, err := p.Get()
					if err != nil {
						fmt.Println("Pool Get error:", err)
						return
					}
					defer p.Put(c)
					resp, err := c.WriteRead(data)
					if err != nil {
						c.Close()
						fmt.Println("Pool WriteRead error:", err)
						return
					}
					fmt.Println("Pool WriteRead:", string(resp))
				}()
			}
		}
	}
}

func AsyncPool(ctx context.Context) {
	p, err := pool.NewAsyncPool("127.0.0.1:7777",
		NewAsyncHandler(),
		gnet.WithPoolSize(0, 5),
		//gnet.WithPoolGetTimeout(10*time.Second),
		//gnet.WithPoolIdleTimeout(30*time.Minute),
		//gnet.WithHeartbeat([]byte{0}, 30*time.Second),
	)
	if err != nil {
		fmt.Println("New pool error:", err)
		return
	}
	defer p.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	data := []byte("Hello async")

	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for j := 0; j < 10; j++ {
				go func() {
					c, err := p.Get()
					if err != nil {
						fmt.Println("AsyncPool Get error:", err)
						return
					}
					defer p.Put(c)
					if err = c.Write(data); err != nil {
						fmt.Println("AsyncPool Write error:", err)
						return
					}
				}()
			}
		}
	}
}
