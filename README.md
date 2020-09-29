# gnet

## Installation

```
go get -u github.com/izhw/gnet
```

## Features
Building net service quickly with functional options
* [x] TCP Server and Client
* [x] Connection pools
* [ ] gRPC Server and Client
* [ ] WebSocket Server and Client

## Quick start

#### Examples

* [tcp-server](https://github.com/izhw/gnet/tree/master/examples/tcp/server)
* [tcp-client](https://github.com/izhw/gnet/tree/master/examples/tcp/client)
* [conn-pool](https://github.com/izhw/gnet/tree/master/examples/tcp/pool)
* [multi](https://github.com/izhw/gnet/tree/master/examples/tcp/multi)

#### TCP Server

```go
package main

import (
    "fmt"
    
    "github.com/izhw/gnet"
    "github.com/izhw/gnet/gcore"
    "github.com/izhw/gnet/logger"
)

type ServerHandler struct {
    *gcore.NetEventHandler
}

func (h *ServerHandler) OnReadMsg(c gcore.Conn, data []byte) error {
    fmt.Println("server read msg:", string(data))
    c.Write(data)
    return nil
}

func main() {
    log := logger.GlobalSimpleLogger()
    svc := gnet.NewService(
        gcore.WithServiceType(gcore.ServiceTCPServer),
        gcore.WithAddr("0.0.0.0:7777"),
        gcore.WithEventHandler(&ServerHandler{}),
        gcore.WithLogger(log),
    )
    s := svc.Server()
    if err := s.Init(); err != nil {
        log.Fatal("server init error:", err)
    }
    log.Fatal("server exit:", s.Serve())
}
```

#### TCP Client

* Sync mode
```go
package main

import (
    "fmt"
    "os"
    
    "github.com/izhw/gnet"
    "github.com/izhw/gnet/gcore"
)

func main() {
    svc := gnet.NewService(
        gcore.WithServiceType(gcore.ServiceTCPClient),
        gcore.WithAddr("127.0.0.1:7777"),
    )
    c := svc.Client()
    if err := c.Init(); err != nil {
        fmt.Println("client init error:", err)
        os.Exit(1)
    }
    defer c.Close()
    // todo...
    data := []byte("Hello world")
    resp, err := c.WriteRead(data)
    if err != nil {
        fmt.Println("client WriteRead error:", err)
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
    "github.com/izhw/gnet/gcore"
    "github.com/izhw/gnet/logger"
)

type AsyncHandler struct {
    *gcore.NetEventHandler
}

func (h *AsyncHandler) OnReadMsg(c gcore.Conn, data []byte) error {
    fmt.Println("async client read msg:", string(data))
    return nil
}

func main() {
    log := logger.GlobalSimpleLogger()
    svc := gnet.NewService(
        gcore.WithServiceType(gcore.ServiceTCPAsyncClient),
        gcore.WithAddr("127.0.0.1:7777"),
        gcore.WithEventHandler(&AsyncHandler{}),
        gcore.WithLogger(log),
    )
    c := svc.Client()
    if err := c.Init(); err != nil {
        log.Fatal("client init error:", err)
    }
    defer c.Close()
    
    data := []byte("Hello world")
    for i := 0; i < 10; i++ {
        if err := c.Write(data); err != nil {
            log.Fatal("client write err:", err)
        }
        time.Sleep(1 * time.Second)
    }
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
    "github.com/izhw/gnet/gcore"
)

func main() {
    svc := gnet.NewService(
        gcore.WithServiceType(gcore.ServiceTCPPool),
        gcore.WithAddr("127.0.0.1:7777"),
        gcore.WithPoolSize(5, 10),
    )
    p := svc.Pool()
    if err := p.Init(); err != nil {
        fmt.Println("pool init error:", err)
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
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/logger"
)

type AsyncHandler struct {
	*gcore.NetEventHandler
}

func (h *AsyncHandler) OnReadMsg(c gcore.Conn, data []byte) error {
	fmt.Println("Pool read msg:", string(data))
	return nil
}

func main() {
    log := logger.GlobalSimpleLogger()
    svc := gnet.NewService(
        gcore.WithServiceType(gcore.ServiceTCPAsyncPool),
        gcore.WithAddr("127.0.0.1:7777"),
        gcore.WithEventHandler(&AsyncHandler{}),
        gcore.WithPoolSize(0, 10),
        gcore.WithPoolGetTimeout(5*time.Second),
        gcore.WithPoolIdleTimeout(30*time.Minute),
        //gcore.WithHeartbeat([]byte{0}, 30*time.Second),
    )
    p := svc.Pool()
    if err := p.Init(); err != nil {
        log.Fatal("pool init error:", err)
    }
    defer p.Close()
    
    c, err := p.Get()
    if err != nil {
        log.Fatal("Pool Get error:", err)
    }
    defer p.Put(c)
    
    if err = c.Write([]byte("Hello world")); err != nil {
        log.Fatal("Pool Write error:", err)
    }
    time.Sleep(3 * time.Second)
}
```

#### Multi

You can build `server`, `client` and `pool` in one `service` by setting the `ServiceType`:
`gcore.WithServiceType(gcore.ServiceTCPServer|gcore.ServiceTCPAsyncClient|gcore.ServiceTCPAsyncPool)`

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/logger"
)

type ServerHandler struct {
    *gcore.NetEventHandler
}

func (h *ServerHandler) OnReadMsg(c gcore.Conn, data []byte) error {
    fmt.Println("server read msg:", string(data))
    c.Write(data)
    return nil
}

type AsyncHandler struct {
    *gcore.NetEventHandler
}

func (h *AsyncHandler) OnReadMsg(c gcore.Conn, data []byte) error {
    fmt.Println("multi client read msg:", string(data))
    return nil
}

func main() {

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    log := logger.GlobalSimpleLogger()
    svc := gnet.NewService(
        gcore.WithServiceType(gcore.ServiceTCPServer|gcore.ServiceTCPAsyncClient),
        gcore.WithLogger(log),
    )
    
    // client
    c := svc.Client()
    err := c.Init(
        gcore.WithAddr("127.0.0.1:7777"),
        gcore.WithEventHandler(&AsyncHandler{}),
    )
    if err != nil {
        log.Fatal("client init error:", err)
    }
    defer c.Close()
    
    go StartClient(ctx, c)
    
    // server
    s := svc.Server()
    err = s.Init(
        gcore.WithAddr("0.0.0.0:7778"),
        gcore.WithEventHandler(&ServerHandler{}),
    )
    if err != nil {
        log.Fatal("server init error:", err)
    }
    log.Fatal("Exit:", s.Serve())
}

func StartClient(ctx context.Context, c gnet.Client) {
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
```

#### Functional options for gnet
for example:
```go
    svc := gnet.NewService(
        gcore.WithServiceType(gcore.ServiceTCPServer),
        gcore.WithAddr("0.0.0.0:7777"),
        gcore.WithEventHandler(&ServerHandler{}),
        gcore.WithLogger(logger.GlobalSimpleLogger()),
        gcore.WithHeaderCodec(&protocol.CodecProtoVarint{}),
        gcore.WithReadTimeout(2 * time.Minute),
        gcore.WithConnNumLimit(1000),
        gcore.WithHeartbeat([]byte{0}, 30 * time.Second)
        ...
    )
```
See `gcore/options.go` for more options

