package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"github.com/Bieroid/dictionary/dto"
)

type OuterService struct {
	client  			*http.Client
	baseUrl 			*url.URL
	retries 	int
}

const (
	ErrTranslationNotFound = "перевод слова отсутствует"
)

func NewOuterService(scheme string, host string, path string, maxRequestRetries int) *OuterService {
	req := &OuterService{
		client: &http.Client{},
		baseUrl: &url.URL{
			Scheme: scheme,
			Host:   host,
			Path:   path,
		},
		retries: maxRequestRetries,
	}
	return req
}

func (req *OuterService) fetchTranslationFromService(languageName string, wordToTranslate string) (translate string, err error) {
	params := url.Values{}
	encodedMessage := "переведи слово " + wordToTranslate + " на язык " + languageName + " и дай ответ одним словом на языке " + languageName
	params.Add("answer", encodedMessage)
	req.baseUrl.RawQuery = params.Encode()
	path := req.baseUrl.String()
	err = errors.New(ErrTranslationNotFound)

	for i := 0; i < req.retries; i++ {
		resp, ok := http.Get(path)
		if ok != nil {
			return
		}
		defer resp.Body.Close()

		body, ok := io.ReadAll(resp.Body)
		if ok != nil {
			return
		}

		var response dto.TranslationResponse
		if ok = json.Unmarshal(body, &response); ok != nil {
			return
		}
		
		translate = trimString(response.Message)

		if len(strings.Fields(translate)) == 1 {
			return translate, nil
		}
	}

	return
}

func trimString(response string) (translate string) {
	replacer := strings.NewReplacer(`"`, "", `\`, "", `.`, "")
	return strings.ToLower(replacer.Replace(response))
}
