package main

import (
	"fmt"
	"secure-password-check/core"
	"secure-password-check/core/corrector"
	"secure-password-check/core/dictionaries"
	"secure-password-check/core/logger"
	"secure-password-check/core/modules/dictionary"
	"secure-password-check/core/modules/entropy"
	"secure-password-check/core/modules/regulars"
	"secure-password-check/core/modules/translit"
	"secure-password-check/core/parser"
)

func main() {
	logger := logger.NewZapEchoLogger()

	config, err := ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	printFunc := printFromConfig(config)

	leakedPasswords, err := dictionaries.NewDictionaryFromFile(config.Dictionaries.LeakedPasswords.LocalDict)
	var dictionaryChecker core.Checker
	if err != nil {
		fmt.Println("WARN: error while reading leaked password dict: dictionary check is disabled!")
		dictionaryChecker = core.NewNoneChecker()
	} else {
		dictionaryChecker = dictionary.NewChecker(
			printFunc,
			leakedPasswords,
		)
	}

	userInputs, err := leakedPasswords.GetWords()
	if err != nil {
		userInputs = nil
	}

	englDict, err := englishDictFromConfig(config)
	if err != nil {
		logger.Fatal(err)
	}

	rusDict, err := russianDictFromConfig(config)
	if err != nil {
		logger.Fatal(err)
	}

	frequencyDict, err := parser.GetDictionaryWithFrequency(config.Dictionaries.Frequency.LocalDict)
	var translitChecker core.Checker
	if err != nil {
		fmt.Printf("WARN: frequensy dictionary parsing error %s: translit check is disabled!\n", err.Error())
		translitChecker = core.NewNoneChecker()
	} else {
		translitChecker = translit.NewChecker(
			printFunc,
			corrector.Corrector{Dict: frequencyDict, Letters: []rune("абвгдежзийклмнопрстуфхцчшщъыьэюя")},
			englDict,
			rusDict,
			config.Translator.MinLengthToCheckDict,
		)
	}

	checkerCfg := checkersConfig{
		regexpChecker: regulars.NewChecker(
			printFunc,
			regulars.Config{MinLength: config.Regulars.MinLength,
				MaxSameSeqenceSymbols: config.Regulars.MaxSameSeqenceSymbols,
				AdditionalRegExprs:    config.Regulars.AdditionalRegexps},
		),
		leakedPasswordsDict: leakedPasswords,
		dictionaryChecker:   dictionaryChecker,
		entropyChecker: entropy.NewChecker(
			printFunc,
			config.Entropy.CrackTimeThreshold,
			userInputs,
		),
		translitChecker: translitChecker,
	}

	result := checkPassword(&checkerCfg, config.Password, printFunc)
	if result {
		fmt.Println("Password is secure!")
	} else {
		fmt.Println("Password is insecure!")
	}
}

func checkPassword(cfg *checkersConfig, password string, printFunc func(string, ...any)) bool {
	printFunc("DEBUG: starting stage 1: regexp checker")
	// stage 1: regexp checker
	if !cfg.regexpChecker.IsSecure(password) {
		return false
	}

	printFunc("DEBUG: starting stage 2: leaked passwords")
	// stage 2: leaked passwords
	if !cfg.dictionaryChecker.IsSecure(password) {
		return false
	}

	printFunc("DEBUG: starting stage 3: calculate crack time using entropy")
	// stage 3: calculate crack time using entropy
	if !cfg.entropyChecker.IsSecure(password) {
		return false
	}

	printFunc("DEBUG: starting stage 4: translit")
	// stage 4: translit
	return cfg.translitChecker.IsSecure(password)
}

type checkersConfig struct {
	regexpChecker       core.Checker
	leakedPasswordsDict dictionaries.Dictionary
	dictionaryChecker   core.Checker
	entropyChecker      core.Checker
	translitChecker     core.Checker
}
