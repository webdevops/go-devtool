package logger

import (
	"log"
	"os"
	"fmt"
	"strings"
)

const (
	LogPrefix = ""
	prefixMain = ":: "
	prefixSub  = "   -> "
	prefixItem  = "     - "
	prefixCmd  = "      $ "
	prefixErr  = "[ERROR] "
)

type Logger struct {
	*log.Logger
}

var (
	logger *Logger
	Verbose bool
	CommandName string
)

func GetInstance(commandName string, flags int) *Logger {
	CommandName = commandName

	if logger == nil {
		logger = &Logger{log.New(os.Stdout, LogPrefix, flags)}
	}
	return logger
}

func (log Logger) Verbose(message string, sprintf ...interface{}) {
	if Verbose {
		if len(sprintf) > 0 {
			message = fmt.Sprintf(message, sprintf...)
		}

		log.Println(message)
	}
}

func (log Logger) Main(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	log.Println(prefixMain + message)
}

func (log Logger) Step(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	log.Println(prefixSub + message)
}


func (log Logger) Item(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	log.Println(prefixItem + message)
}


func (log Logger) Command(message string) {
	log.Println(prefixCmd + message)
}

func (log Logger) Printlnf(message string, sprintf ...interface{}) {
	log.Printf(message + "\n", sprintf...)
}

func (log Logger) FatalExit(exitCode int, message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	log.Fatal(message)
	os.Exit(exitCode)
}


// Log error object as message
func (log Logger) FatalErrorExit(exitCode int, err error) {

	if CommandName != "" {
		cmdline := fmt.Sprintf("%s %s", CommandName, strings.Join(os.Args[1:], " "))
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Command: %s", cmdline))
	}

	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s %s", prefixErr, err))

	os.Exit(exitCode)
}
