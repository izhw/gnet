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

type Conn interface {
	Close()
	Closed() bool
	Write(data []byte) error
	RemoteAddr() string
	SetTag(tag string)
	GetTag() string
}

type ConnStatusCallback interface {
	OnClosed(c Conn)
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
