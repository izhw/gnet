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
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/izhw/gnet"
	"github.com/izhw/gnet/tcp/internal"
)

var _ gnet.Conn = &Conn{}

type Conn struct {
	s         *server
	conn      *net.TCPConn
	buffer    *internal.ReaderBuffer
	sendChan  chan []byte
	closeChan chan struct{}
	wwg       sync.WaitGroup
	rwg       sync.WaitGroup
	closed    int32
	tag       string
}

func newConn(ctx context.Context, s *server, conn *net.TCPConn) *Conn {
	c := &Conn{
		s:         s,
		conn:      conn,
		sendChan:  make(chan []byte, 100),
		closeChan: make(chan struct{}),
	}
	c.buffer = internal.NewReaderBuffer(c.conn, int(s.opts.InitReadBufLen), int(s.opts.MaxReadBufLen))
	c.wwg.Add(1)
	go c.handleWriteLoop(ctx)
	c.rwg.Add(1)
	go c.handleReadLoop(ctx)
	return c
}

func (c *Conn) Read(buf []byte) (n int, err error) {
	return 0, gnet.ErrConnInvalidCall
}

func (c *Conn) ReadFull(buf []byte) (n int, err error) {
	return 0, gnet.ErrConnInvalidCall
}

func (c *Conn) WriteRead(req []byte) (body []byte, err error) {
	return nil, gnet.ErrConnInvalidCall
}

func (c *Conn) Write(data []byte) error {
	if len(data) > 0 {
		select {
		case <-c.closeChan:
			return gnet.ErrConnClosed
		case c.sendChan <- data:
		}
	}
	return nil
}

func (c *Conn) Close() (err error) {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return
	}
	close(c.closeChan)
	c.wwg.Wait()
	for len(c.sendChan) > 0 {
		data := <-c.sendChan
		if err := c.write(data); err != nil {
			c.s.handler.OnWriteError(c, data, err)
		}
	}
	err = c.conn.Close()
	c.rwg.Wait()
	c.buffer.Release()
	c.s.handler.OnClosed(c)
	c.s.onConnClose()
	c.s = nil
	return
}

func (c *Conn) Closed() bool {
	if atomic.LoadInt32(&c.closed) == 1 {
		return true
	}
	return false
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetTag(tag string) {
	c.tag = tag
}

func (c *Conn) GetTag() string {
	return c.tag
}

func (c *Conn) getReadDeadLine() (t time.Time) {
	if c.s.opts.ReadTimeout > 0 {
		t = time.Now().Add(c.s.opts.ReadTimeout)
	}
	return
}

func (c *Conn) getWriteDeadLine() (t time.Time) {
	if c.s.opts.WriteTimeout > 0 {
		t = time.Now().Add(c.s.opts.WriteTimeout)
	}
	return
}

func (c *Conn) handleReadLoop(ctx context.Context) {
	defer func() {
		c.rwg.Done()
		c.Close()
	}()

	h := c.s.handler
	h.OnOpened(c)

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.closeChan:
			return
		default:
		}

		if err := c.conn.SetReadDeadline(c.getReadDeadLine()); err != nil {
			c.s.opts.Logger.Warnf("TCP conn SetReadDeadline error:[%v]", err)
		}
		if _, err := c.buffer.ReadFromReader(); err != nil {
			select {
			case <-c.closeChan:
				return
			default:
			}
			if err != io.EOF {
				c.s.opts.Logger.Debugf("TCP conn read error:[%v]", err)
			}
			return
		}
		for c.buffer.Len() > 0 {
			bodyLen, headerLen := c.s.opts.HeaderCodec.Decode(c.buffer.Data())
			if headerLen == 0 {
				break
			}
			msgLen := bodyLen + headerLen
			if msgLen > c.s.opts.MaxReadBufLen {
				c.s.opts.Logger.Errorf("msg len:%d greater than max:%d", msgLen, c.s.opts.MaxReadBufLen)
				return
			}
			if uint32(c.buffer.Len()) < msgLen {
				break
			}
			buf := make([]byte, bodyLen)
			c.buffer.Read(int(headerLen), int(bodyLen), buf)
			if err := h.OnReadMsg(c, buf); err != nil {
				c.s.opts.Logger.Infof("TcpConn OnReadMsg error:[%v]", err)
				return
			}
		}
	}
}

func (c *Conn) handleWriteLoop(ctx context.Context) {
	defer func() {
		c.wwg.Done()
		c.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.closeChan:
			return
		case data, ok := <-c.sendChan:
			if !ok {
				return
			}
			if err := c.write(data); err != nil {
				c.s.handler.OnWriteError(c, data, err)
				return
			}
		}
	}
}

func (c *Conn) write(data []byte) (err error) {
	header := c.s.opts.HeaderCodec.Encode(data)
	_ = c.conn.SetWriteDeadline(c.getWriteDeadLine())
	if _, err = c.conn.Write(header); err != nil {
		return
	}
	_, err = c.conn.Write(data)
	return
}
