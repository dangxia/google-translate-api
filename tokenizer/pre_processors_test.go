package tokenizer

import (
	"testing"
)

func TestPreProcessors(t *testing.T) {
	f, err := CreateToneMarks()
	if err != nil {
		t.Fatal(err)
	}

	input := "lorem!ipsum?"
	expect := "lorem! ipsum? "
	actual := f(input)

	if expect != actual {
		t.Fatalf("input: %s, expected: %s, actual: %s", input, expect, actual)
	}

	f, err = CreateEndOfLine()
	if err != nil {
		t.Fatal(err)
	}

	input = `test-
ing`
	expect = "testing"
	actual = f(input)

	if expect != actual {
		t.Fatalf("input: %s, expected: %s, actual: %s", input, expect, actual)
	}

	f, err = CreateAbbreviations()
	if err != nil {
		t.Fatal(err)
	}

	input = `jr. sr. dr.`
	expect = "jr sr dr"
	actual = f(input)

	if expect != actual {
		t.Fatalf("input: %s, expected: %s, actual: %s", input, expect, actual)
	}

	fs, err := CreateWorSub()
	if err != nil {
		t.Fatal(err)
	}

	input = `Esq. Bacon`
	expect = "Esquire Bacon"

	actual = input
	for _, f := range fs {
		actual = f(actual)
	}

	if expect != actual {
		t.Fatalf("input: %s, expected: %s, actual: %s", input, expect, actual)
	}
}
