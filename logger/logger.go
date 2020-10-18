package logger

import (
	"fmt"
	"time"
)

func Log(msg string) {
	colorCyan := "\033[36m"
	colorReset := "\033[0m"

	t := time.Now()
	formated := t.Format("2006.01.02 15:04:05")
	fmt.Printf(colorCyan + "[%s]" + colorReset + " %s \n", formated, msg)
}

func LogError(msg string) {
	colorCyan := "\033[36m"
	colorReset := "\033[0m"
	colorRed := "\033[31m"

	t := time.Now()
	formated := t.Format("2006.01.02 15:04:05")
	fmt.Printf(colorCyan + "[%s]" + colorReset + colorRed + " %s \n" + colorReset, formated, msg)
}
