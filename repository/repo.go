package repository

import (
	"strconv"
)

type translator struct { // сделать ее приватной и переименовать DONE
	language string
	dict map[string][]string // сделать значения слайсом стрингов для синонимов + проверки на дубли DONE
}

var dictionary []translator;

func init() { // переделать инициализацию в 1 строку (сразу с данными) // DONE
	mEng := map[string][]string{"hello": {"привет"}}
	strEng := translator{language: "английский", dict: mEng}
	dictionary = append(dictionary, strEng)
}

func IsLanguangeAviable(language string) (string, bool) { // Возврат DONE
	lang := ""
	
	for i := 0; i < len(dictionary); i++ {
		if (dictionary[i].language == language) || (strconv.Itoa(i + 1) == language) { // Сделать независимость language от регистра и в переменной и в стракте // DONE
			if strconv.Itoa(i) == language {
				lang = language
			}
			return lang, false
		}
	}
	return lang, true
}

func FindTranslate(language string, word string) string { // Разобраться в исправленном коде // DONE
	for i := 0; i < len(dictionary); i++ {
		if (dictionary[i].language == language) || (strconv.Itoa(i + 1) == language) { // Сделать независимость language от регистра и в переменной и в стракте //DONE
			for trans, value := range dictionary[i].dict {
				for _, seekWord := range value {
					if seekWord == word {
						return trans
					}
				}
			}
		}
	}
	return ""
}

func AddLanguage(lang string) bool {
	if _, ok := IsLanguangeAviable(lang); ok {
		strNew := translator{language: lang, dict: make(map[string][]string)} // инициализация в 1 строчку //DONE
		dictionary = append(dictionary, strNew)
		return true
	}
	return false
}

func AddTranslate(lang string, word string, trans string) bool {
	for i := 0; i < len(dictionary); i++ {
		if (dictionary[i].language == lang) || (strconv.Itoa(i + 1) == lang) {
			for _, value := range dictionary[i].dict {
				for _, seekWord := range value {
					if seekWord == word {
						return false
					}
				}
			}
			dictionary[i].dict[trans] = append(dictionary[i].dict[trans], word)
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
