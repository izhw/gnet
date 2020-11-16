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

package delay

import (
	"time"
)

type Delay interface {
	GetDelay() time.Duration
	Reset()
}

// double delay in [min, max]
type tempDelay struct {
	d   time.Duration // d *= 2
	min time.Duration // default 5ms
	max time.Duration // default 1s
}

func NewTempDelay(min, max time.Duration) Delay {
	if min == 0 {
		min = 5 * time.Millisecond
	}
	if max == 0 {
		max = time.Second
	}
	if min > max {
		min = max
	}
	return &tempDelay{
		d:   0,
		min: min,
		max: max,
	}
}

func (d *tempDelay) GetDelay() time.Duration {
	switch d.d {
	case 0:
		d.d = d.min
	case d.max:
	default:
		d.d <<= 1
		if d.d > d.max {
			d.d = d.max
		}
	}
	return d.d
}

func (d *tempDelay) Reset() {
	d.d = 0
}
