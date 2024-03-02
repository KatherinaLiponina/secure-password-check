package translit

import (
	"secure-password-check/core"
	"secure-password-check/core/corrector"
	"secure-password-check/core/dictionaries"
	"secure-password-check/core/split"
	"secure-password-check/core/translator"
	"secure-password-check/core/utils"
	"strings"
)

type translitChecker struct {
	print                func(string, ...any)
	correction           corrector.Corrector
	englDict             dictionaries.Dictionary
	rusDict              dictionaries.Dictionary
	minLengthToCheckDict int
}

func NewChecker(print func(string, ...any), correction corrector.Corrector,
	englDict dictionaries.Dictionary, rusDict dictionaries.Dictionary,
	minLengthToCheckDict int) core.Checker {
	return &translitChecker{print: print, correction: correction,
		englDict: englDict, rusDict: rusDict, minLengthToCheckDict: minLengthToCheckDict}
}

func (c *translitChecker) IsSecure(password string) bool {
	clearedPassword := strings.ToLower(split.RemovedSymbols(password))
	c.print("DEBUG: clearedPassword: %s", clearedPassword)
	if !c.checkBaseWord(clearedPassword) {
		return false
	}

	replaced := translator.TranslateWithSymbolReplacements(password)
	c.print("DEBUG: replaced symbols: %s", replaced)
	replaced = strings.ToLower(replaced)
	if !c.checkBaseWord(replaced) {
		return false
	}

	clearedPasswordParts := split.PossibleSplits(clearedPassword)
	c.print("DEBUG: clearedPasswordParts: %v", clearedPasswordParts)

	replacedPasswordParts := split.PossibleSplits(replaced)
	c.print("DEBUG: replacedPasswordParts: %v", replacedPasswordParts)

	parts := utils.MergeSortedArrays(clearedPasswordParts, replacedPasswordParts)
	c.print("DEBUG: parts: %v", parts)
	if !c.checkWords(parts) {
		return false
	}
	c.print("INFO: translit check is passed")

	return true
}

func (c *translitChecker) checkBaseWord(word string) bool {
	// case 1: hello
	if !c.checkViaEnglishDictionary(word) {
		return false
	}

	// case 2: ghbdtn
	if !c.checkViaKeyboardTranslation(word) {
		return false
	}

	// case 3: privet
	if !c.checkViaTransliteration(word) {
		return false
	}
	return true
}

func (c *translitChecker) checkWords(words []string) bool {
	for _, word := range words {
		if !c.checkViaEnglishDictionary(word) {
			return false
		}
	}

	for _, word := range words {
		if !c.checkViaKeyboardTranslation(word) {
			return false
		}
	}

	for _, word := range words {
		if !c.checkViaTransliteration(word) {
			return false
		}
	}
	return true
}

func (c *translitChecker) checkViaEnglishDictionary(word string) bool {
	if len([]rune(word)) >= c.minLengthToCheckDict && c.englDict.IsPresent(word) {
		c.print("ERROR: Part of password is in an english dictionary: %s", word)
		return false
	}
	return true
}

func (c *translitChecker) checkViaKeyboardTranslation(word string) bool {
	tranlated := translator.TranslateKeybord(word)
	if len([]rune(tranlated)) >= c.minLengthToCheckDict && c.rusDict.IsPresent(tranlated) {
		c.print("ERROR: Keyboard was switch for part: %s, origin is %s", word, tranlated)
		return false
	}
	return true
}

func (c *translitChecker) checkViaTransliteration(word string) bool {
	replaced := translator.ReplaceLatinWithCyrillic(word)
	corrected := c.correction.Correct(replaced)
	if len([]rune(corrected)) >= c.minLengthToCheckDict && c.rusDict.IsPresent(corrected) {
		c.print("DEBUG: %s was translated as %s and fixed as %s", word, replaced, corrected)
		c.print("ERROR: Part of password is in russian dictionary: %s", corrected)
		return false
	}
	return true
}
