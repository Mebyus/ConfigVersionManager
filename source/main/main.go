package main

import (
	"fmt"
	"os"
	"strings"
)

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

func readFileToString(filepath string) (string, error) {
	var (
		result    string
		filebytes []byte
	)

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Unable to read local config file: ", err)
		return result, err
	}

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Unable to retrieve local config file stats: ", err)
		return result, err
	}

	filebytes = make([]byte, stat.Size())
	_, err = file.Read(filebytes)
	if err != nil {
		fmt.Println("Unable to read local config file: ", err)
		return result, err
	}
	result = string(filebytes)

	err = file.Close()
	if err != nil {
		fmt.Println("Unable to close local conf file: ", err)
		return result, err
	}

	return result, nil
}

func saveLinesToFile(filepath string, lines []string) error {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Unable to create output file: ", err)
		return err
	}

	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		fmt.Println("Unable to write output to file: ", err)
		return err
	}

	return nil
}

func main() {
	var (
		keeppath  string
		confpath  string
		localpath string
		mergepath string
	)

	keeppath = "./app.conf.keep"
	confpath = "./app.conf"
	localpath = "./app.conf.local"
	mergepath = "./app.conf.merge"

	keepstr, _ := readFileToString(keeppath)
	confstr, _ := readFileToString(confpath)
	localstr, _ := readFileToString(localpath)

	confLines := stringToLines(confstr)
	localLines := stringToLines(localstr)
	keepLines := stringToLines(keepstr)

	for _, line := range keepLines {
		index := findLine(localLines, line)
		if index != -1 {
			substituteLine(&confLines, line, localLines[index])
		}
	}

	_ = saveLinesToFile(mergepath, confLines)
}
