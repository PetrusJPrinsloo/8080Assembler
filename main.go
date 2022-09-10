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

func getBytes(s string) string {
	cleaned := stripByte(s)
	//red first 2 characters
	firstTwo := cleaned[:2]
	secondTwo := cleaned[2:4]
	return secondTwo + "" + firstTwo
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
	opcode := strings.ToUpper(parts[1])

	//get the rest of the parts
	var args string
	if len(parts) > 2 {
		args = parts[2]
	}

	var bytes []string
	//switch on the opcode
	switch opcode {
	case "NOP":
		assembled = "00"
		break
	case "LXI":
		bytes = convertToBytes(args)
		switch parts[2] {
		case "B":
			assembled = "01" + bytes[1] + bytes[0]
			break
		case "D":
			assembled = "11" + bytes[1] + bytes[0]
			break
		case "H":
			assembled = "21" + bytes[1] + bytes[0]
			break
		case "SP":
			assembled = "31" + bytes[1] + bytes[0]
			break
		}
		break
	case "STAX":
		switch parts[2] {
		case "B":
			assembled = "02"
			break
		case "D":
			assembled = "12"
			break
		}
		break
	case "INX":
		switch parts[2] {
		case "B":
			assembled = "03"
			break
		case "D":
			assembled = "13"
			break
		case "H":
			assembled = "23"
			break
		case "SP":
			assembled = "33"
			break
		}
		break
	case "INR":
		switch parts[2] {
		case "B":
			assembled = "04"
			break
		case "C":
			assembled = "0C"
			break
		case "D":
			assembled = "14"
			break
		case "E":
			assembled = "1C"
			break
		case "H":
			assembled = "24"
			break
		case "L":
			assembled = "2C"
			break
		case "M":
			assembled = "34"
			break
		case "A":
			assembled = "3C"
			break
		}
		break
	case "DCR":
		switch parts[2] {
		case "B":
			assembled = "05"
			break
		case "C":
			assembled = "0D"
			break
		case "D":
			assembled = "15"
			break
		case "E":
			assembled = "1D"
			break
		case "H":
			assembled = "25"
			break
		case "L":
			assembled = "2D"
			break
		case "M":
			assembled = "35"
			break
		case "A":
			assembled = "3D"
			break
		}
		break
	case "MVI":
		bytes = convertToBytes(args)
		switch parts[2] {
		case "B":
			assembled = "06" + bytes[1]
			break
		case "C":
			assembled = "0E" + bytes[1]
			break
		case "D":
			assembled = "16" + bytes[1]
			break
		case "E":
			assembled = "1E" + bytes[1]
			break
		case "H":
			assembled = "26" + bytes[1]
			break
		case "L":
			assembled = "2E" + bytes[1]
			break
		case "M":
			assembled = "36" + bytes[1]
			break
		case "A":
			assembled = "3E" + bytes[1]
			break
		}
		break
	case "RLC":
		assembled = "07"
		break
	case "DAD":
		switch parts[2] {
		case "B":
			assembled = "09"
			break
		case "D":
			assembled = "19"
			break
		case "H":
			assembled = "29"
			break
		case "SP":
			assembled = "39"
			break
		}
		break
	case "LDAX":
		switch parts[2] {
		case "B":
			assembled = "0A"
			break
		case "D":
			assembled = "1A"
			break
		}
		break
	case "DCX":
		switch parts[2] {
		case "B":
			assembled = "0B"
			break
		case "D":
			assembled = "1B"
			break
		case "H":
			assembled = "2B"
			break
		case "SP":
			assembled = "3B"
			break
		}
		break
	case "RRC":
		assembled = "0F"
		break
	case "RAL":
		assembled = "17"
		break
	case "RAR":
		assembled = "1F"
		break
	case "DAA":
		assembled = "27"
		break
	case "CMA":
		assembled = "2F"
		break
	case "STC":
		assembled = "37"
		break
	case "CMC":
		assembled = "3F"
		break
	case "MOV":
		regsters := convertToBytes(parts[2])
		switch regsters[0] {
		case "B":
			switch regsters[1] {
			case "B":
				assembled = "40"
				break
			case "C":
				assembled = "41"
				break
			case "D":
				assembled = "42"
				break
			case "E":
				assembled = "43"
				break
			case "H":
				assembled = "44"
				break
			case "L":
				assembled = "45"
				break
			case "M":
				assembled = "46"
				break
			case "A":
				assembled = "47"
				break
			}
			break
		case "C":
			switch regsters[1] {
			case "B":
				assembled = "48"
				break
			case "C":
				assembled = "49"
				break
			case "D":
				assembled = "4A"
				break
			case "E":
				assembled = "4B"
				break
			case "H":
				assembled = "4C"
				break
			case "L":
				assembled = "4D"
				break
			case "M":
				assembled = "4E"
				break
			case "A":
				assembled = "4F"
				break
			}
			break
		case "D":
			switch regsters[1] {
			case "B":
				assembled = "50"
				break
			case "C":
				assembled = "51"
				break
			case "D":
				assembled = "52"
				break
			case "E":
				assembled = "53"
				break
			case "H":
				assembled = "54"
				break
			case "L":
				assembled = "55"
				break
			case "M":
				assembled = "56"
				break
			case "A":
				assembled = "57"
				break
			}
			break
		case "E":
			switch regsters[1] {
			case "B":
				assembled = "58"
				break
			case "C":
				assembled = "59"
				break
			case "D":
				assembled = "5A"
				break
			case "E":
				assembled = "5B"
				break
			case "H":
				assembled = "5C"
				break
			case "L":
				assembled = "5D"
				break
			case "M":
				assembled = "5E"
				break
			case "A":
				assembled = "5F"
				break
			}
			break
		case "H":
			switch regsters[1] {
			case "B":
				assembled = "60"
				break
			case "C":
				assembled = "61"
				break
			case "D":
				assembled = "62"
				break
			case "E":
				assembled = "63"
				break
			case "H":
				assembled = "64"
				break
			case "L":
				assembled = "65"
				break
			case "M":
				assembled = "66"
				break
			case "A":
				assembled = "67"
				break
			}
			break
		case "L":
			switch regsters[1] {
			case "B":
				assembled = "68"
				break
			case "C":
				assembled = "69"
				break
			case "D":
				assembled = "6A"
				break
			case "E":
				assembled = "6B"
				break
			case "H":
				assembled = "6C"
				break
			case "L":
				assembled = "6D"
				break
			case "M":
				assembled = "6E"
				break
			case "A":
				assembled = "6F"
				break
			}
			break
		case "M":
			switch regsters[1] {
			case "B":
				assembled = "70"
				break
			case "C":
				assembled = "71"
				break
			case "D":
				assembled = "72"
				break
			case "E":
				assembled = "73"
				break
			case "H":
				assembled = "74"
				break
			case "L":
				assembled = "75"
				break
			case "A":
				assembled = "77"
				break
			}
			break
		case "A":
			switch regsters[1] {
			case "B":
				assembled = "78"
				break
			case "C":
				assembled = "79"
				break
			case "D":
				assembled = "7A"
				break
			case "E":
				assembled = "7B"
				break
			case "H":
				assembled = "7C"
				break
			case "L":
				assembled = "7D"
				break
			case "M":
				assembled = "7E"
				break
			case "A":
				assembled = "7F"
				break
			}
		}
	case "HLT":
		assembled = "76"
		break
	case "ADD":
		switch parts[2] {
		case "B":
			assembled = "80"
			break
		case "C":
			assembled = "81"
			break
		case "D":
			assembled = "82"
			break
		case "E":
			assembled = "83"
			break
		case "H":
			assembled = "84"
			break
		case "L":
			assembled = "85"
			break
		case "M":
			assembled = "86"
			break
		case "A":
			assembled = "87"
			break
		}
		break
	case "ADC":
		switch parts[2] {
		case "B":
			assembled = "88"
			break
		case "C":
			assembled = "89"
			break
		case "D":
			assembled = "8A"
			break
		case "E":
			assembled = "8B"
			break
		case "H":
			assembled = "8C"
			break
		case "L":
			assembled = "8D"
			break
		case "M":
			assembled = "8E"
			break
		case "A":
			assembled = "8F"
			break
		}
		break
	case "SUB":
		switch parts[2] {
		case "B":
			assembled = "90"
			break
		case "C":
			assembled = "91"
			break
		case "D":
			assembled = "92"
			break
		case "E":
			assembled = "93"
			break
		case "H":
			assembled = "94"
			break
		case "L":
			assembled = "95"
			break
		case "M":
			assembled = "96"
			break
		case "A":
			assembled = "97"
			break
		}
		break
	case "SBB":
		switch parts[2] {
		case "B":
			assembled = "98"
			break
		case "C":
			assembled = "99"
			break
		case "D":
			assembled = "9A"
			break
		case "E":
			assembled = "9B"
			break
		case "H":
			assembled = "9C"
			break
		case "L":
			assembled = "9D"
			break
		case "M":
			assembled = "9E"
			break
		case "A":
			assembled = "9F"
			break
		}
		break
	case "ANA":
		switch parts[2] {
		case "B":
			assembled = "A0"
			break
		case "C":
			assembled = "A1"
			break
		case "D":
			assembled = "A2"
			break
		case "E":
			assembled = "A3"
			break
		case "H":
			assembled = "A4"
			break
		case "L":
			assembled = "A5"
			break
		case "M":
			assembled = "A6"
			break
		case "A":
			assembled = "A7"
			break
		}
		break
	case "XRA":
		switch parts[2] {
		case "B":
			assembled = "A8"
			break
		case "C":
			assembled = "A9"
			break
		case "D":
			assembled = "AA"
			break
		case "E":
			assembled = "AB"
			break
		case "H":
			assembled = "AC"
			break
		case "L":
			assembled = "AD"
			break
		case "M":
			assembled = "AE"
			break
		case "A":
			assembled = "AF"
			break
		}
		break
	case "ORA":
		switch parts[2] {
		case "B":
			assembled = "B0"
			break
		case "C":
			assembled = "B1"
			break
		case "D":
			assembled = "B2"
			break
		case "E":
			assembled = "B3"
			break
		case "H":
			assembled = "B4"
			break
		case "L":
			assembled = "B5"
			break
		case "M":
			assembled = "B6"
			break
		case "A":
			assembled = "B7"
			break
		}
		break
	case "CMP":
		switch parts[2] {
		case "B":
			assembled = "B8"
			break
		case "C":
			assembled = "B9"
			break
		case "D":
			assembled = "BA"
			break
		case "E":
			assembled = "BB"
			break
		case "H":
			assembled = "BC"
			break
		case "L":
			assembled = "BD"
			break
		case "M":
			assembled = "BE"
			break
		case "A":
			assembled = "BF"
			break
		}
		break
	case "RNZ":
		assembled = "C0"
		break
	case "POP":
		switch parts[2] {
		case "B":
			assembled = "C1"
			break
		case "D":
			assembled = "D1"
			break
		case "H":
			assembled = "E1"
			break
		case "PSW":
			assembled = "F1"
			break
		}
		break
	case "JNZ":
		assembled = "C2" + getBytes(parts[2])
		break
	case "JMP":
		assembled = "C3" + getBytes(parts[2])
		break
	case "CNZ":
		assembled = "C4" + getBytes(parts[2])
		break
	case "PUSH":
		switch parts[2] {
		case "B":
			assembled = "C5"
			break
		case "D":
			assembled = "D5"
			break
		case "H":
			assembled = "E5"
			break
		case "PSW":
			assembled = "F5"
			break
		}
		break
	case "ADI":
		assembled = "C6" + parts[2]
		break
	case "RST":
		switch parts[2] {
		case "0":
			assembled = "C7"
			break
		case "1":
			assembled = "CF"
			break
		case "2":
			assembled = "D7"
			break
		case "3":
			assembled = "DF"
			break
		case "4":
			assembled = "E7"
			break
		case "5":
			assembled = "EF"
			break
		case "6":
			assembled = "F7"
			break
		case "7":
			assembled = "FF"
			break
		}
		break
	case "RZ":
		assembled = "C8"
		break
	case "RET":
		assembled = "C9"
		break
	case "JZ":
		assembled = "CA" + getBytes(parts[2])
		break
	case "CZ":
		assembled = "CC" + getBytes(parts[2])
		break
	case "CALL":
		assembled = "CD" + getBytes(parts[2])
		break
	case "ACI":
		assembled = "CE" + parts[2]
		break
	case "RNC":
		assembled = "D0"
		break
	case "JNC":
		assembled = "D2" + getBytes(parts[2])
		break
	case "OUT":
		assembled = "D3" + parts[2]
		break
	case "CNC":
		assembled = "D4" + getBytes(parts[2])
		break
	case "SUI":
		assembled = "D6" + parts[2]
		break
	case "RC":
		assembled = "D8"
		break
	case "JC":
		assembled = "DA" + getBytes(parts[2])
		break
	case "IN":
		assembled = "DB" + parts[2]
		break
	case "CC":
		assembled = "DC" + getBytes(parts[2])
		break
	case "SBI":
		assembled = "DE" + parts[2]
		break
	case "RPO":
		assembled = "E0"
		break
	case "JPO":
		assembled = "E2" + getBytes(parts[2])
		break
	case "XTHL":
		assembled = "E3"
		break
	case "CPO":
		assembled = "E4" + getBytes(parts[2])
		break
	case "ANI":
		assembled = "E6" + parts[2]
		break
	case "RPE":
		assembled = "E8"
		break
	case "PCHL":
		assembled = "E9"
		break
	case "JPE":
		assembled = "EA" + getBytes(parts[2])
		break
	case "XCHG":
		assembled = "EB"
		break
	case "CPE":
		assembled = "EC" + getBytes(parts[2])
		break
	case "XRI":
		assembled = "EE" + parts[2]
		break
	case "RP":
		assembled = "F0"
		break
	case "JP":
		assembled = "F2" + getBytes(parts[2])
		break
	case "DI":
		assembled = "F3"
		break
	case "CP":
		assembled = "F4" + getBytes(parts[2])
		break
	case "ORI":
		assembled = "F6" + parts[2]
		break
	case "RM":
		assembled = "F8"
		break
	case "SPHL":
		assembled = "F9"
		break
	case "JM":
		assembled = "FA" + getBytes(parts[2])
		break
	case "EI":
		assembled = "FB"
		break
	case "CM":
		assembled = "FC" + getBytes(parts[2])
		break
	case "CPI":
		assembled = "FE" + parts[2]
		break

	default:
		fmt.Println("Invalid opcode")
		//os.Exit(1)
	}
	return assembled
}
