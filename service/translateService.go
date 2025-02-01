package service

import (
	"errors"
	"strings"
	"github.com/Bieroid/dictionary/dto"
	"github.com/Bieroid/dictionary/repository"
)

const (
	ErrLanguageAlreadyExist = "указанный язык уже существует в словаре"
	ErrWordNotExist         = "язык отсутствует в словаре"
)

type TranslateService struct {
	repo  *repository.Repository
	outer *OuterService
}

func NewTranslateService(repo *repository.Repository, outer *OuterService) *TranslateService {
	translateService := &TranslateService{repo: repo, outer: outer}
	return translateService
}

func (s *TranslateService) Translate(languageName, wordToTranslate string) (translate []string, err error) {
	translate, err = s.repo.FindTranslation(languageName, wordToTranslate)
	if err != nil || translate != nil {
		return
	}

	translationFromServer, err := s.outer.fetchTranslationFromService(languageName, wordToTranslate)
	if err != nil {
		return
	}

	if words := strings.Fields(translationFromServer); len(words) == 2 || len(words) == 1 {
		s.AddTranslationIntoBD(languageName, wordToTranslate, translationFromServer)
		translate = append(translate, translationFromServer)
	}
	
	if translate == nil {
		err = errors.New(ErrTranslationNotFound)
	}

	return
}

func (s *TranslateService) AddLanguage(languageName string) error {
	if s.repo.AddLanguage(languageName) {
		return nil
	}

	err := errors.New(ErrLanguageAlreadyExist)
	return err
}

func (s *TranslateService) AddTranslationIntoBD(languageName, translate, wordToTranslate string) (httpStatusCode int, err error) {
	return s.repo.AddTranslate(languageName, wordToTranslate, translate)
}

func (s *TranslateService) LanguageCall(language string) (languageName string, err error) {
	languageName, ok := s.repo.IsLanguageExists(language)

	if ok {
		err = errors.New(ErrWordNotExist)
	} else {
		err = nil
	}
	return
}

func (s *TranslateService) GetLanguages() (languages []dto.LanguageInfo) {
	languages = s.repo.FetchLanguages()
	return
}
