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
	"net"
)

type Factory func() (Conn, error)

// Server
// e.g. TCP server, WebSocket server
type Server interface {
	// Serve starts handling events for Server.
	Serve() error
	// Stop can stop the service whenever you want to
	// it is also called automatically when an interrupt signal arrives
	Stop()
	// ConnNum returns the number of currently active connections
	ConnNum() uint32
}

// connection pool
type Pool interface {
	// Get gets a Conn from the pool, creates an Conn if necessary,
	// removes it from the Pool, and returns it to the caller.
	Get() (conn Conn, err error)

	// Put adds conn to the pool.
	// The conn returned by Get should be passed to Put once and only once,
	// whether it's closed or not
	Put(conn Conn)

	// Close closes the pool and all connections in the pool
	Close()
}

type Conn interface {
	// Read reads data from the connection, only for sync Client.
	Read(buf []byte) (n int, err error)
	// ReadFull reads exactly len(buf) bytes from Conn into buf, only for sync Client.
	// It returns the number of bytes copied and an error if fewer bytes were read.
	// On return, n == len(buf) if and only if err == nil.
	ReadFull(buf []byte) (n int, err error)
	// WriteRead writes the request and reads the response, only for sync Client.
	// HeaderCodec(in Options) is used
	// returning msg body, without header
	WriteRead(req []byte) (body []byte, err error)

	// Write writes data to the connection.
	Write(data []byte) error
	// Close closes the connection.
	Close() error
	// Closed
	Closed() bool
	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
	// SetTag sets a tag to Conn
	SetTag(tag string)
	// GetTag gets the tag
	GetTag() string
}

// EventHandler Conn events callback
type EventHandler interface {
	// OnOpened a new Conn has been opened
	OnOpened(c Conn)
	// OnClosed c has been closed
	OnClosed(c Conn)
	// OnReadMsg read one msg
	// data: body data
	// if err != nil, conn will be closed
	OnReadMsg(c Conn, data []byte) (err error)
	// OnWriteError an error occurred while writing data to c
	OnWriteError(c Conn, data []byte, err error)
}

func DefaultEventHandler() EventHandler {
	return &NetEventHandler{}
}

// NetEventHandler is a built-in implementation for EventHandler
type NetEventHandler struct {
}

func (h *NetEventHandler) OnOpened(c Conn) {
}

func (h *NetEventHandler) OnClosed(c Conn) {
}

func (h *NetEventHandler) OnReadMsg(c Conn, data []byte) (err error) {
	return nil
}

func (h *NetEventHandler) OnWriteError(c Conn, data []byte, err error) {
}
