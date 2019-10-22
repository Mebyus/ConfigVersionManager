package executor

import (
	"../../log"
	"../../trace"
	"../command"
)

type IExecutor interface {
	Match(name string) bool
	Load(command *command.Command)
	Validate() trace.ITrace
	Execute() trace.ITrace
	ChangeEtraceFactory(factory trace.IErrorTraceFactory)
	ChangeLogger(logger log.ILogger)
}

type Executor struct {
	name          string
	command       *command.Command
	etraceFactory trace.IErrorTraceFactory
	logger        log.ILogger
}

func (executor *Executor) Match(name string) bool {
	return executor.name == name
}

func (executor *Executor) Load(command *command.Command) {
	executor.command = command
}

func (executor *Executor) ChangeEtraceFactory(factory trace.IErrorTraceFactory) {
	executor.etraceFactory = factory
}

func (executor *Executor) ChangeLogger(logger log.ILogger) {
	executor.logger = logger
}
