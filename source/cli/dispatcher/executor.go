package dispatcher

import (
	"../../log"
	"../../trace"
	"../command"
)

type IExecutor interface {
	Match(name string) bool
	BecomeSlave(boss *Dispatcher)
	Load(command *command.Command)
	Validate() trace.ITrace
	Execute() trace.ITrace
	ChangeEtraceFactory(factory trace.IErrorTraceFactory)
	ChangeLogger(logger log.ILogger)
}

type Executor struct {
	Boss          *Dispatcher
	Name          string
	Command       *command.Command
	EtraceFactory trace.IErrorTraceFactory
	Logger        log.ILogger
}

func (executor *Executor) Match(name string) bool {
	return executor.Name == name
}

func (executor *Executor) BecomeSlave(boss *Dispatcher) {
	executor.Boss = boss
}

func (executor *Executor) Load(command *command.Command) {
	executor.Command = command
}

func (executor *Executor) ChangeEtraceFactory(factory trace.IErrorTraceFactory) {
	executor.EtraceFactory = factory
}

func (executor *Executor) ChangeLogger(logger log.ILogger) {
	executor.Logger = logger
}
