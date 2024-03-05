package main

import (
	"flag"
	"fmt"
	"secure-password-check/core/dictionaries"
	"strings"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Verbose     bool `env:"VERBOSE" envDefault:"false"`
	Interactive bool `env:"INTERACTIVE" envDefault:"false"`
	Password    string

	LevelsOfLogging []string `env:"LEVELS_OF_LOGGING" envSeparator:"," envDefault:"ERROR,WARN"`

	Regulars struct {
		MinLength             int      `env:"MIN_PASSWORD_LENGTH" envDefault:"8"`
		MaxSameSeqenceSymbols int      `env:"MAX_SAME_LENGTH_SEQENCE_SYMBOL" envDefault:"3"`
		AdditionalRegexps     []string `env:"ADDITIONAL_REGEXPRS" envSeparator:","`
	}

	Entropy struct {
		CrackTimeThreshold float64 `env:"CRACK_TIME_THRESHOLD" envDefault:"3600"`
	}

	Dictionaries struct {
		English struct {
			UseRemote  bool   `env:"ENGLISH_USE_REMOTE" envDefault:"false"`
			RemoteDict string `env:"ENGLISH_REMOTE_DICTIONARY_URL" envDefault:"https://api.dictionaryapi.dev/api/v2/entries/en/"`
			LocalDict  string `env:"ENGLISH_LOCAL_DICTIONARY_FILENAME" envDefault:"words.txt"`
		}
		Russian struct {
			UseRemote  bool   `env:"RUSSIAN_USE_REMOTE" envDefault:"false"`
			RemoteDict string `env:"RUSSIAN_REMOTE_DICTIONARY_URL" envDefault:"https://dictionary.yandex.net/api/v1/dicservice/lookup?key=dict.1.1.20240301T195227Z.358cc27d6d61c293.57c5635cb7f5ef3d9b43d2e33234e3218b8b83f5&lang=ru-ru&text="`
			LocalDict  string `env:"RUSSIAN_LOCAL_DICTIONARY_FILENAME" envDefault:"russian_new.txt"`
		}
		LeakedPasswords struct {
			LocalDict string `env:"LEAKED_PASSWORDS_DICTIONARY_FILENAME" envDefault:"russkiwlst_top_100.lst"`
		}
		Frequency struct {
			LocalDict string `env:"RUSSIAN_FREQUENCY_DICTIONARY_FILENAME" envDefault:"5000words.txt"`
		}
	}

	Translator struct {
		MinLengthToCheckDict int `env:"MIN_LENGTH_TO_CHECK_DICTIONARY" envDefault:"5"`
	}
}

func ParseConfig() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	verbose := flag.Bool("verbose", false, "print password check errors")
	interactive := flag.Bool("interactive", false, "TODO: add usage")
	// TODO: more flags!
	flag.Parse()

	config.Verbose = *verbose
	config.Interactive = *interactive
	args := flag.Args()
	if len(args) < 1 && !config.Interactive {
		return &config, fmt.Errorf("no password provided")
	}
	if !config.Interactive {
		config.Password = args[0]
	}

	return &config, nil
}

func printFromConfig(cfg *Config) func(string, ...any) {
	if cfg.Verbose {
		return func(str string, args ...any) {
			log := false
			for _, level := range cfg.LevelsOfLogging {
				if strings.HasPrefix(str, level) {
					log = true
					break
				}
			}
			if log {
				fmt.Printf(str+"\n", args...)
			}
		}
	}
	return func(s string, a ...any) {}
}

func englishDictFromConfig(cfg *Config) (dictionaries.Dictionary, error) {
	if cfg.Dictionaries.English.UseRemote {
		dict, err := dictionaries.NewRemoteDictionary(cfg.Dictionaries.English.RemoteDict)
		if err != nil {
			return nil, err
		}
		return dict, nil
	}
	dict, err := dictionaries.NewDictionaryFromFile(cfg.Dictionaries.English.LocalDict)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

func russianDictFromConfig(cfg *Config) (dictionaries.Dictionary, error) {
	if cfg.Dictionaries.Russian.UseRemote {
		dict, err := dictionaries.NewYandexAPIDictionary(cfg.Dictionaries.Russian.RemoteDict)
		if err != nil {
			return nil, err
		}
		return dict, nil
	}
	dict, err := dictionaries.NewDictionaryFromFile(cfg.Dictionaries.Russian.LocalDict)
	if err != nil {
		return nil, err
	}
	return dict, nil
}
