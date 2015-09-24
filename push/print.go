package push

import (
	"bytes"
	"fmt"
)

func printScanMsg(m *Manifest, done bool) {
	fmt.Printf("\rScanning: %d files, %d redirects, %d existing objects", m.fileCount, m.redirCount, m.objCount)
	if done {
		fmt.Printf(", done\n")
	}
}

func printEntryStatus(e Entry, attemptNumber int, err error, verbose bool) {
	if e.Operation() == Skip && !verbose {
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("  ")
	if err == nil {
		buffer.WriteString("\x1b[32m\u2713\x1b[0m ")
	} else {
		buffer.WriteString("\x1b[31m\u2717\x1b[0m ")
	}
	buffer.WriteString(sprintOperationType(e))
	buffer.WriteString(": ")
	buffer.WriteString(sprintDesc(e))
	buffer.WriteString(" ~ ")
	if e.Operation() == Delete {
		buffer.WriteString(formatByteSize(float64(e.Dst.Size)))
	} else {
		buffer.WriteString(formatByteSize(float64(e.Src.Size())))
	}
	if attemptNumber > 1 {
		buffer.WriteString(fmt.Sprintf(" (attempt #%d)", attemptNumber))
	}
	if err != nil {
		buffer.WriteString("\n    ")
		buffer.WriteString(err.Error())
	}
	fmt.Println(buffer.String())
}

func sprintOperationType(e Entry) string {
	switch e.Operation() {
	case Create:
		return "Create"
	case Update:
		return "Update"
	case ForceUpdate:
		return "Force Update"
	case Delete:
		return "Delete"
	default:
		return "Skip"
	}
}

func sprintDesc(e Entry) string {
	if e.Src != nil && e.Src.IsRedirect() {
		return fmt.Sprintf("%s --> %s", e.Src.Key(), e.Src.RedirectUrl())
	}
	return e.Key
}

func printStats(stats stats) {
	fmt.Printf("\rDone: %d files created, %d updated, and %d deleted ~ %s\n", stats.created, stats.updated, stats.deleted, formatByteSize(float64(stats.bytes)))
}

const (
	_          = iota
	kB float64 = 1 << (10 * iota)
	mB
	gB
)

func formatByteSize(b float64) string {
	switch {
	case b >= gB:
		return fmt.Sprintf("%.2fGB", b/gB)
	case b >= mB:
		return fmt.Sprintf("%.2fMB", b/mB)
	case b >= kB:
		return fmt.Sprintf("%.2fKB", b/kB)
	}
	return fmt.Sprintf("%.2fB", b)
}
