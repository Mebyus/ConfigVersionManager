package executor

import (
	"fmt"

	"../../trace"
	"../dispatcher"
)

type HelpExecutor struct {
	dispatcher.Executor
}

func NewHelpExecutor() *HelpExecutor {
	return &HelpExecutor{
		dispatcher.Executor{
			Name: "help",
		},
	}
}

func (executor *HelpExecutor) ShortHelp() string {
	return "short help for " + executor.Name + " command"
}

func (executor *HelpExecutor) Validate() trace.ITrace {
	return nil
}

func (executor *HelpExecutor) Execute() trace.ITrace {
	for _, executor := range executor.Boss.Pool {
		fmt.Printf("    %s\n", executor.ShortHelp())
	}

	return nil
}
