package tokenizer

import (
	"strings"
	"testing"
)

func TestTokenizes(t *testing.T) {
	input := "Lorem? Ipsum!"
	expect := []string{"Lorem?", "Ipsum!"}
	testTokenize(t, input, expect, ToneMarksTokenize)

	input = "Hello, it's 24.5 degrees in the U.K. today. $20,000,000."
	expect = []string{"Hello", "it's 24.5 degrees in the U.K. today", "$20,000,000."}
	testTokenize(t, input, expect, PeriodCommaTokenize)

	input = "It's now 6:30 which means: morning missing:space"
	expect = []string{"It's now 6:30 which means", "morning missing", "space"}
	testTokenize(t, input, expect, ColonTokenize)

	input = `¡()[]¿…‥،;—。，、：`
	expect = make([]string, 0)
	for _, v := range input {
		expect = append(expect, string(v))
	}
	testTokenize(t, input, expect, OtherTokenize)

}

func testTokenize(t *testing.T, input string, expect []string, tokenize Tokenize) {
	actual := tokenize(input)

	if neq(expect, actual) {
		t.Fatalf(
			"input: %s, expected: %s, actual: %s",
			input,
			strings.Join(expect, ","),
			strings.Join(actual, ","),
		)
	}
}

func neq(expect, actual []string) bool {
	if len(expect) != len(actual) {
		return true
	}

	for i, v := range expect {
		if actual[i] != v {
			return true
		}
	}
	return false
}
