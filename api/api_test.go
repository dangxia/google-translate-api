package api

import (
	"github.com/dangxia/google-translate-api/ctx"
	"strings"
	"testing"
)

func TestAnalyzeResult(t *testing.T) {
	ctx, _ := ctx.NewContext()
	api := NewTranslateApi(ctx)

	translator, err := api.CreateTranslator("hello")
	if err != nil {
		t.Fatal(err)
	}

	translation, err := translator.Translate()
	if err != nil {
		t.Fatal(err)
	}

	s, _ := translation.Get()
	println(strings.Join(s, ","))
}
