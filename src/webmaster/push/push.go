package push

import (
  "fmt"
  "webmaster/context"
  "time"
)

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

  for _, entry := range manifest.entries {
    switch entry.Operation() {
    case Create:
      if err := putFile(ctx, entry.Src); err != nil {
        panic(err)
      }
      printCreateMsg(entry.Src)
    case Update:
      if err := putFile(ctx, entry.Src); err != nil {
        panic(err)
      }
      printUpdateMsg(entry.Src)
    case Delete:
      if err := deleteKey(ctx, entry.Key); err != nil {
        panic(err)
      }
      printDeleteMsg(entry.Key)
    case Skip:
      printSkipMsg(entry.Key)
    }
  }
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
