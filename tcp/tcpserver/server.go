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

package tcpserver

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/internal/util"
	"github.com/izhw/gnet/logger"
)

const DefaultAddr = "0.0.0.0:7777"

var _ gnet.Server = &server{}

type server struct {
	opts     gnet.Options
	addr     string
	handler  gnet.EventHandler
	listener net.Listener
	limiter  *util.Limiter
	stopChan chan struct{}
	wg       sync.WaitGroup
	connNum  uint32
	stopped  int32
}

func NewServer(addr string, h gnet.EventHandler, opts ...gnet.Option) gnet.Server {
	if addr == "" {
		addr = DefaultAddr
	}
	s := &server{
		opts:    gnet.DefaultOptions(),
		addr:    addr,
		handler: h,
		stopped: 1,
	}
	for _, o := range opts {
		o(&s.opts)
	}
	return s
}

func (s *server) Serve() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = l
	if s.handler == nil {
		s.handler = gnet.DefaultEventHandler()
	}
	if s.opts.Logger == nil {
		s.opts.Logger = logger.DefaultLogger()
	}
	if s.opts.ConnLimit > 0 {
		s.limiter = util.NewLimiter(s.opts.ConnLimit)
	}

	s.wg.Add(1)
	go s.work()
	s.stopChan = make(chan struct{})
	s.stopped = 0

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	// Wait on signal
	select {
	case <-s.stopChan:
		s.wait()
	case sig := <-c:
		s.Stop()
		return errors.New("signal:" + sig.String())
	}
	return nil
}

func (s *server) wait() {
	s.wg.Wait()
	for atomic.LoadUint32(&s.connNum) > 0 {
	}
}

func (s *server) Stop() {
	if !atomic.CompareAndSwapInt32(&s.stopped, 0, 1) {
		return
	}
	close(s.stopChan)
	s.listener.Close()
	s.wait()
}

func (s *server) ConnNum() uint32 {
	return atomic.LoadUint32(&s.connNum)
}

func (s *server) onConnClose() {
	if s.limiter != nil {
		s.limiter.Revert()
	}
	atomic.AddUint32(&s.connNum, ^uint32(0))
}

func (s *server) work() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		s.wg.Done()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				s.opts.Logger.Warnf("TCP server accept temp error:[%v]", ne)
				time.Sleep(time.Second)
				continue
			}
			select {
			case <-s.stopChan:
				return
			default:
			}
			s.opts.Logger.Warnf("TCP server accept error:[%v]", err)
			return
		}
		if s.limiter != nil && !s.limiter.Allow() {
			conn.Close()
			s.opts.Logger.Warnf("TCP server accepted max num:%d, new conn rejected", s.opts.ConnLimit)
			continue
		}
		tcpConn := conn.(*net.TCPConn)

		// TCP keepalive
		if err = tcpConn.SetKeepAlive(true); err != nil {
			s.opts.Logger.Warnf("TCP server conn:%s SetKeepAlive error:[%v]", tcpConn.RemoteAddr(), err)
		}
		if err = tcpConn.SetKeepAlivePeriod(time.Minute * 1); err != nil {
			s.opts.Logger.Warnf("TCP server conn:%s SetKeepAlivePeriod error:[%v]", tcpConn.RemoteAddr(), err)
		}
		// setting keepalive retry count and interval
		if err := setKeepaliveParameters(tcpConn, 6, 10); err != nil {
			s.opts.Logger.Warnf("TCP server conn:%s setKeepaliveParameters error:[%v]", tcpConn.RemoteAddr(), err)
		}
		// new conn
		newConn(ctx, s, tcpConn)
		atomic.AddUint32(&s.connNum, 1)
	}
}
