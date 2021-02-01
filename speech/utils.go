package speech

import (
	"fmt"
	"github.com/dangxia/google-translate-api/tokenizer"
	"regexp"
	"strings"
)

var (
	blank = regexp.MustCompile(fmt.Sprintf("^[%s]*$", regexp.QuoteMeta(tokenizer.ALL_PUNC)))
)

func cleanTokens(tokens []string) []string {
	result := make([]string, 0)
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if blank.MatchString(token) {
			continue
		}
		result = append(result, token)
	}
	return result
}

func minimize(str, delim string) ([]string, error) {
	if len(str) > GOOGLE_TTS_MAX_CHARS {
		index := strings.LastIndex(str[0:GOOGLE_TTS_MAX_CHARS], delim)
		if index == -1 {
			return nil, fmt.Errorf("string %s can't be cut by %s", str, delim)
		}
		remaining, err := minimize(str[index:], delim)
		if err != nil {
			return nil, err
		}
		return append([]string{str[:index]}, remaining...), nil
	} else {
		return []string{str}, nil
	}
}
