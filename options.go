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
	"time"

	"github.com/izhw/gnet/logger"
	"github.com/izhw/gnet/protocol"
)

type Option func(o *Options)

type Options struct {
	Logger         logger.Logger                 // default: &discardLogger{}
	Encoder        protocol.HeaderEncodeProtocol // default: protocol.EncodeFixed32, nil value means Encoder disabled
	Decoder        protocol.HeaderDecodeProtocol // default: protocol.DecodeFixed32
	ReadTimeout    time.Duration                 // default: 2m, zero value means I/O operations will not time out
	WriteTimeout   time.Duration                 // default: 5s, zero value means I/O operations will not time out
	InitReadBufLen uint32                        // default: 1024, init length of conn reading buf
	MaxReadBufLen  uint32                        // default: network.MaxRWLen
	ConnLimit      uint32                        // default: 0, unlimited, limit of conn num for Server

	HeartData     []byte        // AsyncClient, heartbeat data, should be without header if Encoder not nil
	HeartInterval time.Duration // AsyncClient, heartbeat interval, default: 30s

	StatusCallback ConnStatusCallback
}

func DefaultOptions() Options {
	return Options{
		Logger:         logger.DefaultLogger(),
		Encoder:        protocol.EncodeFixed32,
		Decoder:        protocol.DecodeFixed32,
		ReadTimeout:    2 * time.Minute,
		WriteTimeout:   5 * time.Second,
		InitReadBufLen: 1024,
		MaxReadBufLen:  MaxRWLen,
		ConnLimit:      0,
		HeartData:      nil,
		HeartInterval:  30 * time.Second,
		StatusCallback: nil,
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

// default: protocol.EncodeFixed32 & protocol.DecodeFixed32
func WithParseHeaderProtocol(e protocol.HeaderEncodeProtocol, d protocol.HeaderDecodeProtocol) Option {
	return func(o *Options) {
		o.Encoder = e
		o.Decoder = d
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

// WithHeartbeat for AsyncClient
// data: body data
func WithHeartbeat(data []byte, interval time.Duration) Option {
	return func(o *Options) {
		o.HeartData = data
		if interval > 0 {
			o.HeartInterval = interval
		}
	}
}

// WithConnStatusCallback
func WithConnStatusCallback(cb ConnStatusCallback) Option {
	return func(o *Options) {
		o.StatusCallback = cb
	}
}
