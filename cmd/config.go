package main

import (
	"flag"
	"fmt"
	"secure-password-check/core/dictionaries"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Verbose     bool `env:"VERBOSE" envDefault:"false"`
	Interactive bool `env:"INTERACTIVE" envDefault:"false"`
	Password    string

	Regulars struct {
		MinLength             int `env:"MIN_PASSWORD_LENGTH" envDefault:"8"`
		MaxSameSeqenceSymbols int // TODO: more regulars
	}

	Entropy struct {
		CrackTimeThreshold float64
	}

	Dictionaries struct {
		English struct {
			UseRemote  bool   `env:"ENGLISH_USE_REMOTE" envDefault:"true"`
			RemoteDict string `env:"ENGLISH_REMOTE_DICTIONARY_URL" envDefault:"https://api.dictionaryapi.dev/api/v2/entries/en/"`
			LocalDict  string `env:"ENGLISH_LOCAL_DICTIONARY_FILENAME" envDefault:"russkiwlst_top_100.lst"`
		}
		Russian struct {
			UseRemote  bool   `env:"RUSSIAN_USE_REMOTE" envDefault:"true"`
			RemoteDict string `env:"RUSSIAN_REMOTE_DICTIONARY_URL" envDefault:"https://dictionary.yandex.net/api/v1/dicservice/lookup?key=dict.1.1.20240301T195227Z.358cc27d6d61c293.57c5635cb7f5ef3d9b43d2e33234e3218b8b83f5&lang=ru-ru&text="`
			LocalDict  string `env:"RUSSIAN_LOCAL_DICTIONARY_FILENAME" envDefault:"russkiwlst_top_100.lst"`
		}
		LeakedPasswords struct {
			LocalDict string `env:"ENGLISH_LOCAL_DICTIONARY_FILENAME" envDefault:"russkiwlst_top_100.lst"`
		}
		Frequency struct {
			LocalDict string `env:"ENGLISH_LOCAL_DICTIONARY_FILENAME" envDefault:"5000words.txt"`
		}
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
			fmt.Printf(str+"\n", args)
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
		dict, err := dictionaries.NewRemoteDictionary(cfg.Dictionaries.Russian.RemoteDict)
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
