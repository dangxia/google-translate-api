package api

import (
	"github.com/dangxia/google-translate-api/ctx"
	"strings"
	"testing"
)

func TestAnalyzeResult(t *testing.T) {
	ctx, _ := ctx.NewContextWithOption("com.hk", "http://127.0.0.1:4780")
	api := NewTranslateApi(ctx)

	translator, err := api.CreateTranslator("hello")
	if err != nil {
		t.Fatal(err)
	}

	translation, err := translator.Translate()
	if err != nil {
		t.Fatal(err)
	}

	s, err := translation.Get()
	if err != nil {
		t.Fatal(err)
	}
	println(strings.Join(s, ","))
}
