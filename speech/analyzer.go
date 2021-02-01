package speech

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	TR_PARAGRAPH_START = `["wrb.fr"`
	TR_PARAGRAPH_END   = `"generic"]`
)

func Analyze(input string) (string, error) {
	i := strings.Index(input, TR_PARAGRAPH_START)
	if i == -1 {
		return "", fmt.Errorf("%s not found", TR_PARAGRAPH_START)
	}
	input = input[i:]

	i = strings.LastIndex(input, TR_PARAGRAPH_END)
	if i == -1 {
		return "", fmt.Errorf("%s not found", TR_PARAGRAPH_END)
	}
	input = input[:i+len(TR_PARAGRAPH_END)]

	paragraph := make([]interface{}, 0)
	err := json.Unmarshal([]byte(input), &paragraph)
	if err != nil {
		return "", fmt.Errorf("paragraph json Unmarshal failed, %+v", err)
	}

	if len(paragraph) < 3 {
		return "", fmt.Errorf("paragraph length < 3")
	}

	innerJson, ok := paragraph[2].(string)
	if !ok {
		return "", fmt.Errorf("inner json not string")
	}

	section := make([]interface{}, 0)
	err = json.Unmarshal([]byte(innerJson), &section)
	if err != nil {
		return "", fmt.Errorf("section json Unmarshal failed, %+v", err)
	}

	if len(section) == 0 {
		return "", fmt.Errorf("section is empty")
	}

	result, ok := section[0].(string)
	if !ok {
		return "", fmt.Errorf("speech base64 is not string")
	}

	return result, nil
}
