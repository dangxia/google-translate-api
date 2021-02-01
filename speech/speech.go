package speech

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/dangxia/google-translate-api/ctx"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const GOOGLE_TTS_MAX_CHARS = 100
const GOOGLE_TTS_RPC = "jQ1olc"

type Speech interface {
	Request() ([]string, error)

	Save(path string) error
}

func NewSpeech(ctx ctx.Context, text, lang string, slowly bool) Speech {
	return &speech{
		ctx:    ctx,
		text:   text,
		lang:   lang,
		slowly: slowly,
	}
}

type speech struct {
	ctx ctx.Context

	text string
	lang string

	slowly bool

	tokens      []string
	tokensOnce  sync.Once
	tokensError error

	payloads   []string
	requestErr error

	requestOnce sync.Once
}

func (me *speech) Save(path string) error {

	payloads, err := me.Request()
	if err != nil {
		return err
	}

	bb := make([]byte, 0)

	for _, pl := range payloads {
		plBB, err := base64.StdEncoding.DecodeString(pl)
		if err != nil {
			return err
		}
		bb = append(bb, plBB...)
	}

	return ioutil.WriteFile(path, bb, 0644)
}

func (me *speech) Request() ([]string, error) {
	me.requestOnce.Do(me.doRequest)
	return me.payloads, me.requestErr
}

func (me *speech) doRequest() {
	tokens, err := me.getTokens()
	if err != nil {
		me.requestErr = err
		return
	}

	results := make([]string, len(tokens))
	errs := make([]error, len(tokens))

	wg := sync.WaitGroup{}
	wg.Add(len(tokens))

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		index := i
		go func() {
			defer wg.Done()
			str, err := me.requestItem(token)
			if err != nil {
				errs[index] = err
			} else {
				results[index] = str
			}
		}()
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			me.requestErr = err
			return
		}
	}
	me.payloads = results
}

func (me *speech) processTokens() {
	me.text = strings.TrimSpace(me.text)

	for _, pp := range me.ctx.PreProcessors() {
		me.text = pp(me.text)
	}

	if len(me.text) < GOOGLE_TTS_MAX_CHARS {
		me.tokens = cleanTokens([]string{me.text})
		return
	}

	me.tokens = me.ctx.Tokenize()(me.text)
	me.tokens = cleanTokens(me.tokens)

	minTokens := make([]string, 0)
	for _, token := range me.tokens {
		splits, err := minimize(token, " ")
		if err != nil {
			me.tokensError = err
			return
		}
		minTokens = append(minTokens, splits...)
	}
	me.tokens = cleanTokens(minTokens)
}

func (me *speech) getTokens() ([]string, error) {
	me.tokensOnce.Do(me.processTokens)
	return me.tokens, me.tokensError
}

func (me *speech) prepareParameters(text string) (string, error) {
	var parameters []interface{}
	if me.slowly {
		parameters = []interface{}{text, me.lang, true, "null"}
	} else {
		parameters = []interface{}{text, me.lang, nil, "null"}
	}
	escaped, err := json.Marshal(parameters)

	if err != nil {
		return "", err
	}

	parameters = []interface{}{
		[]interface{}{
			[]interface{}{
				GOOGLE_TTS_RPC,
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

func (me *speech) requestItem(text string) (string, error) {
	data, err := me.prepareParameters(text)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		me.ctx.TranslateUrl(),
		bytes.NewBuffer([]byte(data)),
	)
	if err != nil {
		return "", err
	}

	me.ctx.DecorateHeader(req.Header)

	response, err := me.ctx.HttpClient().Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}
	result, err := Analyze(string(all))
	if err != nil {
		return "", err
	}
	return result, nil
}
