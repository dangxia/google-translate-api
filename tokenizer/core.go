package tokenizer

import (
	"fmt"
	"regexp"
	"strings"
)

type TokenizerRegexp interface {
	Regexp() *regexp.Regexp

	Pattern() string
}

type tokenizerRegexp struct {
	regexp  *regexp.Regexp
	pattern string
}

func (tr *tokenizerRegexp) Regexp() *regexp.Regexp {
	return tr.regexp
}

func (tr *tokenizerRegexp) Pattern() string {
	return tr.pattern
}

func NewTokenizerRegexp(args []string, patternFunc func(string) string) (TokenizerRegexp, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("args is empty")
	}
	alts := make([]string, 0)
	for _, arg := range args {
		alt := patternFunc(arg)
		alts = append(alts, alt)
	}

	pattern := strings.Join(alts, "|")

	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &tokenizerRegexp{
		regexp:  r,
		pattern: pattern,
	}, nil
}

type PreProcessor = func(text string) string

type PreProcessorRegex struct {
	regexes []*regexp.Regexp
	repl    func(string) string
}

func (pp *PreProcessorRegex) run(text string) string {
	for _, rg := range pp.regexes {
		text = rg.ReplaceAllStringFunc(text, pp.repl)
	}
	return text
}

func NewPreProcessor(searchArgs []string, searchFunc func(string) string, repl func(string) string) (*PreProcessorRegex, error) {
	regexes := make([]*regexp.Regexp, 0)
	for _, arg := range searchArgs {
		tokenizerRegexp, err := NewTokenizerRegexp([]string{arg}, searchFunc)
		if err != nil {
			return nil, err
		}
		regexes = append(regexes, tokenizerRegexp.Regexp())
	}

	return &PreProcessorRegex{
		regexes: regexes,
		repl:    repl,
	}, nil
}
