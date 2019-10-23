package executor

import (
	"os"

	"../../filework"
	"../../log"
	"../../trace"
	"../dispatcher"
)

type BackupExecutor struct {
	dispatcher.Executor
}

func NewBackupExecutor() *BackupExecutor {
	return &BackupExecutor{
		dispatcher.Executor{
			Name: "backup",
		},
	}
}

func (executor *BackupExecutor) Validate() trace.ITrace {
	return nil
}

func (executor *BackupExecutor) Execute() trace.ITrace {
	filepath := "./app.conf"
	backpath := "./backup/app.conf"
	backdir := "./backup"

	if newpath, ok := executor.Command.Extract("f"); ok {
		filepath = newpath
	}

	if err := os.MkdirAll(backdir, os.ModeDir); err != nil {
		etrace := executor.EtraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to create directory")
		return etrace
	}

	reader := filework.NewFileReader(executor.EtraceFactory)

	filebytes, etrace := reader.Read(filepath)
	if etrace != nil {
		etrace.Add("Tried to get content of the origin file")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		executor.Logger.LogTrace(etrace)
	}

	writer := filework.NewFileWriter(executor.EtraceFactory)
	writer.Save(backpath, filebytes)
	if etrace != nil {
		etrace.Add("Tried to save backup file")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		executor.Logger.LogTrace(etrace)
	}

	return nil
}
