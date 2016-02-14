package translate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	YandexTranslateURL = "https://translate.yandex.net/api/v1.5/tr.json"

	yandexTranslateAPIGetLangs  = "getLangs"
	yandexTranslateAPIDetect    = "detect"
	yandexTranslateAPITranslate = "translate"
)

type Languages struct {
	Directions []string          `json:"dirs"`
	Langs      map[string]string `json:"langs"`
}

func (langs *Languages) CanTranslate(from, to string) bool {
	for _, dir := range langs.Directions {
		components := strings.Split(dir, "-")
		if components[0] == from && components[1] == to {
			return true
		}
	}

	return false
}

type Translation struct {
	Language string   `json:"lang"`
	Text     []string `json:"text"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Translator struct {
	TranslateHTML bool

	client *http.Client
	url    string
	key    string
}

func New(key string) *Translator {
	return &Translator{
		client: &http.Client{},
		url:    YandexTranslateURL,
		key:    key,
	}
}

func (tr *Translator) Languages(ui string) (*Languages, error) {
	values := tr.baseValues()
	if ui != "" {
		values.Set("ui", ui)
	}

	var resp struct {
		Response
		*Languages
	}
	if err := tr.request(yandexTranslateAPIGetLangs, values, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 && resp.Message != "" {
		return nil, NewError(resp.Code, resp.Message)
	}

	return resp.Languages, nil
}

func (tr *Translator) Detect(text string) (string, error) {
	values := tr.baseValues()
	values.Set("text", text)

	var resp struct {
		Response
		Lang string `json:"lang"`
	}
	if err := tr.request(yandexTranslateAPIDetect, values, &resp); err != nil {
		return "", err
	}

	if resp.Code != 200 {
		return "", NewError(resp.Code, resp.Message)
	}

	return resp.Lang, nil
}

func (tr *Translator) Translate(origin, dest, text string) (*Translation, error) {
	if dest == "" {
		return nil, fmt.Errorf("translator: no language destination set")
	}

	values := tr.baseValues()
	values.Set("lang", dest)
	if origin != "" {
		values.Set("lang", origin+"-"+dest)
	}
	values.Set("text", text)

	var resp struct {
		Response
		*Translation
	}
	if err := tr.request(yandexTranslateAPITranslate, values, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, NewError(resp.Code, resp.Message)
	}

	return resp.Translation, nil
}

func (tr *Translator) request(path string, values url.Values, v interface{}) error {
	resp, err := tr.client.PostForm(tr.url+"/"+path, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

func (tr *Translator) baseValues() url.Values {
	return url.Values{"key": {tr.key}}
}

type Error struct {
	Code    int
	Message string
}

func NewError(code int, msg string) *Error {
	return &Error{code, msg}
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Message)
}
