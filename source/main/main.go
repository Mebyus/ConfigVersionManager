package main

import (
	"fmt"
	"os"
	"strings"
)

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
		keeppath string
		confpath string
	)

	keeppath = "./app.conf.local"
	confpath = "./app.conf"

	keepstr, _ := readFileToString(keeppath)
	confstr, _ := readFileToString(confpath)

	fmt.Println(keepstr)

	for n, s := range strings.Split(confstr, "\n") {
		fmt.Println(n, s)
	}
}
