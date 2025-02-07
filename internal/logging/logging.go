package logging

/* Copied & adapted from github/hashicorp/terraform */

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/hashicorp/go-hclog"
)

// These are the environmental variables that determine if we log, and if
// we log whether or not the log should go to a file.
const (
	envLog     = "PYGO_LOG"
	envLogFile = "PYGO_LOG_PATH"
)

var (
	// ValidLevels are the log level names that Terraform recognizes.
	ValidLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "OFF"}

	// logger is the global hclog logger
	logger hclog.Logger

	// logWriter is a global writer for logs, to be used with the std log package
	logWriter io.Writer

	// initialize our cache of panic output
	panics = &panicRecorder{
		panics:   make(map[string][]string),
		maxLines: 100,
	}
)

func init() {
	logger = newHCLogger("")
	logWriter = logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})

	// set up the default std library logger to use our output
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(logWriter)
}

// SetupTempLog adds a new log sink which writes all logs to the given file.
func RegisterSink(f *os.File) {
	l, ok := logger.(hclog.InterceptLogger)
	if !ok {
		panic("logger is not an InterceptLogger")
	}

	if f == nil {
		return
	}

	l.RegisterSink(hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Level:  hclog.Trace,
		Output: f,
	}))
}

// LogOutput return the default global log io.Writer
func LogOutput() io.Writer {
	return logWriter
}

// HCLogger returns the default global hclog logger
func HCLogger() hclog.Logger {
	return logger
}

// newHCLogger returns a new hclog.Logger instance with the given name
func newHCLogger(name string) hclog.Logger {
	logOutput := io.Writer(os.Stderr)
	logLevel, json := globalLogLevel()

	if logPath := os.Getenv(envLogFile); logPath != "" {
		f, err := os.OpenFile(logPath, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		} else {
			logOutput = f
		}
	}

	return hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              name,
		Level:             logLevel,
		Output:            logOutput,
		IndependentLevels: true,
		JSONFormat:        json,
	})
}

// NewLogger returns a new logger based in the current global logger, with the
// given name appended.
func NewLogger(name string) hclog.Logger {
	if name == "" {
		panic("logger name required")
	}
	return &logPanicWrapper{
		Logger: logger.Named(name),
	}
}

// CurrentLogLevel returns the current log level string based the environment vars
func CurrentLogLevel() string {
	ll, _ := globalLogLevel()
	return strings.ToUpper(ll.String())
}

func globalLogLevel() (hclog.Level, bool) {
	var json bool
	envLevel := strings.ToUpper(os.Getenv(envLog))
	if envLevel == "JSON" {
		json = true
	}
	return parseLogLevel(envLevel), json
}

func parseLogLevel(envLevel string) hclog.Level {
	if envLevel == "" {
		return hclog.Off
	}
	if envLevel == "JSON" {
		envLevel = "TRACE"
	}

	logLevel := hclog.Trace
	if isValidLogLevel(envLevel) {
		logLevel = hclog.LevelFromString(envLevel)
	} else {
		fmt.Fprintf(os.Stderr, "[WARN] Invalid log level: %q. Defaulting to level: TRACE. Valid levels are: %+v",
			envLevel, ValidLevels)
	}

	return logLevel
}

// IsDebugOrHigher returns whether or not the current log level is debug or trace
func IsDebugOrHigher() bool {
	level, _ := globalLogLevel()
	return level == hclog.Debug || level == hclog.Trace
}

func isValidLogLevel(level string) bool {
	for _, l := range ValidLevels {
		if level == string(l) {
			return true
		}
	}

	return false
}

// PluginOutputMonitor creates an io.Writer that will warn about any writes in
// the default logger. This is used to catch unexpected output from plugins,
// notifying them about the problem as well as surfacing the lost data for
// context.
func PluginOutputMonitor(source string) io.Writer {
	return pluginOutputMonitor{
		source: source,
		log:    logger,
	}
}

// pluginOutputMonitor is an io.Writer that logs all writes as
// "unexpected data" with the source name.
type pluginOutputMonitor struct {
	source string
	log    hclog.Logger
}

func (w pluginOutputMonitor) Write(d []byte) (int, error) {
	// Limit the write size to 1024 bytes We're not expecting any data to come
	// through this channel, so accidental writes will usually be stray fmt
	// debugging statements and the like, but we want to provide context to the
	// provider to indicate what the unexpected data might be.
	n := len(d)
	if n > 1024 {
		d = append(d[:1024], '.', '.', '.')
	}

	w.log.Warn("unexpected data", w.source, strings.TrimSpace(string(d)))
	return n, nil
}
