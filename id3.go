// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package id3

import (
	"errors"
	"github.com/gerow/id3-go/v1"
	"github.com/gerow/id3-go/v2"
	"os"
)

// Tagger represents the metadata of a tag
type Tagger interface {
	Title() string
	Artist() string
	Album() string
	Year() string
	Genre() string
	Comments() []string
	SetTitle(string)
	SetArtist(string)
	SetAlbum(string)
	SetYear(string)
	SetGenre(string)
	AllFrames() []v2.Framer
	Frames(string) []v2.Framer
	Frame(string) v2.Framer
	DeleteFrames(string) []v2.Framer
	AddFrames(...v2.Framer)
	Bytes() []byte
	Padding() uint
	Size() int
	Version() string
}

// File represents the tagged file
type File struct {
	Tagger
	originalSize int
	file         *os.File
	readOnly     bool
}

// Opens a new tagged file
func Open(name string) (*File, error) {
	fi, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	file := &File{
		file:     fi,
		readOnly: false,
	}

	file, err = commonOpen(file, fi)

	return file, err
}

func OpenReadOnly(name string) (*File, error) {
	fi, err := os.OpenFile(name, os.O_RDONLY, 0000)
	if err != nil {
		return nil, err
	}

	file := &File{
		file:     fi,
		readOnly: true,
	}

	file, err = commonOpen(file, fi)

	return file, err
}

func commonOpen(file *File, fi *os.File) (*File, error) {
	if v2Tag := v2.ParseTag(fi); v2Tag != nil {
		file.Tagger = v2Tag
		file.originalSize = v2Tag.Size()
	} else if v1Tag := v1.ParseTag(fi); v1Tag != nil {
		file.Tagger = v1Tag
	} else {
		return nil, errors.New("Open: unknown tag format")
	}

	return file, nil
}

// Saves any edits to the tagged file
func (f *File) Close() error {
	defer f.file.Close()

	if f.readOnly {
		return nil
	}

	switch f.Tagger.(type) {
	case (*v1.Tag):
		if _, err := f.file.Seek(-v1.TagSize, os.SEEK_END); err != nil {
			return err
		}
	case (*v2.Tag):
		if f.Size() > f.originalSize {
			start := int64(f.originalSize + v2.HeaderSize)
			offset := int64(f.Tagger.Size() - f.originalSize)

			if err := shiftBytesBack(f.file, start, offset); err != nil {
				return err
			}
		}

		if _, err := f.file.Seek(0, os.SEEK_SET); err != nil {
			return err
		}
	default:
		return errors.New("Close: unknown tag version")
	}

	if _, err := f.file.Write(f.Tagger.Bytes()); err != nil {
		return err
	}

	return nil
}
