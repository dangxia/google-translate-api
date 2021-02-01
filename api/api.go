package api

import (
	"github.com/dangxia/google-translate-api/ctx"
	"strings"
)

type TranslateApi interface {
	CreateTranslator(text string) (Translator, error)

	CreateLangTranslator(text, lang string) (Translator, error)

	GetCtx() ctx.Context
}

func NewTranslateApi(ctx ctx.Context) TranslateApi {
	return &translateApi{
		Ctx: ctx,
	}
}

type translateApi struct {
	Ctx ctx.Context
}

func (me *translateApi) GetCtx() ctx.Context {
	return me.Ctx
}

func (me *translateApi) CreateTranslator(text string) (Translator, error) {
	return me.CreateLangTranslator(text, me.Ctx.DefaultSourceLang())
}

func (me *translateApi) CreateLangTranslator(text, lang string) (Translator, error) {
	lang = strings.ToLower(lang)
	if err := me.Ctx.IsSupported(lang); err != nil {
		return nil, err
	}

	return &translator{
		ctx:  me.Ctx,
		text: text,
		lang: lang,
	}, nil
}
