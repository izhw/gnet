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

type ServiceType uint32

const (
	ServiceTCPServer ServiceType = 1 << iota
	ServiceTCPClient
	ServiceTCPAsyncClient
	ServiceTCPPool
	ServiceTCPAsyncPool
)

func (t ServiceType) TCPServerType() bool {
	if t&ServiceTCPServer != 0 {
		return true
	}
	return false
}

func (t ServiceType) TCPClientType() bool {
	if t&ServiceTCPClient != 0 {
		return true
	}
	return false
}

func (t ServiceType) TCPAsyncClientType() bool {
	if t&ServiceTCPAsyncClient != 0 {
		return true
	}
	return false
}

func (t ServiceType) TCPPoolType() bool {
	if t&ServiceTCPPool != 0 {
		return true
	}
	return false
}

func (t ServiceType) TCPAsyncPoolType() bool {
	if t&ServiceTCPAsyncPool != 0 {
		return true
	}
	return false
}
