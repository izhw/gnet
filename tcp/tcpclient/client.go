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

type Client struct {
	opts   gnet.Options
	conn   net.Conn
	buffer *internal.ReaderBuffer
	closed int32
	tag    string
}

func NewClient(addr string, opts ...gnet.Option) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		opts: gnet.DefaultOptions(),
		conn: conn,
	}
	for _, o := range opts {
		o(&c.opts)
	}
	c.buffer = internal.NewReaderBuffer(c.conn, int(c.opts.InitReadBufLen), int(c.opts.MaxReadBufLen))
	return c, nil
}

func (c *Client) Close() {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return
	}
	_ = c.conn.Close()
	if c.opts.StatusCallback != nil {
		c.opts.StatusCallback.OnClosed(c)
	}
}

func (c *Client) Closed() bool {
	if atomic.LoadInt32(&c.closed) == 1 {
		return true
	}
	return false
}

func (c *Client) SetTag(tag string) {
	c.tag = tag
}

func (c *Client) GetTag() string {
	return c.tag
}

func (c *Client) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Client) getReadDeadLine() (t time.Time) {
	if c.opts.ReadTimeout > 0 {
		t = time.Now().Add(c.opts.ReadTimeout)
	}
	return
}

func (c *Client) getWriteDeadLine() (t time.Time) {
	if c.opts.WriteTimeout > 0 {
		t = time.Now().Add(c.opts.WriteTimeout)
	}
	return
}

// Read not using Decoder
func (c *Client) Read(buf []byte) (n int, err error) {
	_ = c.conn.SetReadDeadline(c.getReadDeadLine())
	return c.conn.Read(buf)
}

// ReadFull not using Decoder
// On return, n == len(buf) if and only if err == nil.
func (c *Client) ReadFull(buf []byte) (n int, err error) {
	_ = c.conn.SetReadDeadline(c.getReadDeadLine())
	return io.ReadFull(c.conn, buf)
}

// Write data should be without header if Encoder != nil
func (c *Client) Write(data []byte) error {
	if c.opts.Encoder != nil {
		data = c.opts.Encoder(data)
	}
	_ = c.conn.SetWriteDeadline(c.getWriteDeadLine())
	if _, err := c.conn.Write(data); err != nil {
		return err
	}
	return nil
}

// WriteRead using Encoder(if Encoder != nil) and Decoder
// returning msg body, without header
func (c *Client) WriteRead(req []byte) (body []byte, err error) {
	if c.opts.Encoder != nil {
		req = c.opts.Encoder(req)
	}
	_ = c.conn.SetWriteDeadline(c.getWriteDeadLine())
	if _, err := c.conn.Write(req); err != nil {
		return nil, fmt.Errorf("write:%w", err)
	}

	_ = c.conn.SetReadDeadline(c.getReadDeadLine())
	for {
		if _, err := c.buffer.ReadFromReader(); err != nil {
			return nil, fmt.Errorf("read:%w", err)
		}
		bodyLen, headerLen := c.opts.Decoder(c.buffer.Data())
		if headerLen == 0 {
			continue
		}
		msgLen := bodyLen + headerLen
		if msgLen > c.opts.MaxReadBufLen {
			return nil, gnet.ErrMsgTooLarge
		}
		if uint32(c.buffer.Len()) < msgLen {
			continue
		}
		buf := make([]byte, bodyLen)
		c.buffer.Read(int(headerLen), int(bodyLen), buf)
		return buf, nil
	}
}
