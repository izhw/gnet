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
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/izhw/gnet/gcore"
	"github.com/izhw/gnet/tcp/internal"
)

var _ gcore.Conn = &AsyncClient{}

type AsyncClient struct {
	opts      gcore.Options
	conn      net.Conn
	buffer    *internal.ReaderBuffer
	sendChan  chan []byte
	closeChan chan struct{}
	wwg       sync.WaitGroup
	rwg       sync.WaitGroup
	closed    int32
	tag       string
}

func NewAsyncClient() *AsyncClient {
	return &AsyncClient{}
}

func (c *AsyncClient) WithOptions(opts gcore.Options) {
	c.opts = opts
}

func (c *AsyncClient) Init(opts ...gcore.Option) error {
	for _, opt := range opts {
		opt(&c.opts)
	}
	conn, err := net.Dial("tcp", c.opts.Addr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.buffer = internal.NewReaderBuffer(c.conn, int(c.opts.InitReadBufLen), int(c.opts.MaxReadBufLen))
	c.sendChan = make(chan []byte, 100)
	c.closeChan = make(chan struct{})
	c.wwg.Add(1)
	if len(c.opts.HeartData) > 0 {
		go c.handleWriteLoopWithHeartbeat()
	} else {
		go c.handleWriteLoop()
	}
	c.rwg.Add(1)
	go c.handleReadLoop()
	return nil
}

func (c *AsyncClient) Read(buf []byte) (n int, err error) {
	return 0, gcore.ErrConnInvalidCall
}

func (c *AsyncClient) ReadFull(buf []byte) (n int, err error) {
	return 0, gcore.ErrConnInvalidCall
}

func (c *AsyncClient) WriteRead(req []byte) (body []byte, err error) {
	return nil, gcore.ErrConnInvalidCall
}

// Write data should be without header if Encoder != nil
func (c *AsyncClient) Write(data []byte) error {
	if len(data) > 0 {
		select {
		case <-c.closeChan:
			return gcore.ErrConnClosed
		default:
			c.sendChan <- data
		}
	}
	return nil
}

func (c *AsyncClient) Close() (err error) {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return
	}
	close(c.closeChan)
	c.wwg.Wait()
	for len(c.sendChan) > 0 {
		data := <-c.sendChan
		if err := c.write(data); err != nil {
			c.opts.Handler.OnWriteError(c, data, err)
		}
	}
	err = c.conn.Close()
	c.rwg.Wait()
	c.buffer.Release()
	c.opts.Handler.OnClosed(c)
	return
}

func (c *AsyncClient) Closed() bool {
	if atomic.LoadInt32(&c.closed) == 1 {
		return true
	}
	return false
}

func (c *AsyncClient) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *AsyncClient) SetTag(tag string) {
	c.tag = tag
}

func (c *AsyncClient) GetTag() string {
	return c.tag
}

func (c *AsyncClient) handleReadLoop() {
	defer func() {
		c.rwg.Done()
		c.Close()
	}()

	h := c.opts.Handler
	h.OnOpened(c)

	for {
		select {
		case <-c.closeChan:
			return
		default:
		}
		if _, err := c.buffer.ReadFromReader(); err != nil {
			select {
			case <-c.closeChan:
				return
			default:
			}
			if err != io.EOF {
				c.opts.Logger.Debugf("TCP client read error:[%v]", err)
			}
			return
		}
		for c.buffer.Len() > 0 {
			bodyLen, headerLen := c.opts.HeaderCodec.Decode(c.buffer.Data())
			if headerLen == 0 {
				break
			}
			msgLen := bodyLen + headerLen
			if msgLen > c.opts.MaxReadBufLen {
				c.opts.Logger.Warnf("msg len:%d greater than max:%d", msgLen, c.opts.MaxReadBufLen)
				return
			}
			if uint32(c.buffer.Len()) < msgLen {
				break
			}
			buf := make([]byte, bodyLen)
			c.buffer.Read(int(headerLen), int(bodyLen), buf)
			if err := h.OnReadMsg(c, buf); err != nil {
				c.opts.Logger.Infof("TCP client OnReadMsg error:[%v]", err)
				return
			}
		}
	}
}

func (c *AsyncClient) handleWriteLoop() {
	defer func() {
		c.wwg.Done()
		c.Close()
	}()

	for {
		select {
		case <-c.opts.Ctx.Done():
			return
		case <-c.closeChan:
			return
		case data, ok := <-c.sendChan:
			if !ok {
				return
			}
			err := c.write(data)
			if err != nil {
				c.opts.Handler.OnWriteError(c, data, err)
				return
			}
		}
	}
}

func (c *AsyncClient) handleWriteLoopWithHeartbeat() {
	timer := time.NewTimer(c.opts.HeartInterval)
	defer func() {
		timer.Stop()
		c.wwg.Done()
		c.Close()
	}()

	for {
		select {
		case <-c.opts.Ctx.Done():
			return
		case <-c.closeChan:
			return
		case data, ok := <-c.sendChan:
			if !ok {
				return
			}
			err := c.write(data)
			if err != nil {
				c.opts.Handler.OnWriteError(c, data, err)
				return
			}
			continue
		default:
		}
		// heartbeat
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(c.opts.HeartInterval)
		select {
		case <-c.opts.Ctx.Done():
			return
		case <-c.closeChan:
			return
		case data, ok := <-c.sendChan:
			if !ok {
				return
			}
			err := c.write(data)
			if err != nil {
				c.opts.Handler.OnWriteError(c, data, err)
				return
			}
		case <-timer.C:
			if err := c.write(c.opts.HeartData); err != nil {
				if err != io.EOF {
					c.opts.Logger.Infof("TCP client write heartbeat error:[%v]", err)
				}
				return
			}
		}
	}
}

func (c *AsyncClient) getWriteDeadLine() (t time.Time) {
	if c.opts.WriteTimeout > 0 {
		t = time.Now().Add(c.opts.WriteTimeout)
	}
	return
}

func (c *AsyncClient) write(data []byte) (err error) {
	data = c.opts.HeaderCodec.Encode(data)
	_ = c.conn.SetWriteDeadline(c.getWriteDeadLine())
	_, err = c.conn.Write(data)
	return
}
