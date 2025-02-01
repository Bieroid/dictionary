package dto

type TranslationRequest struct {
	Language    string `json:"language"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

type LanguageRequest struct {
	Language string `json:"language"`
}

type TranslationResponse struct {
	Message string `json:"message"`
}

type LanguageInfo struct {
	Name string `json:"language"`
	Date string `json:"date"`
}
