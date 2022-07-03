package utils

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type Logger interface {
	Printf(string, ...interface{})
	Println(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})
}

type CustomLogger struct {
	fileLogger   *log.Logger
	stdoutLogger *log.Logger
}

func NewLogger(stdout bool, filepath string, flags int) Logger {
	logger := &CustomLogger{}

	// Colored stdout compatibility for Windows
	var output io.Writer
	if runtime.GOOS == "windows" {
		output = colorable.NewColorableStdout()
	} else {
		output = os.Stdout
	}

	// Initialize stdout logger if enabled
	if stdout {
		logger.stdoutLogger = log.New(output, "", flags)
	}

	// If filepath is "/dev/null" return, else format logfile path string
	if filepath == "/dev/null" {
		return logger
	} else if filepath == "" {
		filepath = path.Join(binDir, "logs/client.utils")
	} else if !path.IsAbs(filepath) {
		filepath = path.Join(binDir, filepath)
	}

	// Create logfile parent directories
	dir := path.Dir(filepath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		err = errors.New("error creating logfile directory: " + err.Error())
		panic(err)
	}

	// Create logfile
	logFile, err := os.Create(filepath)
	if err != nil {
		err = errors.New("error creating logfile: " + err.Error())
		panic(err)
	}

	// Initialize file logger with new logfile
	logger.fileLogger = log.New(logFile, "", flags)

	return logger
}

func (logger *CustomLogger) output(calldepth int, color func(interface{}) aurora.Value, prefix, str string) {
	if logger == nil {
		return
	}
	calldepth++

	if logger.fileLogger != nil {
		err := logger.fileLogger.Output(calldepth, prefix+str)

		if err != nil {
			err = errors.New("error logging to file: " + err.Error())
			panic(err)
		}
	}

	if logger.stdoutLogger != nil {
		if color != nil {
			prefix = color(prefix).String()
		}
		// Don't print long strings in stdout, truncate them to 400 chars.
		if len(str) > 403 {
			str = str[0:400] + "..."
		}

		err := logger.stdoutLogger.Output(calldepth, prefix+str)

		if err != nil {
			err = errors.New("error logging to stdout: " + err.Error())
			panic(err)
		}
	}
}

func (logger *CustomLogger) Printf(format string, v ...interface{}) {
	logger.output(2, aurora.Green, "INFO ", fmt.Sprintf(format, v...))
}

func (logger *CustomLogger) Println(v ...interface{}) {
	logger.output(2, aurora.Green, "INFO ", fmt.Sprintln(v...))
}

func (logger *CustomLogger) Warnf(format string, v ...interface{}) {
	logger.output(2, aurora.Red, "WARN ", fmt.Sprintf(format, v...))
}

func (logger *CustomLogger) Warnln(v ...interface{}) {
	logger.output(2, aurora.Red, "WARN ", fmt.Sprintln(v...))
}

var binDir = func() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Clean(filepath.Dir(ex) + "/")
}()
