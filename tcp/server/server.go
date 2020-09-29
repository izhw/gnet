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

package server

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

	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/internal/util/limter"
)

const DefaultAddr = "0.0.0.0:7777"

var _ gcore.Server = &Server{}

type Server struct {
	opts     gcore.Options
	listener net.Listener
	limiter  limter.Limiter
	stopChan chan struct{}
	wg       sync.WaitGroup
	heartLen uint32
	connNum  uint32
	stopped  int32
}

func NewServer() *Server {
	return &Server{
		stopped: 1,
	}
}

func (s *Server) WithOptions(opts gcore.Options) {
	s.opts = opts
}

func (s *Server) Init(opts ...gcore.Option) error {
	for _, opt := range opts {
		opt(&s.opts)
	}
	if s.opts.Addr == "" {
		s.opts.Addr = DefaultAddr
	}
	l, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	s.listener = l
	if s.opts.ConnLimit > 0 {
		s.limiter = limter.NewLimiter(s.opts.ConnLimit)
	}
	s.stopChan = make(chan struct{})
	s.heartLen = uint32(len(s.opts.HeartData))
	s.stopped = 0

	return nil
}

func (s *Server) Serve() error {
	if atomic.LoadInt32(&s.stopped) == 1 {
		return errors.New("server uninitialized")
	}
	s.wg.Add(1)
	go s.work()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-s.opts.Ctx.Done():
		s.Stop()
		return s.opts.Ctx.Err()
	case <-s.stopChan:
		s.wait()
	case sig := <-c:
		s.Stop()
		return errors.New("signal:" + sig.String())
	}
	return nil
}

func (s *Server) wait() {
	s.wg.Wait()
	for atomic.LoadUint32(&s.connNum) > 0 {
	}
}

func (s *Server) Stop() {
	if !atomic.CompareAndSwapInt32(&s.stopped, 0, 1) {
		return
	}
	close(s.stopChan)
	s.listener.Close()
	s.wait()
}

func (s *Server) ConnNum() uint32 {
	return atomic.LoadUint32(&s.connNum)
}

func (s *Server) onConnClose() {
	if s.limiter != nil {
		s.limiter.Revert()
	}
	atomic.AddUint32(&s.connNum, ^uint32(0))
}

// isHeartBeat called when len(data) == len(s.opts.HeartData)
func (s *Server) isHeartBeat(data []byte) bool {
	for i := 0; i < len(s.opts.HeartData); i++ {
		if s.opts.HeartData[i] != data[i] {
			return false
		}
	}
	return true
}

func (s *Server) work() {
	ctx, cancel := context.WithCancel(s.opts.Ctx)
	defer func() {
		cancel()
		s.wg.Done()
		s.Stop()
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
			s.opts.Logger.Errorf("TCP server accept error:[%v]", err)
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
