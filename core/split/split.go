package split

import "strings"

func SplitPassword(password string) []string {
	var parts []string

	var part string
	var prevSymbolClass SymbolClass
	for _, symbol := range password {
		class := DetermineClass(symbol)
		if class == prevSymbolClass {
			part += string(symbol)
		} else {
			if len(part) > 0 && prevSymbolClass == Letter {
				parts = append(parts, strings.ToLower(part))
			}
			part = string(symbol)
			prevSymbolClass = class
		}
	}
	return parts
}

type SymbolClass int

const (
	Letter SymbolClass = iota
	Number
	SpecialSymbol
)

func DetermineClass(symbol rune) SymbolClass {
	if symbol >= '0' && symbol <= '9' {
		return Number
	}
	if symbol >= 'a' && symbol <= 'z' {
		return Letter
	}
	if symbol >= 'A' && symbol <= 'Z' {
		return Letter
	}
	return SpecialSymbol
}

func PossibleSplits(word string) []string {
	res := make([]string, 0)
	wordRuned := []rune(word)
	for n := 3; n < len(wordRuned); n++ {
		res = append(res, splitOnN(wordRuned, n)...)
	}
	return res
}

func splitOnN(word []rune, n int) []string {
	res := make([]string, 0, len(word) / n + 1)
	for i := 0; i < len(word) - n; i+=n {
		res = append(res, string(word[i:i+n]))
	}
	res = append(res, string(word[(len(word) / n * n):]))
	return res
}
