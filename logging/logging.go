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
	Log(level rune, inp ...interface{})
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

// this will panic if you give an invalid file name
func NewColorLoggerWithFile(name string) Logger {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Error opening file for logging")

	}
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	buf := bytes.Buffer{}
	return &ColorLogger{
		green:  green,
		red:    red,
		file:   file,
		mu:     sync.Mutex{},
		msgBuf: buf,
	}

}
func NewLoggerWithoutFile() Logger {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	return &LoggerWithoutFile{
		green: green,
		red:   red,
	}

}
func (l *ColorLogger) Log(level rune, inp ...interface{}) {
	switch level {
	case 'E':
		{
			go func() {
				l.mu.Lock()
				fmt.Fprint(&l.msgBuf, inp...)
				fmt.Fprintf(&l.msgBuf, "\n")
				l.mu.Unlock()
				if l.writes >= 10 {
					l.writeToFile()
					l.writes = 0
				} else {

					l.writes = l.writes + 1
				}

			}()

			log.Println(l.red(inp...))

		}
	case 'I':
		{

			log.Println(inp...)

		}
	case 'S':
		{

			log.Println(l.green(inp...))

		}
	case 'P':
		{
			l.mu.Lock()
			fmt.Fprint(&l.msgBuf, inp...)
			fmt.Fprint(&l.msgBuf, "\n")
			l.mu.Unlock()
			if l.writes >= 10 {
				l.writeToFile()
				l.writes = 0
			} else {
				l.writes = l.writes + 1
			}
			log.Fatal(l.red(inp...))
		}
	default:
		{

			log.Println(inp...)

		}

	}
}
func (l *LoggerWithoutFile) Log(level rune, inp ...interface{}) {
	switch level {
	case 'E':
		{
			log.Println(l.red(inp...))

		}
	case 'I':
		{

			log.Println(inp...)

		}
	case 'S':
		{

			log.Println(l.green(inp...))

		}
	case 'P':
		{
			log.Fatal(l.red(inp...))
		}
	default:
		{
			log.Println(inp...)
		}

	}
}
