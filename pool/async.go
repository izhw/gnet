// Copyright (c) 2020 izhw
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pool

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/internal/util/limter"
	"github.com/izhw/gnet/tcp/client"
)

var _ gcore.Pool = &AsyncPool{}

type AsyncPool struct {
	opts      gcore.Options
	factory   func() (gcore.Conn, error)
	connChan  chan gcore.Conn
	closeChan chan struct{}
	limiter   limter.Limiter
	cancel    context.CancelFunc
	closed    int32
}

func NewAsyncPool() *AsyncPool {
	return &AsyncPool{}
}

func (p *AsyncPool) WithOptions(opts gcore.Options) {
	p.opts = opts
}

func (p *AsyncPool) Init(opts ...gcore.Option) error {
	for _, opt := range opts {
		opt(&p.opts)
	}
	if p.opts.Addr == "" {
		return gcore.ErrPoolInvalidAddr
	}
	if p.opts.PoolMaxSize == 0 {
		p.opts.PoolMaxSize = 16
	}
	p.connChan = make(chan gcore.Conn, p.opts.PoolMaxSize)
	p.closeChan = make(chan struct{})
	p.limiter = limter.NewTimeoutLimiter(p.opts.PoolMaxSize, p.opts.PoolGetTimeout)
	ctx, cancel := context.WithCancel(p.opts.Ctx)
	go func() {
		<-ctx.Done()
		p.Close()
	}()
	p.cancel = cancel
	p.opts.Ctx = ctx

	p.factory = func() (gcore.Conn, error) {
		c := client.NewAsyncClient()
		c.WithOptions(p.opts)
		if err := c.Init(); err != nil {
			return nil, err
		}
		return c, nil
	}
	for i := 0; i < int(p.opts.PoolInitSize); i++ {
		conn, err := p.createConn()
		if err != nil {
			p.Close()
			return fmt.Errorf("pool:%w", err)
		}
		p.connChan <- conn
	}
	return nil
}

func (p *AsyncPool) createConn() (gcore.Conn, error) {
	conn, err := p.factory()
	if err != nil {
		return nil, err
	}
	conn.SetTag(p.opts.Tag)
	return conn, nil
}

func (p *AsyncPool) Get() (conn gcore.Conn, err error) {
	if !p.limiter.Allow() {
		return nil, gcore.ErrPoolTimeout
	}
	defer func() {
		if err != nil {
			p.limiter.Revert()
		}
	}()

	for {
		select {
		case <-p.closeChan:
			return nil, gcore.ErrPoolClosed
		default:
		}

		select {
		case conn, ok := <-p.connChan:
			if !ok {
				return nil, gcore.ErrPoolClosed
			}
			if conn.Closed() {
				continue
			}
			return conn, nil
		default:
			return p.createConn()
		}
	}
}

func (p *AsyncPool) Put(conn gcore.Conn) {
	if conn == nil {
		return
	}
	p.limiter.Revert()

	select {
	case <-p.closeChan:
		return
	default:
	}
	if conn.Closed() {
		return
	}
	p.connChan <- conn
}

// Close
func (p *AsyncPool) Close() {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return
	}
	p.cancel()
	close(p.closeChan)
}
