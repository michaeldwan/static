package push

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/michaeldwan/static/context"
	"github.com/michaeldwan/static/printer"
)

type stats struct {
	created int
	updated int
	deleted int
	skipped int
	bytes   int64
}

func Perform(ctx *context.Context) {
	printer.Infof("Source: %s\n", ctx.Config.SourceDirectory)
	printer.Infof("Destination: %s (%s)\n", ctx.Config.S3Bucket, ctx.Config.S3Region)
	fmt.Println(ctx.Config.MaxAge("index.html"))
	manifest := newManifest(ctx)
	printScanMsg(manifest, false)
	ticker := time.NewTicker(50 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			printScanMsg(manifest, false)
		}
	}()
	manifest.scan()
	ticker.Stop()
	printScanMsg(manifest, true)

	stats := &stats{}
	printer.Debugln("Parallel uploads:", ctx.Flags.Concurrency)
	if ctx.Flags.DryRun {
		printer.Infoln("*** Dry Run, operations are simulated ***")
	}

	queue := make(chan Entry, len(manifest.entries))
	for _, entry := range manifest.entries {
		queue <- *entry
	}
	close(queue)
	abort := make(chan error)
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(len(manifest.entries))

	for worker := 0; worker < ctx.Flags.Concurrency; worker++ {
		go func() {
			for {
				select {
				case entry, ok := <-queue:
					if !ok {
						return
					}
					if err := pushEntry(ctx, entry, stats); err != nil {
						abort <- err
					}
					wg.Done()
				case <-done:
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(abort)
	}()

	if err := <-abort; err != nil {
		close(done)
		printer.Infoln(err)
		os.Exit(1)
		return
	}

	printStats(*stats)
}

func pushEntry(ctx *context.Context, e Entry, stats *stats) error {
	for attempt := 1; attempt <= 3; attempt++ {
		var err error
		switch e.Operation() {
		case Create:
			if err = putFile(ctx, e.Src); err != nil {
				stats.bytes += e.Src.Size()
				stats.created++
			}
		case Update, ForceUpdate:
			if err = putFile(ctx, e.Src); err != nil {
				stats.bytes += e.Src.Size()
				stats.updated++
			}
		case Delete:
			if err = deleteKey(ctx, e.Key); err != nil {
				stats.deleted++
			}
		case Skip:
			stats.skipped++
			err = nil
		}
		printEntryStatus(e, attempt, err, true)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("Failed to push %s after 3 attempts. Aborting.\n", e.Key)
}
