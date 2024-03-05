package translator

func TranslateKeybord(password string) string {
	translator := map[rune]rune{
		'q': 'й',
		'w': 'ц',
		'e': 'у',
		'r': 'к',
		't': 'е',
		'y': 'н',
		'u': 'г',
		'i': 'ш',
		'o': 'щ',
		'p': 'з',
		'a': 'ф',
		's': 'ы',
		'd': 'в',
		'f': 'а',
		'g': 'п',
		'h': 'р',
		'j': 'о',
		'k': 'л',
		'l': 'д',
		'z': 'я',
		'x': 'ч',
		'c': 'с',
		'v': 'м',
		'b': 'и',
		'n': 'т',
		'm': 'ь',
	}

	var translatedPassword string
	for _, symbol := range password {
		translatedPassword += string(translator[symbol])
	}

	return translatedPassword
}

func ReplaceLatinWithCyrillic(password string) string {
	letters := map[rune]rune{
		'a': 'а',
		'b': 'б',
		'c': 'ц',
		'd': 'д',
		'e': 'е',
		'f': 'ф',
		'g': 'г',
		'h': 'х',
		'i': 'и',
		'j': 'й',
		'k': 'к',
		'l': 'л',
		'm': 'м',
		'n': 'н',
		'o': 'о',
		'p': 'п',
		'q': 'к',
		'r': 'р',
		's': 'с',
		't': 'т',
		'u': 'у',
		'v': 'в',
		'w': 'в',
		'x': 'х',
		'y': 'ю',
		'z': 'з',
	}

	diletters := map[string]rune{
		"sh": 'ш',
		"ch": 'ч',
		"zh": 'ж',
		"ya": 'я',
	}

	var translatedPassword string
	symbols := []rune(password)
	for i := 0; i < len(symbols); i++ {
		if i+1 < len(symbols) {
			if sym, ok := diletters[string(symbols[i])+string(symbols[i+1])]; ok {
				translatedPassword += string(sym)
				i++
				continue
			}
		}
		sym, ok := letters[symbols[i]]
		if !ok {
			// unexpected symbol -> drop
			continue
		}
		translatedPassword += string(sym)
	}
	return translatedPassword
}

func TranslateWithSymbolReplacements(password string) string {
	translator := map[rune]rune{
		'4': 'a',
		'@': 'a',
		'8': 'b',
		'3': 'e',
		'9': 'g',
		'1': 'i',
		'7': 'l',
		'0': 'o',
		'5': 's',
		'2': 'z',
	}

	var translatedPassword string
	for _, symbol := range password {
		s, ok := translator[symbol]
		if ok {
			translatedPassword += string(s)
		} else if symbol >= 'a' && symbol <= 'z' || symbol >= 'A' && symbol <= 'Z' {
			translatedPassword += string(symbol)
		}
	}

	return translatedPassword
}
