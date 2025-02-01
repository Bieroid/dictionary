package comand

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/Bieroid/dictionary/dto"
	"github.com/Bieroid/dictionary/service"
)

const (
	ErrInvalidInput = "некорректный ввод"
	ErrWordTooLong  = "Длина введенного слова превышает допустимое"
	ErrNoLanguages  = "Нет языков в БД"
)

// TranslateHandle представляет обработчик для работы с переводами
type TranslateHandle struct {
	service *service.TranslateService
}

// NewTranslateHandle создает новый экземпляр TranslateHandle
func NewTranslateHandle(service *service.TranslateService) *TranslateHandle {
	translateHandle := &TranslateHandle{service: service}
	return translateHandle
}

// TranslateWord переводит слово на указанный язык
// @Summary Перевести слово
// @Description Переводит слово на указанный язык
// @Tags translate
// @Accept json
// @Produce json
// @Param language query string true "Язык для перевода"
// @Param word query string true "Слово для перевода"
// @Success 200 {object} map[string][]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /dictionary/api/translate [get]
func (h *TranslateHandle) TranslateWord(c echo.Context) error {
	languageName := c.QueryParam("language")
	wordToTranslate := c.QueryParam("word")

	if languageName == "" || wordToTranslate == "" {
		return c.JSON(http.StatusBadRequest, errorMessageForm(ErrInvalidInput))
	}

	translations, err := h.service.Translate(languageName, wordToTranslate)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorMessageForm(err.Error()))
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, map[string][]string{"translations": translations})
}

// AddTranslation добавляет новый перевод в базу данных
// @Summary Добавить перевод
// @Description Добавляет новый перевод в базу данных
// @Tags translate
// @Accept json
// @Produce json
// @Param input body dto.TranslationRequest true "Данные для добавления перевода"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /dictionary/api/words [post]
func (h *TranslateHandle) AddTranslation(c echo.Context) error {
	var input dto.TranslationRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorMessageForm(ErrInvalidInput))
	}

	if input.Language == "" || input.Word == "" || input.Translation == "" {
		return c.JSON(http.StatusBadRequest, errorMessageForm(ErrInvalidInput))
	}

	if len(input.Word) > 255 {
		return c.JSON(http.StatusBadRequest, errorMessageForm(ErrWordTooLong))
	}

	errID, err := h.service.AddTranslationIntoBD(input.Language, input.Word, input.Translation)
	if err != nil {
		return c.JSON(errID, errorMessageForm(err.Error()))
	}

	return c.JSON(http.StatusCreated, statusMessageForm("Слово '"+input.Word+"' успешно добавлено в язык '"+input.Language+"'"))
}

// AddLanguage добавляет новый язык в базу данных
// @Summary Добавить язык
// @Description Добавляет новый язык в базу данных
// @Tags languages
// @Accept json
// @Produce json
// @Param input body dto.LanguageRequest true "Данные для добавления языка"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /dictionary/api/languages [post]
func (h *TranslateHandle) AddLanguage(c echo.Context) error {
	var input dto.LanguageRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorMessageForm(ErrInvalidInput))
	}

	if input.Language == "" {
		return c.JSON(http.StatusBadRequest, errorMessageForm(ErrInvalidInput))
	}

	err := h.service.AddLanguage(input.Language)
	if err != nil {
		return c.JSON(http.StatusConflict, errorMessageForm(err.Error()))
	}

	return c.JSON(http.StatusCreated, statusMessageForm("Язык '"+input.Language+"' успешно добавлен"))
}

// GetLanguages возвращает список доступных языков
// @Summary Получить языки
// @Description Возвращает список языков, доступных в базе данных
// @Tags languages
// @Produce json
// @Success 200 {object} map[string][]dto.LanguageInfo
// @Failure 404 {object} map[string]string
// @Router /dictionary/api/languages [get]
func (h *TranslateHandle) GetLanguages(c echo.Context) error {
	languages := h.service.GetLanguages()
	if languages == nil {
		return c.JSON(http.StatusNotFound, errorMessageForm(ErrNoLanguages))
	}

	return c.JSON(http.StatusOK, map[string][]dto.LanguageInfo{"languages": languages})
}

func errorMessageForm(s string) map[string]string {
	return map[string]string{"error": s}
}

func statusMessageForm(s string) map[string]string {
	return map[string]string{"message": s}
}
