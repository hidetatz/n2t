package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var symbolTable = map[string]int16{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SCREEN": 16384,
	"KBD":    24576,
}

var comps = map[string]string{
	// a = 0
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",
	// a = 1
	"M":   "1110000",
	"!M":  "1110001",
	"-M":  "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
}

var dests = map[string]string{
	"":    "000", // dest can be empty
	"M":   "001",
	"D":   "010",
	"MD":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "110",
	"AMD": "111",
}

var jumps = map[string]string{
	"":    "000", // jump can be empty
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

func main() {
	filename := os.Args[1]
	if filename == "" {
		fmt.Fprintf(os.Stderr, "file name must be passed\n")
		os.Exit(1)
	}

	if !strings.HasSuffix(filename, ".asm") {
		fmt.Fprintf(os.Stderr, "passed file is not an assembly file\n")
		os.Exit(1)
	}

	asm, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open file: %v\n", err)
		os.Exit(1)
	}
	defer asm.Close()

	write := func(l string) {
		fmt.Println(l)
	}

	buff := []string{}
	scanner := bufio.NewScanner(asm)
	index := 0
	// First, scan program just to create symbol table.
	// It is necessary to resolve the instruction address on memory.
	// In this first loop, it removes all the comments and create symbol table, then
	// register labels symbols to the table if found.
	for scanner.Scan() {
		line := scanner.Text()

		// Trim just in case
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		// skip comment line
		if strings.HasPrefix(line, "//") {
			continue
		}

		if strings.Contains(line, "//") {
			line = strings.Split(line, "//")[0]
			// TrimSpace does not trim mid-spaces. We need to trim spaces between instruction and "//" here
			// using Regex will resolve this problem but it seems like this is sufficient for now
			line = strings.TrimSpace(line)
		}

		// Register label symbol to the table
		if strings.HasPrefix(line, "(") {
			l := strings.TrimLeft(line, "(")
			l = strings.TrimRight(l, ")")
			// index is the current line. It should be also used as the instruction pointer.
			symbolTable[l] = int16(index)
		}

		// store line to the buff for the second loop.
		buff = append(buff, line)
		index++
	}

	if err = scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan file: %v\n", err)
		os.Exit(1)
	}

	// variables should be stored from address 16 since 0~15 is already used by R0-R15
	addr := 16
	// In the second loop, it parses lines and converts them into the Hack machine code.
	for _, line := range buff {
		// labels symbol is already resolved. Just skip it
		if strings.HasPrefix(line, "(") {
			continue
		}

		// A-instruction starts with "@".
		// In this if condition, the line can be one of:
		// * constant number (e.g. 0, 2, 5)
		// * symbol in symbol table (e.g. R0, SP, or label symbols)
		// * variable
		// Machine code format is: 0xxx xxxx xxxx xxxx (while xxx... is the binary value).
		if strings.HasPrefix(line, "@") {
			line = strings.TrimLeft(line, "@")
			a, ok := symbolTable[line]
			if ok {
				// symbol is found in the table.
				// write the instruction address.
				write("0" + fmt.Sprintf("%015b", a))
				continue
			}

			if !ok {
				if n, err := strconv.Atoi(line); err == nil {
					// constant number. Just write them.
					write("0" + fmt.Sprintf("%015b", n))
				} else {
					// a variable. write the address which starts from 16.
					write("0" + fmt.Sprintf("%015b", addr))
					addr++
				}
				continue
			}
		}

		// Else, C-instruction.
		// the format is one of:
		// * dest=comp;jump
		// * dest=comp
		// * comp;jump
		// Each machine code is defined in comps, dests, and jumps.
		// Machine code format is: 111c cccc ccdd djjj.
		var dest, comp, jump string
		if strings.Contains(line, "=") && strings.Contains(line, ";") {
			spl := strings.Split(line, "=")
			dest = spl[0]
			sp := strings.Split(spl[1], ";")
			comp = sp[0]
			jump = sp[1]
		} else if strings.Contains(line, "=") {
			spl := strings.Split(line, "=")
			dest = spl[0]
			comp = spl[1]
		} else {
			spl := strings.Split(line, ";")
			comp = spl[0]
			jump = spl[1]
		}
		write("111" + comps[comp] + dests[dest] + jumps[jump])
	}
}
