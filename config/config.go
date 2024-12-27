package config

import (
	"fmt"
	"log"
	"time"
)

const TIME_FORMAT = "[2006-01-02 15:04:05]"
const DEFAULT_COLOR = 0xfff79e

type LogWriter struct{}

func (lw *LogWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("%s %s", time.Now().UTC().Format(TIME_FORMAT), string(bytes))
}

func Init() {
	log.SetFlags(0)
	log.SetOutput(new(LogWriter))
}
