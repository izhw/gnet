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

package gcore

import (
	"net"
)

type Conn interface {
	// Init initiates Conn with options
	Init(opts ...Option) error
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
