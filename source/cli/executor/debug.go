package executor

import (
	"../../trace"
	"../dispatcher"
)

type DebugExecutor struct {
	dispatcher.Executor
}

func NewDebugExecutor() *DebugExecutor {
	return &DebugExecutor{
		dispatcher.Executor{
			Name: "debug",
		},
	}
}

func (executor *DebugExecutor) ShortHelp() string {
	return "short help for " + executor.Name + " command"
}

func (executor *DebugExecutor) Validate() trace.ITrace {
	return nil
}

func (executor *DebugExecutor) Execute() trace.ITrace {
	errorTraceFactory := trace.NewErrorTraceFactory("detailed")
	executor.Boss.ChangeEtraceFactory(errorTraceFactory)

	return nil
}
