// Package graylog provides a custom Logrus hook for sending logs to Graylog
// using the GELF (Graylog Extended Log Format) protocol over TCP. This package
// allows you to integrate your Go application with Graylog for centralized logging
// and monitoring. It also supports dynamic log level adjustment at runtime.
//
// Features:
// - Sends logs to Graylog using GELF over TCP.
// - Supports dynamic log level adjustment.
// - Captures and sends stack trace information for warnings and above.
// - Customizable service name for log messages.
//
// Usage:
//
// Installation:
// To use this package, first install it via `go get`:
//
//	go get github.com/hydraide/graylog
//
// Initialization:
// To initialize the Graylog logging in your application, create a new instance of the Graylog logger and set the desired log level.
//
// Example:
//
//	package main
//
//	import (
//	  "github.com/hydraide/graylog"
//	  log "github.com/sirupsen/logrus"
//	)
//
//	func main() {
//	  graylogAddress := "xxx.xxx.xxx.xxx:5140" // Replace with your Graylog server's IP address and port
//	  serviceName := "TestService"
//
//	  // Initialize Graylog
//	  g := graylog.New(graylogAddress, serviceName)
//	  // Set log level
//	  g.SetLogLevel("error")
//
//	  // Generate some logs
//	  log.Info("This is an info message")
//	  log.Warn("This is a warning message")
//	  log.Error("This is an error message")
//	}
//
// Setting Log Level:
// The `SetLogLevel` method allows you to dynamically adjust the logging level at runtime. This is particularly useful if you need to change the log level based on configuration changes or central management without restarting the application.
//
// Example usage:
//
//	g.SetLogLevel("debug")
//
// Log Levels:
// The following log levels are supported:
// - `trace`
// - `debug`
// - `info`
// - `warn`
// - `error`
// - `fatal`
// - `panic`
//
// Stack Trace:
// The package captures and includes stack trace information for log entries at the `warn` level and above. This helps in debugging by providing context about where the log entry was generated, without overwhelming the Graylog server with unnecessary information for lower log levels.
//
// Example:
// Below is a full example of setting up and using the Graylog logging library:
//
//	package main
//
//	import (
//	  "github.com/hydraide/graylog"
//	  log "github.com/sirupsen/logrus"
//	)
//
//	func main() {
//	  graylogAddress := "xxx.xxx.xxx.xxx:5140" // Replace with your Graylog server's IP address and port
//	  serviceName := "TestService"
//
//	  // Initialize Graylog
//	  g := graylog.New(graylogAddress, serviceName)
//	  // Set initial log level
//	  g.SetLogLevel("info")
//
//	  // Generate some logs
//	  log.WithFields(log.Fields{
//	    "service": "TestService",
//	    "user":    "test_user",
//	  }).Info("This is an info message")
//
//	  log.WithFields(log.Fields{
//	    "service": "TestService",
//	    "user":    "test_user",
//	  }).Warn("This is a warning message")
//
//	  log.Error("This is an error message")
//
//	  // Dynamically change log level
//	  g.SetLogLevel("error")
//
//	  log.Info("This message will not be logged because the log level is set to error")
//	  log.Error("This error message will be logged")
//	}
package graylog

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"runtime"
	"sync"
)

const (
	logFileName       = "graylog-fallback.log"
	logBackupFileName = "graylog_fallback.log.old"
)

type Graylog interface {
	// SetLogLevel sets the log level for the Graylog logger instance
	SetLogLevel(logLevel string)
	// GetHook A shutdown metódus miatt adjuk vissza
	GetHook() *Hook
}

type graylog struct {
	graylogAddress string
	serviceName    string
	graylogHook    *Hook
}

func New(graylogAddress, serviceName string) Graylog {

	g := &graylog{
		graylogAddress: graylogAddress,
		serviceName:    serviceName,
	}

	ctx, cancel := context.WithCancel(context.Background())
	g.graylogHook = NewGraylogHook(g.graylogAddress, g.serviceName, ctx, cancel)

	log.AddHook(g.graylogHook)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// default log level
	g.SetLogLevel("info")

	return g
}

func (g *graylog) GetHook() *Hook {
	return g.graylogHook
}

func (g *graylog) SetLogLevel(logLevel string) {
	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// GELFMessage represents the structure of a GELF message
type GELFMessage struct {
	Version      string                 `json:"version"`
	Host         string                 `json:"host"`
	ShortMessage string                 `json:"short_message"`
	Timestamp    float64                `json:"timestamp"`
	Level        int                    `json:"level"`
	Stack        []StackTrace           `json:"stack,omitempty"`
	Extra        map[string]interface{} `json:"-"`
}

// StackTrace represents a single frame in a stack trace
type StackTrace struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"func"`
}

// Hook is a custom hook for Logrus to send logs to Graylog
type Hook struct {
	GraylogAddress  string
	Host            string
	logQueue        chan string
	mutex           sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
	connectionError bool
}

// NewGraylogHook creates a new Graylog hook
func NewGraylogHook(address, host string, ctx context.Context, cancel context.CancelFunc) *Hook {

	hook := &Hook{
		GraylogAddress: address,
		Host:           host,
		logQueue:       make(chan string, 1000), // Buffer size
		ctx:            ctx,
		cancel:         cancel,
	}

	go hook.dispatcher()

	return hook

}

func (hook *Hook) dispatcher() {

	for {
		select {
		case <-hook.ctx.Done():
			return
		case message := <-hook.logQueue:

			err := hook.sendToGraylog(message)
			if err != nil {

				if hook.connectionError == false {
					hook.connectionError = true
				}

				hook.saveToLocalFile(message)

			} else {

				if hook.connectionError {
					hook.connectionError = false
					hook.SendFallbackLogsToGraylogServer()
				}

			}

		}
	}

}

func (hook *Hook) Shutdown() {

	// close the queue and the channel
	hook.cancel()
	close(hook.logQueue)

	hook.mutex.Lock()
	defer hook.mutex.Unlock()

}

func (hook *Hook) sendToGraylog(message string) error {

	hook.mutex.Lock()
	defer hook.mutex.Unlock()

	conn, err := net.Dial("tcp", hook.GraylogAddress)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect to Graylog server via TCP: %s", err.Error()))
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("cannot close graylog connection: %v", err)
		}
	}()

	_, err = conn.Write([]byte(message + "\n"))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to send GELF message via TCP: %s", err.Error()))
	}

	return nil

}

func (hook *Hook) saveToLocalFile(message string) {

	hook.mutex.Lock()
	defer hook.mutex.Unlock()

	// retate the local log file if it exceeds the size limit
	err := hook.rotateLocalFile()
	if err != nil {
		fmt.Printf("failed to rotate log file: %v\n", err)
		return
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open fallback log file: %v\n", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close fallback log file: %v\n", err)
		}
	}()

	// write the message to the local file
	_, writeErr := file.WriteString(message + "\n")
	if writeErr != nil {
		fmt.Printf("failed to write log to local file: %v\n", writeErr)
	}

}

func (hook *Hook) rotateLocalFile() error {

	info, err := os.Stat(logFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() < 10*1024*1024 { // 10 MB limit
		return nil
	}

	if err := os.Rename(logFileName, logBackupFileName); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	return nil
}

func (hook *Hook) SendFallbackLogsToGraylogServer() {

	logFiles := []string{logFileName, logBackupFileName}

	for _, fileName := range logFiles {

		// Megnyitjuk a log fájlt olvasásra
		file, err := os.Open(fileName)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("failed to open fallback log file %s: %v\n", fileName, err)
			return
		}

		tempFileName := fileName + ".tmp"
		tempFile, err := os.OpenFile(tempFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("failed to create temp file %s: %v\n", tempFileName, err)
			func() {
				if closeErr := file.Close(); closeErr != nil {
					fmt.Printf("failed to close log file %s: %v\n", fileName, closeErr)
				}
			}()
			return
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			message := scanner.Text()
			err := hook.sendToGraylog(message)
			if err != nil {
				// Ha hiba lép fel, írjuk vissza a megmaradt logokat a temp fájlba
				_, writeErr := tempFile.WriteString(message + "\n")
				if writeErr != nil {
					fmt.Printf("failed to write log to temp file %s: %v\n", tempFileName, writeErr)
				}
			}
		}

		// close all files
		if fileCloseErr := file.Close(); fileCloseErr != nil {
			fmt.Printf("failed to close log file %s: %v\n", fileName, fileCloseErr)
		}

		if tempCloseErr := tempFile.Close(); tempCloseErr != nil {
			fmt.Printf("failed to close temp file %s: %v\n", tempFileName, tempCloseErr)
		}

		if err := os.Rename(tempFileName, fileName); err != nil {
			fmt.Printf("failed to replace original log file %s with temp file %s: %v\n", fileName, tempFileName, err)
		} else {
			// Ha minden logot sikeresen elküldtünk, töröljük a fájlt
			if fileName != logBackupFileName {
				_ = os.Remove(fileName)
			}
		}

	}
}

// Fire is called when a log event is fired
func (hook *Hook) Fire(entry *log.Entry) error {

	// Convert logrus level to GELF level
	level := entry.Level
	gelfLevel := 7 // Default to debug
	switch level {
	case log.PanicLevel, log.FatalLevel:
		gelfLevel = 0
	case log.ErrorLevel:
		gelfLevel = 3
	case log.WarnLevel:
		gelfLevel = 4
	case log.InfoLevel:
		gelfLevel = 6
	case log.DebugLevel:
		gelfLevel = 7
	default:
		gelfLevel = 6
	}

	// Create GELF message
	gelfMessage := GELFMessage{
		Version:      "1.1",
		Host:         hook.Host,
		ShortMessage: entry.Message,
		Timestamp:    float64(entry.Time.UnixNano()) / 1e9,
		Level:        gelfLevel,
		Extra:        entry.Data, // Add additional fields from log entry
	}

	// Capture stack trace only for WARN and above
	if level <= log.WarnLevel {
		gelfMessage.Stack = captureStackTrace()
	}

	// Serialize to JSON with extra fields
	extraFields, err := json.Marshal(gelfMessage.Extra)
	if err != nil {
		return fmt.Errorf("failed to marshal extra fields: %w", err)
	}

	message, err := json.Marshal(gelfMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal GELF message: %w", err)
	}

	// Merge GELF message and extra fields
	finalMessage := fmt.Sprintf(`%s,%s}`, message[:len(message)-1], extraFields[1:])

	// Non-blocking send to logQueue
	select {
	case hook.logQueue <- finalMessage:
		// Successfully enqueued
	default:
		// Queue is full or closed
		fmt.Printf("log queue is full or closed, dropping log: %v", entry.Message)
	}

	return nil
}

// Levels defines on which log levels this hook would trigger
func (hook *Hook) Levels() []log.Level {
	return log.AllLevels
}

// captureStackTrace captures the current stack trace
func captureStackTrace() []StackTrace {
	stack := []StackTrace{}
	callers := make([]uintptr, 10)
	n := runtime.Callers(3, callers)
	frames := runtime.CallersFrames(callers[:n])
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		stack = append(stack, StackTrace{
			File: frame.File,
			Line: frame.Line,
			Func: frame.Function,
		})
	}
	return stack
}
