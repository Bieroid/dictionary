package service

import (
	"github.com/Bieroid/dictionary/repository"
	"errors"
)


func Translate(language, word string) (transl []string, err error) {
	transl = repository.FindTranslate(language, word)
	if transl == nil {
		err = errors.New("перевод данного слова отсутствует в словаре")
	}
	return transl, err
}

func AddLanguange(lang string) error {
	if repository.AddLanguage(lang) {
		return nil
	}
	err := errors.New("указанный язык уже существует в словаре")
	return err
}

func AddWord(lang, trans, word string) error {
	if repository.AddTranslate(lang, word, trans) {
		return nil
	}
	err := errors.New("перевод данного слова уже есть в словаре")
	return err
}

func LanguageCall(language string) (lang string, err error) {
	lang, tmblr := repository.IsLanguangeAviable(language)
	
	if tmblr {
		err = errors.New("язык отсутствует в словаре")
	} else {
		err = nil
	}
	return
}

func GetLanguages() (languages []string) {
	languages = repository.GetLanguages()
	return languages
}
