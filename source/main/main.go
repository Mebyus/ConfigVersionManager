package main

import (
	"os"

	"../cli/command"
	"../cli/dispatcher"
	"../cli/executor"
	"../log"
	"../trace"
)

func main() {
	logger := &log.Logger{
		Level:   log.INFO,
		InfoOut: os.Stdout,
	}
	errorTraceFactory := trace.NewErrorTraceFactory("detailed")
	dispatcher := dispatcher.NewDispatcher(logger, errorTraceFactory)
	command := command.ParseArgs(os.Args, "-")

	dispatcher.Register(executor.NewBackupExecutor())
	dispatcher.Register(executor.NewMergeExecutor())
	dispatcher.Dispatch(command)
}
