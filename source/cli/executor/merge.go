package executor

import (
	"strings"

	"../../filework"
	"../../log"
	"../../trace"
	"../dispatcher"
)

type MergeExecutor struct {
	dispatcher.Executor
}

func NewMergeExecutor() *MergeExecutor {
	return &MergeExecutor{
		dispatcher.Executor{
			Name: "merge",
		},
	}
}

func (executor *MergeExecutor) Validate() trace.ITrace {
	return nil
}

func (executor *MergeExecutor) Execute() trace.ITrace {
	reader := filework.NewFileReader(executor.EtraceFactory)
	writer := filework.NewFileWriter(executor.EtraceFactory)

	keeppath := "./app.conf.keep"
	confpath := "./app.conf"
	localpath := "./app.conf.local"
	mergepath := "./app.conf.merge"

	keepstr, etrace := reader.ReadString(keeppath)
	if etrace != nil {
		etrace.Add("Something went wrong while reading list of local settings")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		executor.Logger.LogTrace(etrace)
	}

	confstr, etrace := reader.ReadString(confpath)
	if etrace != nil {
		etrace.Add("Something went wrong while reading distributed config file")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		executor.Logger.LogTrace(etrace)
	}

	localstr, etrace := reader.ReadString(localpath)
	if etrace != nil {
		etrace.Add("Something went wrong while reading local config file")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		executor.Logger.LogTrace(etrace)
	}

	confLines := stringToLines(confstr)
	localLines := stringToLines(localstr)
	keepLines := stringToLines(keepstr)

	for _, line := range keepLines {
		index := findLine(localLines, line)
		if index != -1 {
			substituteLine(&confLines, line, localLines[index])
		}
	}

	etrace = writer.SaveLines(mergepath, confLines)
	if etrace != nil {
		etrace.Add("Something went wrong while saving merged config file")
		if etrace.SafetyLevel() < log.WARN {
			return etrace
		}
		executor.Logger.LogTrace(etrace)
	}

	return nil
}

/**
Описание:
	Находит индекс первой строки из среза, которая начинается на заданную строку.

Аргументы:
	[1] lines - срез строк, по которым будет производится поиск.
	[2] searchStr - строка, на которую должна начинаться искомая.

Возвращает:
	[1] Индекс (начинается с 0) строки или -1, в случае если подходящей строки в срезе не найдено.

Примечание:
	На время разработки к строке поиска прибавляется " =" для исключения случаев вида:

	abcx = 123
	abc = 123

	В данном случае поиск по "abc" нашел бы строку "abcx = 123", без такого дополнения.

TODO:
	Распознавание ситуации, описанной в примечании, должно быть более интеллектуальным,
	с распознаванием случаев вида:

	abc =123
	abc=123
	abc= 123
	abc = 123
	abc  = 123
	...

	С отбрасыванием случаев вида:

	abcx= 123
	abcx  = 123
	...

	Возможно стоит использовать регулярное выражение.
*/
func findLine(lines []string, searchStr string) int32 {
	for i, line := range lines {
		if strings.HasPrefix(line, searchStr+" =") {
			return int32(i)
		}
	}
	return -1
}

/**
Описание:
	Разделяет строку, содержащую символы новой строки ("\n"), на срез строк.
	Каждый элемент среза это отдельная строка, как если бы она была строкой файла.
	В дополнение к этому, у каждой такой строки обрезаются с краев символы:
	" ", "\r", "\t" (в любом количестве).

Аргументы:
	[1] str - строка для разделения на срез

Возвращает:
	[1] Срез строк

Примечание:
	Если символов новой строки нет, то срез будет состоять из одной строки.
	"\r" - символ перевода каретки, обычно используется в Windows вместе с "\n".
	"\t" - символ табуляции
*/
func stringToLines(str string) []string {
	lines := strings.Split(str, "\n")

	for i, line := range lines {
		lines[i] = strings.Trim(line, " \r\t")
	}
	return lines
}

/**
Описание:
	Заменяет в срезе строку, начинающуюся на заданную, на другую строку.

Аргументы:
	[1] linesPtr - указатель на срез, в котором надо заменить строку.
	[2] searchStr - строка, на которую должна начинаться искомая (см.
		описание функции findLine).
	[3] target - строка, на которую будет произведена замена.

Примечание:
	Если нужная строка не найдена, не делает ничего.
	Заменяется первая подходящая строка, остальные, если они есть, остаются неизменными.
*/
func substituteLine(linesPtr *[]string, searchStr, target string) {
	lines := *linesPtr

	index := findLine(lines, searchStr)
	if index != -1 {
		lines[index] = target
	}
}
