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
	"github.com/izhw/gnet/logger"
	"github.com/izhw/gnet/tcp/tcpclient"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
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
		ticker := time.NewTicker(1 * time.Second)
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

	l := logger.NewSimpleLogger()
	c, err := tcpclient.NewClient("127.0.0.1:7777",
		gnet.WithLogger(l),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	data := []byte("Hello world " + strconv.Itoa(id))
	for i := 0; i < 10; i++ {
		resp, err := c.WriteRead(data)
		if err != nil {
			l.Info(err)
			return
		}
		l.Info("recv resp:", string(resp), i)
		time.Sleep(1 * time.Second)
	}
	l.Info("Client Done:", id)
}

func AsyncClient(id int) {
	l := logger.NewSimpleLogger()
	c, err := tcpclient.NewAsyncClient("127.0.0.1:7777",
		NewAsyncHandler(),
		gnet.WithLogger(l),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.SetTag(strconv.Itoa(id))
	defer c.Close()

	data := []byte("Hello world " + strconv.Itoa(id))
	for i := 0; i < 10; i++ {
		if err := c.Write(data); err != nil {
			l.Info(id, "write err:", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	l.Info("AsyncClient Done:", id)
}
