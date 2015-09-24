package push

import (
	"fmt"
	"sync"
	"time"

	"github.com/michaeldwan/webmaster/context"
)

type stats struct {
	created int
	updated int
	deleted int
	skipped int
	bytes   int64
}

func Perform(ctx *context.Context) {
	defer ctx.Clean()

	fmt.Printf("Source: %s\n", ctx.Config.SourceDirectory)
	fmt.Printf("Destination: %s (%s)\n", ctx.Config.S3Bucket, ctx.Config.S3Region)

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
	fmt.Println("Parallel uploads:", ctx.Flags.Concurrency)
	if ctx.Flags.DryRun {
		fmt.Println("*** Dry Run, operations are simulated ***")
	}
	queue := make(chan Entry, 500)
	var wg sync.WaitGroup
	wg.Add(len(manifest.entries))
	go func() {
		for _, entry := range manifest.entries {
			queue <- *entry
		}
	}()

	for i := 0; i < ctx.Flags.Concurrency; i++ {
		go func() {
			for entry := range queue {
				if err := pushEntry(ctx, entry, stats); err != nil {
					panic(err)
				}
				wg.Done()
			}
		}()
	}

	wg.Wait()

	printStats(*stats)
}

func pushEntry(ctx *context.Context, e Entry, stats *stats) error {
	attempt := 1
	for {
		var err error = nil
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
		if err != nil {
			attempt++
			if attempt > 3 {
				return err
			}
		}
	}
}
