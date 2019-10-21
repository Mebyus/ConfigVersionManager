package executor

import (
	"../../trace"
	"../command"
)

type IExecutor interface {
	Match(name string) bool
	Load(command *command.Command)
	Validate() trace.ITrace
	Execute() trace.ITrace
	ChangeEtraceFactory(factory trace.IErrorTraceFactory)
}

type Executor struct {
	name          string
	command       *command.Command
	etraceFactory trace.IErrorTraceFactory
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
