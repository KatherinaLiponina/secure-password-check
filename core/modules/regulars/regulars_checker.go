package regulars

import (
	"fmt"
	"regexp"
	"secure-password-check/core"
	"secure-password-check/core/utils"
)

type Config struct {
	MinLength             int
	MaxSameSeqenceSymbols int // TODO: unsupported now (coundn't think about optimal regexp)
	AdditionalRegExprs    []string
}

type RegularExpressionStore struct {
	store []*regexp.Regexp
}

func NewRegularExpressionStore(AdditionalRegExprs []string) (RegularExpressionStore, error) {
	baseRegexprs := []string{"[0-9]", "[a-z]", "[A-Z]", "[!@#$%^&*]"}
	// TODO: more regulars maybe?

	regexpStore := RegularExpressionStore{store: make([]*regexp.Regexp, 0,
		len(AdditionalRegExprs)+len(baseRegexprs))}
	var err error

	for _, expr := range baseRegexprs {
		regexpStore.store = append(regexpStore.store, regexp.MustCompile(expr))
	}

	for _, expr := range AdditionalRegExprs {
		regexpr, compileErr := regexp.Compile(expr)
		if compileErr != nil {
			err = fmt.Errorf("%s is not valid regexp, %w", expr, err)
			continue
		}
		regexpStore.store = append(regexpStore.store, regexpr)
	}

	return regexpStore, err
}

type regularsChecker struct {
	print    func(string, ...any)
	config   Config
	regexprs RegularExpressionStore
}

// regulars checker uses regulars expression to check password strength
func NewChecker(print func(string, ...any), config Config) core.Checker {
	regexprs, err := NewRegularExpressionStore(config.AdditionalRegExprs)
	if err != nil {
		print("WARN: some of additional regexp were incorrect and therefore skipped: %s", err.Error())
	}
	return &regularsChecker{print: print, config: config, regexprs: regexprs}
}

func (c *regularsChecker) IsSecure(password string) bool {
	return utils.Any(c.checkPasswordLength(password), c.checkPasswordFeatures(password))
}

func (c *regularsChecker) checkPasswordLength(password string) bool {
	if len(password) < c.config.MinLength {
		c.print("ERROR: password's length is less than %d", c.config.MinLength)
		return false
	}
	c.print("INFO: length check is passed")
	return true
}

func (c *regularsChecker) checkPasswordFeatures(password string) bool {
	result := true
	for _, expr := range c.regexprs.store {
		match := expr.FindString(password) //TODO: find or match
		if len(match) == 0 {
			c.print("ERROR: password doesn't match regexpr %s", expr.String())
			result = false
		}
	}
	if result {
		c.print("INFO: regulars check is passed")
	}
	return result
}
