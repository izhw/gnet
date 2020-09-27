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

package gnet

import (
	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/pool"
	"github.com/izhw/gnet/tcp/tcpclient"
	"github.com/izhw/gnet/tcp/tcpserver"
)

type service struct {
	opts   gcore.Options
	server Server
	client Conn
	pool   Pool
}

func newService(opts ...gcore.Option) Service {
	s := &service{
		opts: gcore.DefaultOptions(),
	}
	for _, opt := range opts {
		opt(&s.opts)
	}
	s.init()
	return s
}

func (s *service) init() {
	if s.opts.ServiceType.TCPServerType() {
		svr := tcpserver.NewServer()
		svr.WithOptions(s.opts)
		s.server = svr
	}
	if s.opts.ServiceType.TCPClientType() {
		c := tcpclient.NewClient()
		c.WithOptions(s.opts)
		s.client = c
	}
	if s.client == nil && s.opts.ServiceType.TCPAsyncClientType() {
		c := tcpclient.NewAsyncClient()
		c.WithOptions(s.opts)
		s.client = c
	}
	if s.opts.ServiceType.TCPPoolType() {
		p := pool.NewPool()
		p.WithOptions(s.opts)
		s.pool = p
	}
	if s.pool == nil && s.opts.ServiceType.TCPAsyncPoolType() {
		p := pool.NewAsyncPool()
		p.WithOptions(s.opts)
		s.pool = p
	}
}

// Server returns the server
func (s *service) Server() Server {
	return s.server
}

// Client returns the client
func (s *service) Client() Conn {
	return s.client
}

// Server returns the server
func (s *service) Pool() Pool {
	return s.pool
}
