package translit

import (
	"secure-password-check/core"
	"secure-password-check/core/corrector"
	"secure-password-check/core/dictionaries"
	"secure-password-check/core/split"
	"secure-password-check/core/translator"
)

type translitChecker struct {
	print      func(string, ...any)
	correction corrector.Corrector
	englDict   dictionaries.Dictionary
	rusDict    dictionaries.Dictionary
}

func NewChecker(print func(string, ...any), correction corrector.Corrector,
	englDict dictionaries.Dictionary, rusDict dictionaries.Dictionary) core.Checker {
		return &translitChecker{print: print, correction: correction, englDict: englDict, rusDict: rusDict}
}

func (c *translitChecker) IsSecure(password string) bool {
	parts := split.SplitPassword(password)
	for _, part := range parts {
		if !c.checkBaseWord(part) {
			return false
		}

		evenlowerParts := split.PossibleSplits(part)
		for _, p := range evenlowerParts {
			if !c.checkBaseWord(p) {
				return false
			}
		}
	}

	replaced := translator.TranslateWithSymbolReplacements(password)
	if !c.checkBaseWord(replaced) {
		return false
	}
	evenlowerParts := split.PossibleSplits(replaced)
	for _, p := range evenlowerParts {
		if !c.checkBaseWord(p) {
			return false
		}
	}

	return true
}

func (c *translitChecker) checkBaseWord(word string) bool {
	// case 1: hello
	if c.englDict.IsPresent(word) {
		c.print("Part of password is in an english dictionary dictionary: %s", word)
		return false
	}
	// case 2: ghbdtn
	if c.rusDict.IsPresent(translator.TranslateKeybord(word)) {
		c.print("Keyboard was switch for part: %s", word)
		return false
	}
	// case 3: privet
	replaced := translator.ReplaceLatinWithCyrillic(word)
	corrected := c.correction.Correct(replaced)
	if c.rusDict.IsPresent(corrected) {
		c.print("Part %s was translated as %s and fixed as %s", word, replaced, corrected)
		return false
	}
	return true
}
