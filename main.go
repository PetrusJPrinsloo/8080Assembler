package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

	//create new file
	asmFile, err := os.OpenFile("output.h", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

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
	var assembled string
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		fmt.Println(scanner.Text())
		assembled = assembleLine(cleanString(line))
		writeToFile(asmFile, assembled)
	}
}

func convertToBytes(s string) []string {
	var r []string
	for i := 0; i < len(s); i += 2 {
		r = append(r, string(s[i]))
	}
	return r
}

func writeToFile(file *os.File, s string) {
	_, err := file.WriteString(s)
	if err != nil {
		log.Fatal(err)
	}
}

func stripByte(s string) string {
	split := strings.Split(s, ",")

	if len(split) == 2 {
		s = split[1]
	}
	//strip prefix
	re := regexp.MustCompile(`[$#]+`)
	return re.ReplaceAllString(s, "")
}

func cleanString(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}

func assembleLine(line string) string {
	var assembled string
	//split line into parts
	parts := strings.Split(line, " ")
	//get the first part of the line
	opcode := strings.ToLower(parts[1])
	//get the rest of the parts
	args := parts[2]
	var bytes []string
	//switch on the opcode
	switch opcode {
	case "nop":
		assembled = "00"
		break
	case "lxi":
		assembled = "01"
		bytes = convertToBytes(args)
		assembled += bytes[1] + bytes[0]
		break
	case "stax":
		assembled = "02"
		//todo: check if B, C, D, E, H, L, M, A
		break
	case "inx":
		assembled = "03"
		//todo: check if B, C, D, E, H, L, M, A
		break
	case "inr":
		assembled = "04"
		//todo: check if B, C, D, E, H, L, M, A
		break
	case "dcr":
		assembled = "05"
		//todo: check if B, C, D, E, H, L, M, A
		break
	default:
		fmt.Println("Invalid opcode")
		//os.Exit(1)
	}
	return assembled
}
