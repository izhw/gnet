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

package tcpclient

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/tcp/internal"
)

type client struct {
	opts   gnet.Options
	conn   net.Conn
	buffer *internal.ReaderBuffer
	closed int32
	tag    string
}

func NewClient(addr string, opts ...gnet.Option) (gnet.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &client{
		opts: gnet.DefaultOptions(),
		conn: conn,
	}
	for _, o := range opts {
		o(&c.opts)
	}
	c.buffer = internal.NewReaderBuffer(c.conn, int(c.opts.InitReadBufLen), int(c.opts.MaxReadBufLen))
	return c, nil
}

// Read
func (c *client) Read(buf []byte) (n int, err error) {
	_ = c.conn.SetReadDeadline(c.getReadDeadLine())
	return c.conn.Read(buf)
}

// ReadFull
// On return, n == len(buf) if and only if err == nil.
func (c *client) ReadFull(buf []byte) (n int, err error) {
	_ = c.conn.SetReadDeadline(c.getReadDeadLine())
	return io.ReadFull(c.conn, buf)
}

// WriteRead using HeaderCodec
// returning msg body, without header
func (c *client) WriteRead(data []byte) (body []byte, err error) {
	header := c.opts.HeaderCodec.Encode(data)
	_ = c.conn.SetWriteDeadline(c.getWriteDeadLine())
	if _, err := c.conn.Write(header); err != nil {
		return nil, fmt.Errorf("write:%w", err)
	}
	if _, err := c.conn.Write(data); err != nil {
		return nil, fmt.Errorf("write:%w", err)
	}

	_ = c.conn.SetReadDeadline(c.getReadDeadLine())
	for {
		if _, err := c.buffer.ReadFromReader(); err != nil {
			return nil, fmt.Errorf("read:%w", err)
		}
		bodyLen, headerLen := c.opts.HeaderCodec.Decode(c.buffer.Data())
		if headerLen == 0 {
			continue
		}
		msgLen := bodyLen + headerLen
		if msgLen > c.opts.MaxReadBufLen {
			return nil, gnet.ErrTooLarge
		}
		if uint32(c.buffer.Len()) < msgLen {
			continue
		}
		buf := make([]byte, bodyLen)
		c.buffer.Read(int(headerLen), int(bodyLen), buf)
		return buf, nil
	}
}

// Write using HeaderCodec
func (c *client) Write(data []byte) error {
	header := c.opts.HeaderCodec.Encode(data)
	_ = c.conn.SetWriteDeadline(c.getWriteDeadLine())
	if _, err := c.conn.Write(header); err != nil {
		return err
	}
	if _, err := c.conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *client) Close() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}
	return c.conn.Close()
}

func (c *client) Closed() bool {
	if atomic.LoadInt32(&c.closed) == 1 {
		return true
	}
	return false
}

func (c *client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *client) SetTag(tag string) {
	c.tag = tag
}

func (c *client) GetTag() string {
	return c.tag
}

func (c *client) getReadDeadLine() (t time.Time) {
	if c.opts.ReadTimeout > 0 {
		t = time.Now().Add(c.opts.ReadTimeout)
	}
	return
}

func (c *client) getWriteDeadLine() (t time.Time) {
	if c.opts.WriteTimeout > 0 {
		t = time.Now().Add(c.opts.WriteTimeout)
	}
	return
}
