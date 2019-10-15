package main

import (
	"fmt"
	"os"
	"strings"
)

func findLine(lines []string, searchStr string) int32 {
	for i, line := range lines {
		if strings.HasPrefix(line, searchStr+" =") {
			return int32(i)
		}
	}
	return -1
}

func stringToLines(str string) []string {
	lines := strings.Split(str, "\n")

	for i, line := range lines {
		lines[i] = strings.Trim(line, " \r")
	}
	return lines
}

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

func main() {
	var (
		keeppath  string
		confpath  string
		localpath string
	)

	keeppath = "./app.conf.keep"
	confpath = "./app.conf"
	localpath = "./app.conf.local"

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

	for i, line := range confLines {
		fmt.Println(i, line)
	}
}
