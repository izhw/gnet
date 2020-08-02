# gnet

## Installation

```
go get -u github.com/izhw/gnet
```

## Features
Building net service quickly with functional options
* [x] TCP Server and Client
* [x] Connection pool
* [ ] gRPC Srever and Client
* [ ] WebSocket Server and Client

## Quick start

#### Examples

* [tcp-server](https://github.com/izhw/gnet/tree/master/examples/tcp/server)
* [tcp-client](https://github.com/izhw/gnet/tree/master/examples/tcp/client)
* [conn-pool](https://github.com/izhw/gnet/tree/master/examples/tcp/pool)

#### TCP Server

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/izhw/gnet"
    "github.com/izhw/gnet/tcp/tcpserver"
)

type ServerHandler struct {
    *gnet.NetEventHandler
}

func (h *ServerHandler) OnReadMsg(c gnet.Conn, data []byte) error {
    fmt.Println("read msg:", string(data))
    c.Write(data)
    return nil
}

func main() {
    s := tcpserver.NewServer("0.0.0.0:7777", &ServerHandler{})
    log.Fatal("Exit:", s.Serve())
}
```

#### TCP Client
* Sync mode
```go
package main

import (
    "fmt"
    "os"
    
    "github.com/izhw/gnet/tcp/tcpclient"
)

func main() {
    c, err := tcpclient.NewClient("127.0.0.1:7777")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer c.Close()
    // todo...
    data := []byte("Hello world")
    resp, err := c.WriteRead(data)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("recv resp:", string(resp))
}
```
* Async mode
```go
package main

import (
    "fmt"
    "os"
    "time"
    
    "github.com/izhw/gnet"
    "github.com/izhw/gnet/tcp/tcpclient"
)

type AsyncHandler struct {
    *gnet.NetEventHandler
}

func (h *AsyncHandler) OnReadMsg(c gnet.Conn, data []byte) error {
    fmt.Println("read msg:", string(data))
    return nil
}

func main() {
    c, err := tcpclient.NewAsyncClient("127.0.0.1:7777", &AsyncHandler{})
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer c.Close()
    
    data := []byte("Hello world")
    for i := 0; i < 10; i++ {
        if err := c.Write(data); err != nil {
            fmt.Println("write err:", err)
            os.Exit(1)
        }
        time.Sleep(1 * time.Second)
    }
    fmt.Println("AsyncClient done")
}
```

#### Connection pool
* Sync mode
```go
package main

import (
    "fmt"
    "os"
    
    "github.com/izhw/gnet"
    "github.com/izhw/gnet/pool"
)

func main() {
    p, err := pool.NewPool("127.0.0.1:7777",
        gnet.WithPoolSize(5, 10),
        //gnet.WithPoolGetTimeout(10*time.Second),
        //gnet.WithPoolIdleTimeout(30*time.Minute),
        //gnet.WithHeartbeat([]byte{0}, 30*time.Second),
    )
    if err != nil {
        fmt.Println("New pool error:", err)
        os.Exit(1)
    }
    defer p.Close()
    
    c, err := p.Get()
    if err != nil {
        fmt.Println("Pool Get error:", err)
        os.Exit(1)
    }
    defer p.Put(c)
    
    resp, err := c.WriteRead([]byte("Hello world"))
    if err != nil {
        fmt.Println("Pool WriteRead error:", err)
        os.Exit(1)
    }
    fmt.Println("Pool WriteRead:", string(resp))
}
```
* Async mode
```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/pool"
)

type AsyncHandler struct {
	*gnet.NetEventHandler
}

func (h *AsyncHandler) OnReadMsg(c gnet.Conn, data []byte) error {
	fmt.Println("Pool read msg:", string(data))
	return nil
}

func main() {
	p, err := pool.NewAsyncPool("127.0.0.1:7777",
		&AsyncHandler{},
		gnet.WithPoolSize(0, 10),
		//gnet.WithPoolGetTimeout(10*time.Second),
		//gnet.WithPoolIdleTimeout(30*time.Minute),
		//gnet.WithHeartbeat([]byte{0}, 30*time.Second),
	)
	if err != nil {
		fmt.Println("New pool error:", err)
		os.Exit(1)
	}
	defer p.Close()

	c, err := p.Get()
	if err != nil {
		fmt.Println("Pool Get error:", err)
		os.Exit(1)
	}
	defer p.Put(c)

	if err := c.Write([]byte("Hello world")); err != nil {
		fmt.Println("Pool Write error:", err)
		os.Exit(1)
	}
	time.Sleep(3 * time.Second)
}
```

#### Functional options for gnet
for example:
```go
    s := tcpserver.NewServer("0.0.0.0:7777",
        NewServerHandler(),
        gnet.WithLogger(logger.NewSimpleLogger()),
        gnet.WithConnNumLimit(50),
        gnet.WithHeaderCodec(&protocol.CodecProtoVarint{}),
        ...,
    )
```
See `options.go` for more options

