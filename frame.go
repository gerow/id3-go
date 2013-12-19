// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

const (
	FrameHeaderSize = 10
)

type FrameType struct {
	id          string
	description string
	constructor func(FrameHead, []byte) Framer
}

type Framer interface {
	Id() string
	Size() int
	String() string
}

type FrameHead struct {
	FrameType
	statusFlags byte
	formatFlags byte
	size        int32
}

func (h FrameHead) Id() string {
	return h.id
}

func (h FrameHead) Size() int {
	return int(h.size)
}

type DataFrame struct {
	FrameHead
	Data []byte
}

func NewDataFrame(head FrameHead, data []byte) Framer {
	return &DataFrame{head, data}
}

func (f DataFrame) String() string {
	return "<binary data>"
}