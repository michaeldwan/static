package push

import (
  "fmt"
  "webmaster/context"
  "time"
)

func Perform(ctx *context.Context, dryRun bool) {
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

  if dryRun {
    fmt.Println("*** Dry Run, operations are simulated ***")
  }

  var created int
  var updated int
  var deleted int
  var bytes int64

  for _, entry := range manifest.entries {
    switch entry.Operation() {
    case Create:
      if !dryRun {
        if err := putFile(ctx, entry.Src); err != nil {
          panic(err)
        }
      }
      created++
      bytes += entry.Src.Size()
      printCreateMsg(entry.Src)
    case Update:
      if !dryRun {
        if err := putFile(ctx, entry.Src); err != nil {
          panic(err)
        }
      }
      updated++
      bytes += entry.Src.Size()
      printUpdateMsg(entry.Src)
    case Delete:
      if !dryRun {
        if err := deleteKey(ctx, entry.Key); err != nil {
          panic(err)
        }
      }
      deleted++
      printDeleteMsg(entry.Key)
    case Skip:
      printSkipMsg(entry.Key)
    }
  }

  fmt.Printf("\rDone: %d files created, %d updated, and %d deleted ~ %db\n", created, updated, deleted, bytes)
}

func printScanMsg(m *Manifest, done bool) {
	fmt.Printf("\rScanning: %d files, %d redirects, %d existing objects", m.fileCount, m.redirCount, m.objCount)
	if done {
		fmt.Printf(", done\n")
	}
}

func printCreateMsg(file *File) {
  if file.IsRedirect() {
    fmt.Printf(" [+] %s --> %s ~ %db\n", file.Key(), file.RedirectUrl(), file.Size())
  } else {
    fmt.Printf(" [+] %s ~ %db\n", file.Key(), file.Size())
  }
}

func printUpdateMsg(file *File) {
  if file.IsRedirect() {
    fmt.Printf(" [>] %s --> %s ~ %db\n", file.Key(), file.RedirectUrl(), file.Size())
  } else {
    fmt.Printf(" [>] %s ~ %db\n", file.Key(), file.Size())
  }
}

func printDeleteMsg(key string) {
  fmt.Printf(" [-] %s ~ delete\n", key)
}

func printSkipMsg(key string) {
  fmt.Printf(" [ ] %s ~ skip\n", key)
}
