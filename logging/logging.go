package logging

import (
	"github.com/fatih/color"
	"log"
)

type Logger interface {
	Err(a ...interface{})
	Success(a ...interface{})
	Info(a ...interface{})
	Fatal(a ...interface{})
}

type LoggerWithoutFile struct {
	green func(a ...interface{}) string
	red   func(a ...interface{}) string
}

func NewLoggerWithoutFile() Logger {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	return &LoggerWithoutFile{
		green: green,
		red:   red,
	}

}
func (l *LoggerWithoutFile) Err(a ...interface{}) {
	log.Println(l.red(a...))
}
func (l *LoggerWithoutFile) Success(a ...interface{}) {
	log.Println(l.green(a...))
}
func (l *LoggerWithoutFile) Info(a ...interface{}) {
	log.Println(a...)
}
func (l *LoggerWithoutFile) Fatal(a ...interface{}) {
	log.Fatal(l.red(a...))
}
