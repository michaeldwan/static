package push

import (
	"bytes"
)

type Operation int

const (
	Create Operation = iota
	Update
	Delete
	Skip
)

type Entry struct {
	Key string
	Src *File
	Dst *DestObject
}

func (e *Entry) Operation() Operation {
	if e.Src != nil && e.Dst == nil {
		return Create
	} else if e.Src == nil && e.Dst != nil {
		return Delete
	} else if e.Src != nil && e.Dst != nil && !bytes.Equal(e.Src.Digest(), e.Dst.Digest) {
		return Update
	} else {
		return Skip
	}
}
