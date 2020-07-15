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

package protocol

import (
	"encoding/binary"
)

// HeaderDecodeProtocol returns the integer value and its length
// It returns (0, 0) if there is a parse error
// b: [header...]
// v: the value of header, is the length of body
type HeaderDecodeProtocol func(b []byte) (v uint32, n uint32)

// DecodeFixed32 parses b as a big-endian uint32, reporting its length
func DecodeFixed32(b []byte) (v uint32, n uint32) {
	if len(b) < 4 {
		return 0, 0
	}
	v = binary.BigEndian.Uint32(b)
	return v, 4
}

// ProtoDecodeProtocol parses a protobuf varint encoded integer from b,
// returning the integer value and the length of the varint
//func DecodeProtoVarint(b []byte) (uint32, uint32) {
//	v64, n := proto.DecodeVarint(b)
//	return uint32(v64), uint32(n)
//}
