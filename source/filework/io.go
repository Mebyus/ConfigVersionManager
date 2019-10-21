package filework

import (
	"os"
	"strings"

	"../log"
	"../trace"
)

type FileReader struct {
	etraceFactory trace.IErrorTraceFactory
}

type FileWriter struct {
	etraceFactory trace.IErrorTraceFactory
}

func NewFileReader(etraceFactory trace.IErrorTraceFactory) *FileReader {
	return &FileReader{
		etraceFactory: etraceFactory,
	}
}

func NewFileWriter(etraceFactory trace.IErrorTraceFactory) *FileWriter {
	return &FileWriter{
		etraceFactory: etraceFactory,
	}
}

func (reader *FileReader) Read(filepath string) ([]byte, trace.ITrace) {
	file, err := os.Open(filepath)
	if err != nil {
		etrace := reader.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to open file for reading")
		return nil, etrace
	}

	stat, err := file.Stat()
	if err != nil {
		etrace := reader.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to retrieve file information")
		return nil, etrace
	}

	filebytes := make([]byte, stat.Size())
	_, err = file.Read(filebytes)
	if err != nil {
		etrace := reader.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to read file content")
		return nil, etrace
	}

	err = file.Close()
	if err != nil {
		etrace := reader.etraceFactory.CreateTrace(err, log.WARN)
		etrace.Add("Tried to close file")
		return filebytes, etrace
	}

	return filebytes, nil
}

func (reader *FileReader) ReadString(filepath string) (string, trace.ITrace) {
	bytes, etrace := reader.Read(filepath)
	if etrace != nil {
		etrace.Add("Tried to get bytes from file")
		if etrace.SafetyLevel() < log.WARN {
			return "", etrace
		}
	}
	result := string(bytes)

	return result, etrace
}

func (writer *FileWriter) Save(filepath string, bytes []byte) trace.ITrace {
	file, err := os.Create(filepath)
	if err != nil {
		etrace := writer.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to create file")
		return etrace
	}

	if _, err := file.Write(bytes); err != nil {
		etrace := writer.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to save bytes to file")
		return etrace
	}

	if err = file.Close(); err != nil {
		etrace := writer.etraceFactory.CreateTrace(err, log.WARN)
		etrace.Add("Tried to close file")
		return etrace
	}

	return nil
}

func (writer *FileWriter) SaveLines(filepath string, lines []string) trace.ITrace {
	file, err := os.Create(filepath)
	if err != nil {
		etrace := writer.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to create file")
		return etrace
	}

	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		etrace := writer.etraceFactory.CreateTrace(err, log.ERROR)
		etrace.Add("Tried to write to file")
		return etrace
	}

	if err = file.Close(); err != nil {
		etrace := writer.etraceFactory.CreateTrace(err, log.WARN)
		etrace.Add("Tried to close file")
		return etrace
	}

	return nil
}
