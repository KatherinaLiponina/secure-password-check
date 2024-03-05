package corrector

import (
	"strings"

	"golang.org/x/exp/maps"
)

type Corrector struct {
	Dict    map[string]float64
	Letters []rune
}

func (c *Corrector) Correct(word string) string {
	word = strings.ToLower(word)
	return c.correction(word)
}

func (c *Corrector) correction(word string) string {
	candidates := c.candidates(word)
	var answer string
	var max float64
	for _, cand := range candidates { // TODO: refactor
		if c.Dict[cand] > max {
			answer = cand
			max = c.Dict[cand]
		}
	}
	return answer
}

func (c *Corrector) candidates(word string) []string {
	return append(c.known([]string{word}), c.known(c.edits1(word))...)
}

func (c *Corrector) known(words []string) []string {
	knownWords := make([]string, 0)
	for _, word := range words {
		if _, ok := c.Dict[word]; ok {
			knownWords = append(knownWords, word)
		}
	}
	return knownWords
}

func (c *Corrector) edits1(word string) []string {
	letters := c.Letters
	k := len(letters)

	runedWord := []rune(word)

	splits := make([]split, 0, len(runedWord))
	for i := 0; i < len(runedWord)+1; i++ {
		splits = append(splits, split{left: runedWord[:i], right: []rune(word)[i:]})
	}

	set := make(map[string]struct{}, 2*(k+1)*len(runedWord)+k-1)

	// deletes
	for i := 0; i < len(splits)-1; i++ {
		set[string(splits[i].left)+string(splits[i].right[1:])] = struct{}{}
	}

	// transposes
	for i := 0; i < len(splits) && len(splits[i].right) > 1; i++ {
		set[string(splits[i].left)+string(splits[i].right[1])+
			string(splits[i].right[0])+string(splits[i].right[2:])] = struct{}{}
	}

	// replaces
	for i := 0; i < len(splits)-1; i++ {
		for _, letter := range letters {
			set[string(splits[i].left)+string(letter)+string(splits[i].right[1:])] = struct{}{}
		}
	}

	// inserts
	for _, sp := range splits {
		for _, letter := range letters {
			set[string(sp.left)+string(letter)+string(sp.right)] = struct{}{}
		}
	}

	return maps.Keys(set)
}

type split struct {
	left  []rune
	right []rune
}

// func (c *Corrector) edits2(word string) []string {
// 	res := make([]string, 0, 100000)
// 	edits := c.edits1(word)
// 	for _, e1 := range edits {
// 		e2 := c.edits1(e1)
// 		res = append(res, e2...)
// 	}
// 	return res
// }
