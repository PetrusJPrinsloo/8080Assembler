package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: main.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	file, err := os.Open(filename)

	check(err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file")
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		fmt.Println(scanner.Text())
	}
}

func convert16to8bytes(s string) []string {
	var r []string
	for i := 0; i < len(s); i += 2 {
		r = append(r, string(s[i]))
	}
	return r
}
