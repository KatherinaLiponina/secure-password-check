package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetDictionaryFromFile(filename string) (map[string]struct{}, error) {
	// word by line
	dict := make(map[string]struct{})

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		dict[fileScanner.Text()] = struct{}{}
	}

	file.Close()
	return dict, nil
}

func GetDictionaryWithFrequency(filename string) (map[string]float64, error) {
	// number, frequency, word, type
	dict := make(map[string]float64)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 4 {
			return nil, fmt.Errorf("unexpected length of line %s", line)
		}
		val, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse %s", parts[1])
		}
		dict[parts[2]] = val
	}

	file.Close()
	return dict, nil
}
