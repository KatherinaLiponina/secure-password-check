package core

type Checker interface {
	IsSecure(string) bool
}

type multiChecker struct {
	checkers []Checker
}

func NewMultiChecker(checkers ...Checker) Checker {
	return &multiChecker{checkers: checkers}
}

func (c *multiChecker) IsSecure(password string) bool {
	result := true
	for _, checker := range c.checkers {
		secure := checker.IsSecure(password)
		if !secure {
			result = false
		}
	}
	return result
}

type noneChecker struct{}

func NewNoneChecker() Checker {
	return &noneChecker{}
}

func (c *noneChecker) IsSecure(password string) bool {
	return true
}
