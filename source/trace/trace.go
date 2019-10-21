package trace

import (
	"runtime"
	"strconv"
	"strings"
)

type ITrace interface {
	SafetyLevel() uint8
	Add(message string)
	String() string
}

type IErrorTraceFactory interface {
	CreateTrace(base error, level uint8) ITrace
}

type simpleErrorTrace struct {
	Base         error
	Level        uint8
	messageStack []string
}

type detailedErrorTrace struct {
	Base      error
	Level     uint8
	stepStack []*detailedTraceStep
}

type detailedTraceStep struct {
	message  string
	location string
}

type SimpleErrorTraceFactory struct {
}

type DetailedErrorTraceFactory struct {
}

func filename(filepath string) (name string) {
	index := strings.LastIndex(filepath, "/")
	if index == -1 {
		return filepath
	}
	return filepath[index+1:]
}

/**
Unfinished, now it duplicates func filename
*/
func funcname(fullname string) (name string) {
	index := strings.LastIndex(fullname, "/")
	if index == -1 {
		return fullname
	}
	return fullname[index+1:]
}

func NewErrorTraceFactory(factoryType string) (factory IErrorTraceFactory) {
	switch factoryType {
	case "detailed":
		factory = &DetailedErrorTraceFactory{}
	case "simple":
		factory = &SimpleErrorTraceFactory{}
	default:
		factory = &SimpleErrorTraceFactory{}
	}
	return
}

func (factory *SimpleErrorTraceFactory) CreateTrace(base error, level uint8) (trace ITrace) {
	trace = &simpleErrorTrace{
		Base:         base,
		Level:        level,
		messageStack: make([]string, 0),
	}

	return
}

func (factory *DetailedErrorTraceFactory) CreateTrace(base error, level uint8) (trace ITrace) {
	trace = &detailedErrorTrace{
		Base:      base,
		Level:     level,
		stepStack: make([]*detailedTraceStep, 0),
	}

	return
}

func (trace *simpleErrorTrace) SafetyLevel() uint8 {
	return trace.Level
}

func (trace *detailedErrorTrace) SafetyLevel() uint8 {
	return trace.Level
}

func (trace *simpleErrorTrace) Add(message string) {
	trace.messageStack = append(trace.messageStack, message)
}

func (trace *detailedErrorTrace) Add(message string) {
	var location string
	pcounter, filepath, linenumber, ok := runtime.Caller(1)
	if ok {
		details := runtime.FuncForPC(pcounter)
		location = funcname(details.Name()) + " | " + filename(filepath) + ": " + strconv.Itoa(linenumber)
	} else {
		location = "Information about location in the code is unavailable"
	}

	step := &detailedTraceStep{
		message:  message,
		location: location,
	}
	trace.stepStack = append(trace.stepStack, step)
}

func (trace *simpleErrorTrace) String() (result string) {
	length := len(trace.messageStack)
	indent := ""

	for i := 1; i <= length; i++ {
		result += indent + trace.messageStack[length-i] + ":\n"
		indent += "    "
	}
	result += indent + "[Cause] " + trace.Base.Error()

	return
}

func (trace *detailedErrorTrace) String() (result string) {
	length := len(trace.stepStack)
	indent := ""

	for i := 1; i <= length; i++ {
		result += indent + "[Caller] | " + trace.stepStack[length-i].location + "\n"
		result += indent + trace.stepStack[length-i].message + ":\n"
		indent += "    "
	}
	result += indent + "[Cause] " + trace.Base.Error()

	return
}
