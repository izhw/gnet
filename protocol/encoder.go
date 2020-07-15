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

// HeaderEncodeProtocol returns header+data
// data: body data
type HeaderEncodeProtocol func(data []byte) []byte

// EncodeFixed32 returns header(4 bytes, big-endian uint32)+data
func EncodeFixed32(data []byte) []byte {
	b := make([]byte, 4, 4+len(data))
	binary.BigEndian.PutUint32(b, uint32(len(data)))
	b = append(b, data...)
	return b
}

// EncodeProtoVarint returns header(protobuf varint)+data
//func EncodeProtoVarint(data []byte) []byte {
//	b := proto.EncodeVarint(uint64(len(data)))
//	return append(b, data...)
//}
