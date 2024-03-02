package entropy

import (
	"secure-password-check/core"

	"github.com/nbutton23/zxcvbn-go"
)

type entropyChecker struct {
	print              func(string, ...any)
	crackTimeThreshold float64
	userInputs         []string
}

// entropy checker calculates entropy and find out if it's lower than threshold
func NewChecker(print func(string, ...any), crackTimeThreshold float64, userInputs []string) core.Checker {
	return &entropyChecker{print: print, crackTimeThreshold: crackTimeThreshold, userInputs: userInputs}
}

func (c *entropyChecker) IsSecure(password string) bool {
	res := zxcvbn.PasswordStrength(password, c.userInputs)
	c.print("INFO: entropy: %f, crack time: %s", res.Entropy, res.CrackTimeDisplay)
	if res.CrackTime < c.crackTimeThreshold {
		c.print("ERROR: crack time (%fs) is lower that needed (%fs)", res.CrackTime, c.crackTimeThreshold)
		return false
	}
	c.print("INFO: entropy check is passed")
	return true
}
