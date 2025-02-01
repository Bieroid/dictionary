package repository

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"github.com/Bieroid/dictionary/dto"
)

const (
	ErrNoLanguage               = "язык отсутствует в словаре"
	ErrTranslationAlreadyExists = "указанный перевод уже существует в словаре"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(host string, connStr string) *Repository {
	dataBase, err := sql.Open(host, connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}

	err = dataBase.Ping()
	if err != nil {
		log.Fatal("Ошибка при проверке подключения: ", err)
	}

	translateRepository := &Repository{db: dataBase}
	return translateRepository
}

func (t *Repository) FetchLanguages() (languages []dto.LanguageInfo) {
	rows, err := t.db.Query("SELECT name, data FROM language ORDER BY id")
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var language dto.LanguageInfo
		if err := rows.Scan(&language.Name, &language.Date); err != nil {
			return nil
		}
		languages = append(languages, language)
	}

	return languages
}

func (t *Repository) IsLanguageExists(language string) (string, bool) {
	var languageName string

	err := t.db.QueryRow("SELECT name FROM language WHERE name = $1", language).Scan(&languageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		return "", true
	}

	return languageName, false
}

func (t *Repository) FindTranslation(languageName string, wordToTranslate string) (translations []string, err error) {
	var languageID int

	err = t.db.QueryRow("SELECT id FROM language WHERE name = $1", languageName).Scan(&languageID)
	if err != nil {
		err = errors.New(ErrNoLanguage)
		return
	}

	rows, err := t.db.Query(
		`SELECT translate FROM word JOIN language ON word.language_id = language.id
		 WHERE language.name = $1 AND word.name = $2`, languageName, wordToTranslate)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var translation string
		err = rows.Scan(&translation)
		if err != nil {
			return
		}
		translations = append(translations, translation)
	}

	return
}

func (t *Repository) AddLanguage(languageName string) bool {
	_, tmblr := t.IsLanguageExists(languageName)
	if tmblr {
		return false
	}

	_, err := t.db.Exec("INSERT INTO language (name) VALUES ($1)", languageName)
	return err == nil
}

func (t *Repository) AddTranslate(languageName string, sourceWord string, translation string) (httpStatusCode int, err error) {
	var languageID int

	err = t.db.QueryRow("SELECT id FROM language WHERE name = $1", languageName).Scan(&languageID)
	if err != nil {
		err = errors.New(ErrNoLanguage)
		httpStatusCode = 404
		return
	}

	_, err = t.db.Exec("INSERT INTO word (name, translate, language_id) VALUES ($1, $2, $3)", translation, sourceWord, languageID)
	if err != nil {
		err = errors.New(ErrTranslationAlreadyExists)
		httpStatusCode = 409
	}
	return
}
