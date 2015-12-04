package push

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/michaeldwan/static/printer"
)

func printScanMsg(m *Manifest, done bool) {
	printer.Infof("\rScanning: %d files, %d redirects, %d existing objects", m.fileCount, m.redirCount, m.objCount)
	if done {
		printer.Infof(", done\n")
	}
}

func printEntryStatus(e Entry, attemptNumber int, err error, verbose bool) {
	level := printer.DefaultLevel
	if e.Operation() == Skip {
		level = printer.LevelDebug
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
	printer.Println(level, buffer.String())
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
	printer.Infof("Done: %d files created, %d updated, %d deleted, and %d skipped ~ %s\n", stats.created, stats.updated, stats.deleted, stats.skipped, formatByteSize(float64(stats.bytes)))
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

func printAWSError(err error) {
	if err == nil {
		return
	}

	if awsErr, ok := err.(awserr.Error); ok {
		// Generic AWS error with Code, Message, and original error (if any)
		fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			// A service error occurred
			fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
		}
	} else {
		// This case should never be hit, the SDK should always return an
		// error which satisfies the awserr.Error interface.
		fmt.Println(err.Error())
	}
}
