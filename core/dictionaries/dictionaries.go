package dictionaries

import (
	"fmt"
	"net/http"
	"secure-password-check/core/parser"

	"golang.org/x/exp/maps"
)

type Dictionary interface {
	IsPresent(word string) bool
	GetWords() ([]string, error)
}

type localDictionary struct {
	dict map[string]struct{}
}

func NewDictionaryFromFile(filename string) (Dictionary, error) {
	dict, err := parser.GetDictionaryFromFile(filename)
	if err != nil {
		return nil, err
	}
	return &localDictionary{dict: dict}, nil
}

func (d *localDictionary) IsPresent(word string) bool {
	_, has := d.dict[word]
	return has
}

func (d *localDictionary) GetWords() ([]string, error) {
	return maps.Keys(d.dict), nil
}

type remoteDictionary struct {
	dictURL string
	online  bool
}

func NewRemoteDictionary(url string) (Dictionary, error) {
	dict := remoteDictionary{dictURL: url}
	if !dict.IsPresent("hello") {
		return nil, fmt.Errorf("can't establish connection with %s", url)
	}
	return &dict, nil
}

func (d *remoteDictionary) IsPresent(word string) bool {
	requestURL := fmt.Sprintf("%s%s", d.dictURL, word)
	res, err := http.Get(requestURL)
	if err != nil {
		d.online = false
		return false
	}

	if res.StatusCode == http.StatusOK {
		return true
	}
	if res.StatusCode == http.StatusNotFound {
		return false
	}
	return false
}

func (d *remoteDictionary) GetWords() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}
