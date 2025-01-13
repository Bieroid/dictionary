package repository

import (
	"strconv"
)

type translator struct {
	language string
	dict map[string][]string
}

var dictionary []translator;

func init() {
	mEng := map[string][]string{"привет": {"hello"}}
	strEng := translator{language: "английский", dict: mEng}
	dictionary = append(dictionary, strEng)
}

func IsLanguangeAviable(language string) (string, bool) {
	lang := ""
	
	for i := 0; i < len(dictionary); i++ {
		if (dictionary[i].language == language) || (strconv.Itoa(i + 1) == language) {
			if strconv.Itoa(i) == language {
				lang = language
			}
			return lang, false
		}
	}
	return lang, true
}

func FindTranslate(language string, word string) []string {
	for i := 0; i < len(dictionary); i++ {
		if (dictionary[i].language == language) || (strconv.Itoa(i + 1) == language) {
			if _, ok := dictionary[i].dict[word]; ok {
				return dictionary[i].dict[word]
			}
		}
	}
	return nil
}

func AddLanguage(lang string) bool {
	if _, ok := IsLanguangeAviable(lang); ok {
		strNew := translator{language: lang, dict: make(map[string][]string)}
		dictionary = append(dictionary, strNew)
		return true
	}
	return false
}

func AddTranslate(lang string, word string, trans string) bool {
	for i := 0; i < len(dictionary); i++ {
		if (dictionary[i].language == lang) || (strconv.Itoa(i + 1) == lang) {
			if value, ok := dictionary[i].dict[trans]; ok {
				for _, seekWord := range value {
					if seekWord == word {
						return false
					}
				}
			}
			dictionary[i].dict[trans] = append(dictionary[i].dict[trans], word)
			break
		}
	}
	return true
}

func GetLanguages() (languages []string) {
	for i := 0; i < len(dictionary); i++ {
		languages = append(languages, dictionary[i].language)
	}
	return languages
}
