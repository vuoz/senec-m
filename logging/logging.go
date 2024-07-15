package logging

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

type Logger interface {
	Err(a ...interface{})
	Success(a ...interface{})
	Info(a ...interface{})
	Fatal(a ...interface{})
}
type ColorLogger struct {
	green  func(a ...interface{}) string
	red    func(a ...interface{}) string
	msgBuf bytes.Buffer
	file   *os.File
	mu     sync.Mutex
	writes int
}
type LoggerWithoutFile struct {
	green func(a ...interface{}) string
	red   func(a ...interface{}) string
}

func (l *ColorLogger) writeToFile() error {
	if l.file == nil {
		return fmt.Errorf("file does not exist")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	content := l.msgBuf.String()
	_, err := l.file.Write([]byte(content))
	if err != nil {
		return err
	}

	l.msgBuf.Reset()
	return nil
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
