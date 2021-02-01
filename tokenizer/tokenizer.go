package tokenizer

import (
	"regexp"
	"strings"
)

var (
	TokenizeSeq = []Tokenize{
		BreakTokenize,
		ToneMarksTokenize,
		PeriodCommaTokenize,
		ColonTokenize,
		OtherTokenize,
	}
)

type Tokenize func(text string) []string

func TotalTokenize(text string) []string {
	seq := []string{text}
	for _, t := range TokenizeSeq {
		nextLoop := make([]string, 0)
		for _, item := range seq {
			strs := t(item)
			for _, str := range strs {
				nextLoop = append(nextLoop, str)
			}
		}
		seq = nextLoop
	}
	return seq
}

func BreakTokenize(text string) []string {
	strs := strings.Split(text, "\n")
	result := make([]string, 0)

	for _, str := range strs {
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			continue
		}
		result = append(result, str)
	}

	return result
}

func ToneMarksTokenize(text string) []string {
	r := regexp.MustCompile(TONE_MARKS)
	text = r.ReplaceAllStringFunc(text, func(s string) string {
		return s + "\n"
	})
	return BreakTokenize(text)
}

func PeriodCommaTokenize(text string) []string {
	r := regexp.MustCompile("(\\.[a-zA-Z])?[,.] ")
	text = r.ReplaceAllStringFunc(text, func(s string) string {
		if len(s) > 2 {
			return s
		}
		return "\n"
	})
	return BreakTokenize(text)
}

func ColonTokenize(text string) []string {
	r := regexp.MustCompile("(\\d)?:")
	text = r.ReplaceAllStringFunc(text, func(s string) string {
		if len(s) > 1 {
			return s
		}
		return "\n"
	})
	return BreakTokenize(text)
}

func OtherTokenize(text string) []string {
	r := regexp.MustCompile("[" + regexp.QuoteMeta(OTHER_PUNC) + "]")
	text = r.ReplaceAllStringFunc(text, func(s string) string {
		return s + "\n"
	})
	return BreakTokenize(text)
}
