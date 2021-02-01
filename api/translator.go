package api

import (
	"github.com/dangxia/google-translate-api/ctx"
	"github.com/dangxia/google-translate-api/speech"
	"github.com/dangxia/google-translate-api/translation"
	"strings"
)

type Translator interface {
	Translate() (translation.Translation, error)

	TranslateTo(lang string) (translation.Translation, error)

	Speak() (speech.Speech, error)

	SpeakSlowly(slowly bool) (speech.Speech, error)
}

type translator struct {
	ctx  ctx.Context
	text string
	lang string
}

func (me *translator) Translate() (translation.Translation, error) {
	return me.TranslateTo(me.ctx.DefaultTargetLang())
}

func (me *translator) TranslateTo(lang string) (translation.Translation, error) {
	lang = strings.ToLower(lang)
	if err := me.ctx.IsSupported(lang); err != nil {
		return nil, err
	}
	return translation.NewTranslation(me.lang, lang, me.text, me.ctx), nil
}

func (me *translator) Speak() (speech.Speech, error) {
	return me.SpeakSlowly(me.ctx.DefaultSlowly())
}

func (me *translator) SpeakSlowly(slowly bool) (speech.Speech, error) {
	return speech.NewSpeech(me.ctx, me.text, me.lang, slowly), nil
}
