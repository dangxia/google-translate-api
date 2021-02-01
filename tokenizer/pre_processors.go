package tokenizer

import "fmt"

func CreateToneMarks() (PreProcessor, error) {
	preProcessor, err := NewPreProcessor([]string{TONE_MARKS}, func(s string) string {
		return fmt.Sprintf("%s", s)
	}, func(s string) string {
		return s + " "
	})
	if err != nil {
		return nil, err
	}
	return preProcessor.run, nil
}

func CreateEndOfLine() (PreProcessor, error) {
	preProcessor, err := NewPreProcessor([]string{"-\n"}, func(s string) string {
		return fmt.Sprintf("%s", s)
	}, func(s string) string {
		return ""
	})
	if err != nil {
		return nil, err
	}
	return preProcessor.run, nil
}

func CreateAbbreviations() (PreProcessor, error) {
	preProcessor, err := NewPreProcessor(ABBREVIATIONS, func(s string) string {
		return fmt.Sprintf("%s\\.", s)
	}, func(s string) string {
		return s[0 : len(s)-1]
	})
	if err != nil {
		return nil, err
	}
	return preProcessor.run, nil
}

func CreateWorSub() ([]PreProcessor, error) {
	result := make([]PreProcessor, 0)
	for k, v := range SUB_PAIRS {
		preProcessor, err := NewPreProcessor([]string{k}, func(s string) string {
			return fmt.Sprintf("%s", s)
		}, func(s string) string {
			return v
		})
		if err != nil {
			return nil, err
		}

		result = append(result, preProcessor.run)
	}

	return result, nil
}
