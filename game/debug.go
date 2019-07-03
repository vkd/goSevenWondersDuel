package game

import (
	"log"
	"os"
)

var (
	gameLog = log.New(os.Stdout, "", 0)
	isDebug bool
)

func debug(format string, args ...interface{}) {
	if isDebug {
		gameLog.Printf(format, args...)
	}
}
