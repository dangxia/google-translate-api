package ctx

import (
	"fmt"
	"github.com/dangxia/google-translate-api/tokenizer"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; WOW64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/47.0.2526.106 Safari/537.36"

	CONTENT_TYPE = "application/x-www-form-urlencoded;charset=utf-8"

	TRANSLATE_PATH = "/_/TranslateWebserverUi/data/batchexecute"

	DOMAIN_TEMPLATE = "https://translate.google.%s"
)

type Context interface {
	TranslateUrl() string

	CheckLang() bool

	DefaultSourceLang() string
	DefaultTargetLang() string
	DefaultSlowly() bool

	HttpClient() *http.Client
	DecorateHeader(header http.Header)

	IsSupported(lang string) error

	PreProcessors() []tokenizer.PreProcessor
	Tokenize() tokenizer.Tokenize
}

type context struct {
	client *http.Client

	tld string

	checkLang bool

	defaultSourceLang string
	defaultTargetLang string
	defaultSlowly     bool

	translateUrl     string
	refererUrl       string
	translateUrlLock sync.Once

	preProcessors []tokenizer.PreProcessor
	tokenize      tokenizer.Tokenize
}

func createDefaultPreProcessors() ([]tokenizer.PreProcessor, error) {
	defaultPreProcessors := make([]tokenizer.PreProcessor, 0)

	processor, err := tokenizer.CreateToneMarks()
	if err != nil {
		return nil, err
	}
	defaultPreProcessors = append(defaultPreProcessors, processor)

	processor, err = tokenizer.CreateEndOfLine()
	if err != nil {
		return nil, err
	}
	defaultPreProcessors = append(defaultPreProcessors, processor)

	processor, err = tokenizer.CreateAbbreviations()
	if err != nil {
		return nil, err
	}
	defaultPreProcessors = append(defaultPreProcessors, processor)

	processors, err := tokenizer.CreateWorSub()
	if err != nil {
		return nil, err
	}
	for _, processor := range processors {
		defaultPreProcessors = append(defaultPreProcessors, processor)
	}
	return defaultPreProcessors, nil
}

func NewContext() (Context, error) {
	defaultPreProcessors, err := createDefaultPreProcessors()
	if err != nil {
		return nil, err
	}

	netTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	netClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	ctx := &context{
		tld: "cn",

		checkLang: true,

		defaultSourceLang: "en",
		defaultTargetLang: "zh-cn",
		defaultSlowly:     true,

		client: netClient,

		preProcessors: defaultPreProcessors,
		tokenize:      tokenizer.TotalTokenize,
	}

	return ctx, nil
}

func (me *context) TranslateUrl() string {
	me.translateUrlLock.Do(func() {
		me.refererUrl = fmt.Sprintf(DOMAIN_TEMPLATE, me.tld)
		me.translateUrl = me.refererUrl + TRANSLATE_PATH
	})
	return me.translateUrl
}

func (me *context) DefaultSourceLang() string {
	return me.defaultSourceLang
}

func (me *context) DefaultTargetLang() string {
	return me.defaultTargetLang
}

func (me *context) DefaultSlowly() bool {
	return me.defaultSlowly
}

func (me *context) CheckLang() bool {
	return me.checkLang
}

func (me *context) HttpClient() *http.Client {
	return me.client
}

func (me *context) PreProcessors() []tokenizer.PreProcessor {
	return me.preProcessors
}
func (me *context) Tokenize() tokenizer.Tokenize {
	return me.tokenize
}

func (me *context) DecorateHeader(header http.Header) {
	header.Set("Referer", me.refererUrl)
	header.Set("User-Agent", USER_AGENT)
	header.Set("Content-Type", CONTENT_TYPE)
}

func (me *context) IsSupported(lang string) error {
	lang = strings.ToLower(lang)
	if me.checkLang {
		if _, ok := Langs[lang]; !ok {
			return fmt.Errorf("Language not supported: %s ", lang)
		}
	}
	return nil
}
