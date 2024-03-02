package dictionary

import (
	"secure-password-check/core"
	"secure-password-check/core/dictionaries"
)

type dictionariesChecker struct {
	print func(string, ...any)
	dict  dictionaries.Dictionary
}

func (c *dictionariesChecker) IsSecure(password string) bool {
	if c.dict.IsPresent(password) {
		c.print("ERROR: %s is word from dictionary", password)
		return false
	}
	c.print("INFO: dictionary check is passed")
	return true
}

// dictionary checker looks for password in leaked password dictionary
func NewChecker(print func(string, ...any), dict dictionaries.Dictionary) core.Checker {
	return &dictionariesChecker{dict: dict, print: print}
}
