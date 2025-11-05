package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// MCode represents a message code with predefined messages
type MCode struct {
	Code    string
	Message string
}

// FormatWithOptional formats the message with optional additional message
func (m MCode) FormatWithOptional(optionalMessage string) string {
	if optionalMessage == "" {
		return m.Message
	}
	return fmt.Sprintf("%s: %s", m.Message, optionalMessage)
}

// Message Codes for goxcel
var (
	// System Layer Codes - SY_* (SYstem)
	SYS1 = MCode{"SY-S1", "Application started"}
	SYS2 = MCode{"SY-S2", "Application terminated successfully"}
	SYS3 = MCode{"SY-S3", "Application terminated with error"}
	SYE1 = MCode{"SY-E1", "Unexpected error occurred"}

	// File System Operation Codes - FS_* (File System)
	FSO1 = MCode{"FS-O1", "File system operation success"}
	FSO2 = MCode{"FS-O2", "File system operation failed"}
	FSM1 = MCode{"FS-M1", "Directory creation success"}
	FSM2 = MCode{"FS-M2", "Directory creation failed"}
	FSW1 = MCode{"FS-W1", "File write success"}
	FSW2 = MCode{"FS-W2", "File write failed"}
	FSR1 = MCode{"FS-R1", "File read success"}
	FSR2 = MCode{"FS-R2", "File read failed"}

	// Repository Layer Codes - R_* (Repository)
	RI1 = MCode{"R-I1", "Repository initialization success"}
	RI2 = MCode{"R-I2", "Repository initialization failed"}
	RP1 = MCode{"R-P1", "Parse operation success"}
	RP2 = MCode{"R-P2", "Parse operation failed"}
	RW1 = MCode{"R-W1", "Write operation success"}
	RW2 = MCode{"R-W2", "Write operation failed"}

	// Controller Layer Codes - C_* (Controller)
	CI1 = MCode{"C-I1", "Controller initialization success"}
	CI2 = MCode{"C-I2", "Controller initialization failed"}
	CC1 = MCode{"C-C1", "Command execution success"}
	CC2 = MCode{"C-C2", "Command execution failed"}

	// UseCase Layer Codes - U_* (UseCase)
	UI1 = MCode{"U-I1", "UseCase initialization success"}
	UI2 = MCode{"U-I2", "UseCase initialization failed"}
	UR1 = MCode{"U-R1", "Rendering operation success"}
	UR2 = MCode{"U-R2", "Rendering operation failed"}

	// UseCase Cell Layer Codes - UC_* (UseCase Cell)
	UCE1 = MCode{"UC-E1", "Cell expression expansion started"}
	UCE2 = MCode{"UC-E2", "Cell expression expansion completed"}
	UCT1 = MCode{"UC-T1", "Cell type inference"}
	UCR1 = MCode{"UC-R1", "Cell path resolution"}
	UCS1 = MCode{"UC-S1", "Cell style parsing"}

	// UseCase Sheet Layer Codes - US_* (UseCase Sheet)
	USR1 = MCode{"US-R1", "Sheet rendering started"}
	USR2 = MCode{"US-R2", "Sheet rendering completed"}
	USG1 = MCode{"US-G1", "Grid rendering"}
	USF1 = MCode{"US-F1", "For loop processing"}
	USA1 = MCode{"US-A1", "Anchor positioning"}

	// UseCase Book Layer Codes - UB_* (UseCase Book)
	UBR1 = MCode{"UB-R1", "Book rendering started"}
	UBR2 = MCode{"UB-R2", "Book rendering completed"}
	UBN1 = MCode{"UB-N1", "Data normalization"}

	// Model Layer Codes - M_* (Model)
	MV1 = MCode{"M-V1", "Model validation success"}
	MV2 = MCode{"M-V2", "Model validation failed"}
	MC1 = MCode{"M-C1", "Model conversion success"}
	MC2 = MCode{"M-C2", "Model conversion failed"}

	// GXL Processing Codes - GXL_*
	GXLP1 = MCode{"GXL-P1", "GXL parsing success"}
	GXLP2 = MCode{"GXL-P2", "GXL parsing failed"}
	GXLR1 = MCode{"GXL-R1", "GXL rendering success"}
	GXLR2 = MCode{"GXL-R2", "GXL rendering failed"}

	// XLSX Processing Codes - XLSX_*
	XLSXW1 = MCode{"XLSX-W1", "XLSX generation success"}
	XLSXW2 = MCode{"XLSX-W2", "XLSX generation failed"}
	XLSXR1 = MCode{"XLSX-R1", "XLSX reading success"}
	XLSXR2 = MCode{"XLSX-R2", "XLSX reading failed"}

	// XML Processing Codes - XML_*
	XMLM1 = MCode{"XML-M1", "XML marshaling success"}
	XMLM2 = MCode{"XML-M2", "XML marshaling failed"}
	XMLU1 = MCode{"XML-U1", "XML unmarshaling success"}
	XMLU2 = MCode{"XML-U2", "XML unmarshaling failed"}
)

// LogLevel represents the log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Code      string                 `json:"code"`
	Component string                 `json:"component"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	File      string                 `json:"file,omitempty"`
	Function  string                 `json:"function,omitempty"`
	Line      int                    `json:"line,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Component    string `json:"component" yaml:"component"`
	Service      string `json:"service" yaml:"service"`
	Level        string `json:"level" yaml:"level"`
	Structured   bool   `json:"structured" yaml:"structured"`
	EnableCaller bool   `json:"enable_caller" yaml:"enable_caller"`
	Output       string `json:"output" yaml:"output"`
}

// LOGGER represents the application logger
type LOGGER struct {
	config *LoggerConfig
	level  LogLevel
	output io.Writer
}

// Logger defines the logging interface
type Logger interface {
	DEBUG(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	INFO(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	WARN(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	ERROR(mcode MCode, optionalMessage string, fields ...map[string]interface{})
	FATAL(mcode MCode, optionalMessage string, fields ...map[string]interface{})
}

// NewLogger creates a new logger instance
func NewLogger(loggerConfig LoggerConfig) Logger {
	logger := &LOGGER{
		config: &loggerConfig,
		output: os.Stdout,
	}

	// Set log level
	switch strings.ToUpper(loggerConfig.Level) {
	case "DEBUG":
		logger.level = DEBUG
	case "INFO":
		logger.level = INFO
	case "WARN":
		logger.level = WARN
	case "ERROR":
		logger.level = ERROR
	case "FATAL":
		logger.level = FATAL
	default:
		logger.level = INFO
	}

	// Set output
	switch loggerConfig.Output {
	case "stderr":
		logger.output = os.Stderr
	case "stdout", "":
		logger.output = os.Stdout
	default:
		// File output
		if file, err := os.OpenFile(loggerConfig.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666); err == nil {
			logger.output = file
		} else {
			logger.output = os.Stdout
			logger.ERROR(FSW2, fmt.Sprintf("file: %s, error: %s", loggerConfig.Output, err.Error()))
		}
	}

	return logger
}

// log writes a log entry using MCode
func (rcv *LOGGER) log(level LogLevel, mcode MCode, optionalMessage string, fields map[string]interface{}) {
	if level < rcv.level {
		return
	}

	finalMessage := mcode.FormatWithOptional(optionalMessage)

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level.String(),
		Code:      mcode.Code,
		Component: rcv.config.Component,
		Service:   rcv.config.Service,
		Message:   finalMessage,
		Fields:    fields,
	}

	rcv.writeLogEntry(entry)
}

// writeLogEntry writes the actual log entry to output
func (rcv *LOGGER) writeLogEntry(entry LogEntry) {
	// Add caller information if enabled
	if rcv.config.EnableCaller {
		if pc, file, line, ok := runtime.Caller(4); ok {
			entry.File = file
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}

	// Extract common fields from fields map
	if entry.Fields != nil {
		if traceID, ok := entry.Fields["trace_id"].(string); ok {
			entry.TraceID = traceID
			delete(entry.Fields, "trace_id")
		}
		if requestID, ok := entry.Fields["request_id"].(string); ok {
			entry.RequestID = requestID
			delete(entry.Fields, "request_id")
		}
		if err, ok := entry.Fields["error"].(string); ok {
			entry.Error = err
			delete(entry.Fields, "error")
		}
		if err, ok := entry.Fields["error"].(error); ok {
			entry.Error = err.Error()
			delete(entry.Fields, "error")
		}
	}

	if rcv.config.Structured {
		// JSON format
		if jsonBytes, err := json.Marshal(entry); err == nil {
			fmt.Fprintln(rcv.output, string(jsonBytes))
		} else {
			// Fallback to simple format
			fmt.Fprintf(rcv.output, "[%s] %s [%s] %s/%s: %s\n",
				entry.Timestamp, entry.Level, entry.Code, entry.Component, entry.Service, entry.Message)
		}
	} else {
		// Human-readable format
		fmt.Fprintf(rcv.output, "[%s] %s [%s] %s/%s: %s",
			entry.Timestamp, entry.Level, entry.Code, entry.Component, entry.Service, entry.Message)
		if len(entry.Fields) > 0 {
			if fieldsJSON, err := json.Marshal(entry.Fields); err == nil {
				fmt.Fprintf(rcv.output, " %s", string(fieldsJSON))
			}
		}
		fmt.Fprintln(rcv.output)
	}
}

// DEBUG logs a debug message using MCode
func (rcv *LOGGER) DEBUG(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	rcv.log(DEBUG, mcode, optionalMessage, f)
}

// INFO logs an info message using MCode
func (rcv *LOGGER) INFO(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	rcv.log(INFO, mcode, optionalMessage, f)
}

// WARN logs a warning message using MCode
func (rcv *LOGGER) WARN(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	rcv.log(WARN, mcode, optionalMessage, f)
}

// ERROR logs an error message using MCode
func (rcv *LOGGER) ERROR(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	rcv.log(ERROR, mcode, optionalMessage, f)
}

// FATAL logs a fatal message using MCode and exits
func (rcv *LOGGER) FATAL(mcode MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	rcv.log(FATAL, mcode, optionalMessage, f)
	// Indirect exit to allow test override
	exitFunc(1)
}

// exitFunc allows tests to override process exit behavior for FATAL.
var exitFunc = os.Exit
