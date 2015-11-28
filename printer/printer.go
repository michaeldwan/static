package printer

import (
	"fmt"
	"sync"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	DefaultLevel Level = LevelInfo
)

var (
	mu           sync.Mutex
	currentLevel Level
)

func init() {
	currentLevel = DefaultLevel
}

func SetLevel(level Level) {
	currentLevel = level
}

func Debugln(a ...interface{}) {
	Println(LevelDebug, a...)
}

func Infoln(a ...interface{}) {
	Println(LevelInfo, a...)
}

func Println(requiredLevel Level, a ...interface{}) {
	if requiredLevel < currentLevel {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(a...)
}

func Debugf(format string, a ...interface{}) {
	Printf(LevelDebug, format, a...)
}

func Infof(format string, a ...interface{}) {
	Printf(LevelInfo, format, a...)
}

func Printf(requiredLevel Level, format string, a ...interface{}) {
	if requiredLevel < currentLevel {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf(format, a...)
}
