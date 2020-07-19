package pool

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/internal/util/limter"
	"github.com/izhw/gnet/tcp/tcpclient"
)

type asyncPool struct {
	addr      string
	opts      gnet.Options
	handler   gnet.EventHandler
	factory   gnet.Factory
	connChan  chan gnet.Conn
	closeChan chan struct{}
	limiter   limter.Limiter
	cancel    context.CancelFunc
	closed    int32
}

func NewAsyncPool(addr string, h gnet.EventHandler, opts ...gnet.Option) (gnet.Pool, error) {
	if addr == "" {
		return nil, gnet.ErrPoolInvalidAddr
	}
	if h == nil {
		h = gnet.DefaultEventHandler()
	}
	p := &asyncPool{
		addr:      addr,
		opts:      gnet.DefaultOptions(),
		handler:   h,
		closeChan: make(chan struct{}),
	}
	for _, o := range opts {
		o(&p.opts)
	}
	if p.opts.PoolMaxSize == 0 {
		p.opts.PoolMaxSize = gnet.DefaultPoolSize
	}
	p.connChan = make(chan gnet.Conn, p.opts.PoolMaxSize)
	p.limiter = limter.NewTimeoutLimiter(p.opts.PoolMaxSize, p.opts.PoolGetTimeout)
	if p.opts.Ctx == nil {
		p.opts.Ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(p.opts.Ctx)
	go func() {
		<-ctx.Done()
		p.Close()
	}()
	p.cancel = cancel
	opts = append(opts, gnet.WithContext(ctx))
	p.factory = func() (gnet.Conn, error) {
		return tcpclient.NewAsyncClient(p.addr, p.handler, opts...)
	}
	for i := 0; i < int(p.opts.PoolInitSize); i++ {
		conn, err := p.createConn()
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("pool:%w", err)
		}
		p.connChan <- conn
	}
	return p, nil
}

func (p *asyncPool) createConn() (gnet.Conn, error) {
	conn, err := p.factory()
	if err != nil {
		return nil, err
	}
	conn.SetTag(p.opts.Tag)
	return conn, nil
}

func (p *asyncPool) Get() (conn gnet.Conn, err error) {
	if !p.limiter.Allow() {
		return nil, gnet.ErrPoolTimeout
	}
	defer func() {
		if err != nil {
			p.limiter.Revert()
		}
	}()

	for {
		select {
		case <-p.closeChan:
			return nil, gnet.ErrPoolClosed
		default:
		}

		select {
		case conn, ok := <-p.connChan:
			if !ok {
				return nil, gnet.ErrPoolClosed
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

func (p *asyncPool) Put(conn gnet.Conn) {
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
func (p *asyncPool) Close() {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return
	}
	p.cancel()
	close(p.closeChan)
}
