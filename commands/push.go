package commands

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/michaeldwan/static/printer"
	"github.com/michaeldwan/static/staticlib"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to AWS",
	Long:  `Static will upload changed files to S3, and more..`,
	Run:   push,
}

var (
	dryRun      bool
	forceUpdate bool
	concurrency int
)

type pushStats struct {
	created int
	updated int
	deleted int
	skipped int
	bytes   int64
}

func push(cmd *cobra.Command, args []string) {
	cfg := staticlib.NewConfig(configFilePath)
	staticlib.ConfigureAWS(cfg)
	deployment := staticlib.NewDeployment(cfg)
	defer deployment.Clean()
	printer.Infof("Source: %s\n", cfg.SourceDirectory)
	printer.Infof("Destination: %s (%s)\n", cfg.S3Bucket, cfg.S3Region)

	for stats := range deployment.Compile(forceUpdate) {
		printScanMsg(stats, false)
	}
	printScanMsg(deployment.Manifest.Stats, true)

	printer.Debugln("Parallel uploads:", concurrency)
	if dryRun {
		printer.Infoln("*** Dry Run, operations are simulated ***")
	}

	pusher := staticlib.NewPush(&deployment)
	for result := range pusher.Push(concurrency, forceUpdate, dryRun) {
		printPushEntryRestult(result, verboseOutput)
	}

	if pusher.Err() != nil {
		fmt.Println("Encountered an error while pushing, abort")
		os.Exit(1)
	}

	printStats(pusher.Stats())
	pusher.Invalidate()
}

func init() {
	pushCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Shows which change(s) would be applied but doesn't perform anything.")
	pushCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "push everything")
	pushCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "parallel uploads")
	staticCmd.AddCommand(pushCmd)
}

func printScanMsg(s *staticlib.ManifestStats, done bool) {
	printer.Infof("\rScanning: %d files, %d redirects, %d existing objects", s.FileCount, s.RedirCount, s.ObjCount)
	if done {
		printer.Infof(", done\n")
	}
}

func printStats(stats staticlib.PushStats) {
	printer.Infof("Done: %d files created, %d updated, %d deleted, and %d skipped ~ %s\n", stats.Created, stats.Updated, stats.Deleted, stats.Skipped, formatByteSize(float64(stats.Bytes)))
}

func printPushEntryRestult(result staticlib.PushEntryResult, verbose bool) {
	level := printer.DefaultLevel
	if result.Entry.Operation == staticlib.Skip {
		level = printer.LevelDebug
	}
	var buffer bytes.Buffer
	buffer.WriteString("  ")
	if result.Error == nil {
		buffer.WriteString("\x1b[32m\u2713\x1b[0m ")
	} else {
		buffer.WriteString("\x1b[31m\u2717\x1b[0m ")
	}
	buffer.WriteString(sprintOperationType(*result.Entry))
	buffer.WriteString(": ")
	buffer.WriteString(sprintDesc(*result.Entry))
	buffer.WriteString(" ~ ")
	if result.Entry.Operation == staticlib.Delete {
		buffer.WriteString(formatByteSize(float64(result.Entry.Dst.Size)))
	} else {
		buffer.WriteString(formatByteSize(float64(result.Entry.Src.Size)))
	}
	if result.Entry.Src != nil && len(result.Entry.Src.Notes) > 0 {
		buffer.WriteString(" [")
		notes := strings.Join(result.Entry.Src.Notes, ", ")
		buffer.WriteString(notes)
		buffer.WriteString("]")
	}
	if err := result.Error; err != nil {
		buffer.WriteString("\n    ")
		buffer.WriteString(err.Error())
	}
	printer.Println(level, buffer.String())
}

func sprintOperationType(e staticlib.Entry) string {
	switch e.Operation {
	case staticlib.Create:
		return "Create"
	case staticlib.Update:
		return "Update"
	case staticlib.ForceUpdate:
		return "Force Update"
	case staticlib.Delete:
		return "Delete"
	default:
		return "Skip"
	}
}

func sprintDesc(e staticlib.Entry) string {
	if e.Src != nil && e.Src.IsRedirect() {
		return fmt.Sprintf("%s --> %s", e.Src.Key, e.Src.RedirectUrl)
	}
	return e.Key
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
		return fmt.Sprintf("%.1fGB", b/gB)
	case b >= mB:
		return fmt.Sprintf("%.1fMB", b/mB)
	case b >= kB:
		return fmt.Sprintf("%.1fKB", b/kB)
	}
	return fmt.Sprintf("%.1fB", b)
}
