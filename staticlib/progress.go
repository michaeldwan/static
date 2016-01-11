package staticlib

import (
	"fmt"
	"strconv"
)

// Progress represents the current progress of an operation
type Progress struct {
	Num   int
	Total int
}

func (p Progress) String() string {
	if p.Total > 0 {
		percent := float64(p.Total) / float64(p.Num) * 100
		return fmt.Sprintf("%.f%% (%d/%d)", percent, p.Num, p.Total)
	}
	return strconv.Itoa(p.Num)
}

type Operation struct {
	progressCh chan Progress
	err        error
}

func newOperation() *Operation {
	return &Operation{
		progressCh: make(chan Progress),
	}
}

func (op *Operation) Progress() <-chan Progress {
	return op.progressCh
}

func (op *Operation) Err() error {
	return op.err
}

func (op *Operation) done() {
	close(op.progressCh)
}
