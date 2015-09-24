package main

import "fmt"

type Logger struct {
	// log     io.WriteCloser
	verbose bool
}

func newLogger(verbose bool) Logger {
	l := Logger{
		verbose: verbose,
		// log:     newLogFile(),
	}
	return l
}

// func newLogFile() *io.File {
// 	f, err := ioutil.TempFile("", "webmaster-log")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return f
// }

func (l Logger) Verbose(msg string) {
	if l.verbose {
		fmt.Println(msg)
	}
	// log.WriteString(msg)
}

func (l Logger) Error(msg string) {
	fmt.Println(msg)
	// log.writeString(msg)
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
	// log.WriteString(msg)
}

func (l Logger) Close() {
	// l.log.Close()
}
