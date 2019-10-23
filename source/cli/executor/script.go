package executor

import (
	"strings"

	"../../filework"
	"../../log"
	"../../trace"
	"../command"
	"../dispatcher"
)

type ScriptExecutor struct {
	dispatcher.Executor
}

func NewScriptExecutor() *ScriptExecutor {
	return &ScriptExecutor{
		dispatcher.Executor{
			Name: "script",
		},
	}
}

func (executor *ScriptExecutor) Validate() trace.ITrace {
	return nil
}

func (executor *ScriptExecutor) Execute() trace.ITrace {
	reader := filework.NewFileReader(executor.EtraceFactory)

	scriptPath := "./script.cvms"

	scriptStr, etrace := reader.ReadString(scriptPath)
	if etrace != nil {
		etrace.Add("Something went wrong while reading list of local settings")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		executor.Logger.LogTrace(etrace)
	}

	scriptLines := stringToLines(scriptStr)

	for _, line := range scriptLines {
		args := strings.Split("cvm "+line, " ")
		command := command.ParseArgs(args, "-")
		executor.Boss.Dispatch(command)
	}

	return nil
}
