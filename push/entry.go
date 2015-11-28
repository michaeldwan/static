package push

import (
	"bytes"

	"github.com/michaeldwan/static/context"
)

type Operation int

const (
	Skip Operation = iota
	Create
	Update
	ForceUpdate
	Delete
)

type Entry struct {
	Key       string
	Src       *File
	Dst       *DestObject
	operation Operation
}

func (e *Entry) Operation() Operation { return e.operation }

func (e *Entry) plan(ctx *context.Context) {
	if e.Src != nil && e.Dst == nil {
		e.operation = Create
	} else if e.Src == nil && e.Dst != nil {
		e.operation = Delete
	} else if e.Src != nil && e.Dst != nil && !bytes.Equal(e.Src.Digest(), e.Dst.Digest) {
		e.operation = Update
	} else if ctx.Flags.Force {
		e.operation = ForceUpdate
	} else {
		e.operation = Skip
	}
}
