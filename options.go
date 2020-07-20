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
	"context"
	"time"

	"github.com/izhw/gnet/logger"
	"github.com/izhw/gnet/protocol"
)

type Option func(o *Options)

type Options struct {
	Logger         logger.Logger        // default: &discardLogger{}
	HeaderCodec    protocol.HeaderCodec // default: &protocol.CodecFixed32{}
	ReadTimeout    time.Duration        // default: 2m, zero value means I/O operations will not time out
	WriteTimeout   time.Duration        // default: 5s, zero value means I/O operations will not time out
	InitReadBufLen uint32               // default: 1024, init length of conn reading buf
	MaxReadBufLen  uint32               // default: network.MaxRWLen
	ConnLimit      uint32               // default: 0, unlimited, limit of conn num for Server

	// HeartData heartbeat data, for asyncClient or Pool
	// should be without header if Encoder not nil
	HeartData []byte
	// HeartInterval heartbeat interval, default: 30s
	HeartInterval time.Duration

	// Context specifies a context for the service.
	// Can be used to signal shutdown of the service.
	Ctx context.Context

	// Tag tag for gnet.Conn
	Tag string

	// PoolInitSize number of connections to establish when creating a pool
	PoolInitSize uint32
	// PoolMaxSize max number of connections in pool
	PoolMaxSize uint32
	// PoolGetTimeout timeout for getting a Conn from pool
	PoolGetTimeout time.Duration
	// PoolIdleTimeout Conn max idle duration,
	// if timeout occurs, Conn will be closed and removed from the pool
	// default:0, means conn will not time out.
	PoolIdleTimeout time.Duration
}

func DefaultOptions() Options {
	return Options{
		Logger:         logger.DefaultLogger(),
		HeaderCodec:    &protocol.CodecFixed32{},
		ReadTimeout:    2 * time.Minute,
		WriteTimeout:   5 * time.Second,
		InitReadBufLen: 1024,
		MaxReadBufLen:  MaxRWLen,
		ConnLimit:      0,
		HeartData:      nil,
		HeartInterval:  30 * time.Second,
		PoolInitSize:   0,
		PoolMaxSize:    10,
		PoolGetTimeout: 10 * time.Second,
	}
}

// WithOptions
func WithOptions(opts Options) Option {
	return func(o *Options) {
		*o = opts
	}
}

// default: discardLogger
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// default: &protocol.CodecFixed32{}
func WithHeaderCodec(codec protocol.HeaderCodec) Option {
	return func(o *Options) {
		if codec != nil {
			o.HeaderCodec = codec
		}
	}
}

// timeout: A zero value for t means I/O operations will not time out.
// default: 2m
func WithReadTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = timeout
	}
}

// timeout: A zero value for t means I/O operations will not time out.
// default: 5s
func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = timeout
	}
}

// default: init 1024, max constant.MaxRWLen
func WithBufferLen(init, max uint32) Option {
	return func(o *Options) {
		if max > 0 {
			o.MaxReadBufLen = max
		}
		if init > o.MaxReadBufLen {
			init = o.MaxReadBufLen
		}
		if init > 0 {
			o.InitReadBufLen = init
		}
	}
}

// WithConnNumLimit limit of conn for Server
// default: 0, unlimited
func WithConnNumLimit(limit uint32) Option {
	return func(o *Options) {
		o.ConnLimit = limit
	}
}

// WithHeartbeat for AsyncClient or Pool
// data: body data
func WithHeartbeat(data []byte, interval time.Duration) Option {
	return func(o *Options) {
		o.HeartData = data
		if interval > 0 {
			o.HeartInterval = interval
		}
	}
}

// WithContext
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Ctx = ctx
	}
}

// WithTag
func WithTag(tag string) Option {
	return func(o *Options) {
		o.Tag = tag
	}
}

// WithPoolSize
func WithPoolSize(init, max uint32) Option {
	return func(o *Options) {
		if max > 0 {
			o.PoolMaxSize = max
		}
		if init > o.PoolMaxSize {
			init = o.PoolMaxSize
		}
		if init > 0 {
			o.PoolInitSize = init
		}
	}
}

// WithPoolGetTimeout Get a Conn from pool, failed when timeout occurs
// default: 10s
func WithPoolGetTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.PoolGetTimeout = timeout
	}
}

// WithPoolIdleTimeout
func WithPoolIdleTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.PoolIdleTimeout = timeout
	}
}
