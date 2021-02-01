package translation

import (
	"bytes"
	"encoding/json"
	"github.com/dangxia/google-translate-api/ctx"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

const GOOGLE_TRANSLATE_RPC = "MkEWBc"

type Translation interface {
	Get() ([]string, error)
}

func NewTranslation(sourceLang, targetLang, text string, ctx ctx.Context) Translation {
	return &translation{
		sourceLang: sourceLang,
		targetLang: targetLang,
		text:       text,

		ctx: ctx,
	}
}

type translation struct {
	ctx ctx.Context

	sourceLang, targetLang string
	text                   string

	once sync.Once

	result    []string
	resultErr error
}

func (me *translation) prepareParameters() (string, error) {
	parameters := []interface{}{
		[]interface{}{me.text, me.sourceLang, me.targetLang, true},
		[]interface{}{nil},
	}

	escaped, err := json.Marshal(parameters)
	if err != nil {
		return "", err
	}

	parameters = []interface{}{
		[]interface{}{
			[]interface{}{
				GOOGLE_TRANSLATE_RPC,
				string(escaped),
				nil,
				"generic",
			},
		},
	}

	escaped, err = json.Marshal(parameters)
	if err != nil {
		return "", err
	}

	data := "f.req=" + url.QueryEscape(string(escaped))

	return data, nil
}

func (me *translation) get() ([]string, error) {
	data, err := me.prepareParameters()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		me.ctx.TranslateUrl(),
		bytes.NewBuffer([]byte(data)),
	)
	if err != nil {
		return nil, err
	}

	me.ctx.DecorateHeader(req.Header)

	response, err := me.ctx.HttpClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	list, err := Analyze(string(all))
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (me *translation) Get() ([]string, error) {
	me.once.Do(func() {
		me.result, me.resultErr = me.get()
	})
	return me.result, me.resultErr
}
