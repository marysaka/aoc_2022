package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func expectNoError(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Elf struct {
	cal uint64
}

type ElfSlice []Elf

func (e ElfSlice) Len() int {
	return len(e)
}

func (e ElfSlice) Less(i, j int) bool {
	return e[i].cal < e[j].cal
}

func (e ElfSlice) Swap(i, j int) {
	tmp := e[i]

	e[i] = e[j]
	e[j] = tmp
}

func parseElfs(scanner *bufio.Scanner) (ElfSlice, error) {
	data := ElfSlice{}

	elf := Elf{cal: 0}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			data = append(data, elf)

			elf = Elf{}

			continue
		}

		val, err := strconv.ParseUint(line, 10, 64)

		if err != nil {
			return nil, err
		}

		elf.cal += val
	}

	return data, nil
}

func partGeneric(scanner *bufio.Scanner, elfCount int) {
	data, err := parseElfs(scanner)
	expectNoError(err, "Cannot parse Elfs")
	sort.Sort(sort.Reverse(data))

	var total uint64

	for _, e := range data[0:elfCount] {
		total += e.cal
	}

	fmt.Println(total)
}

func main() {
	var err error

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage ", os.Args[0], "<1|2> <file>")
		os.Exit(1)
	}

	partNum, err := strconv.ParseUint(os.Args[1], 10, 8)
	expectNoError(err, "Part Number must be an integer")

	file, err := os.Open(os.Args[2])
	defer file.Close()
	expectNoError(err, "Cannot open file")

	var data = bufio.NewScanner(file)

	switch partNum {
	case 1:
		partGeneric(data, 1)
		break
	case 2:
		partGeneric(data, 3)
		break
	default:
		fmt.Fprintln(os.Stderr, "Part number must be between 1 and 2")
		os.Exit(1)
		break
	}

	expectNoError(data.Err(), "Error during scanner read")
}
