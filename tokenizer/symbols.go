package tokenizer

const (
	TONE_MARKS = "[?!？！]"

	ALL_PUNC = "?!？！.,¡()[]¿…‥،;:—。，、： \t\n\r\v\f"

	OTHER_PUNC = `¡()[]¿…‥،;—。，、：`
)

var (
	ABBREVIATIONS = []string{
		"dr", "jr", "mr",
		"mrs", "ms", "msgr",
		"prof", "sr", "st",
	}

	SUB_PAIRS = map[string]string{
		"Esq.": "Esquire",
	}
)
