package executor

import (
	"os"

	"../../filework"
	"../../log"
	"../../trace"
)

type BackupExecutor struct {
	Executor
}

func NewBackupExecutor() *BackupExecutor {
	return &BackupExecutor{
		Executor{
			name: "backup",
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

	if err := os.MkdirAll(backdir, os.ModeDir); err != nil {
		etrace := executor.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to create directory")
		return etrace
	}

	reader := filework.NewFileReader(executor.etraceFactory)

	filebytes, etrace := reader.Read(filepath)
	if etrace != nil {
		etrace.Add("Tried to get content of the origin file")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		// logger.Log(etrace.String(), etrace.SafetyLevel())
	}

	writer := filework.NewFileWriter(executor.etraceFactory)
	writer.Save(backpath, filebytes)
	if etrace != nil {
		etrace.Add("Tried to save backup file")
		return etrace
	}

	return nil
}
