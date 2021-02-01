package translation

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	TR_PARAGRAPH_START = `["wrb.fr"`
	TR_PARAGRAPH_END   = `"generic"]`
)

func Analyze(input string) ([]string, error) {
	i := strings.Index(input, TR_PARAGRAPH_START)
	if i == -1 {
		return nil, fmt.Errorf("%s not found", TR_PARAGRAPH_START)
	}
	input = input[i:]

	i = strings.LastIndex(input, TR_PARAGRAPH_END)
	if i == -1 {
		return nil, fmt.Errorf("%s not found", TR_PARAGRAPH_END)
	}
	input = input[:i+len(TR_PARAGRAPH_END)]

	paragraph := make([]interface{}, 0)
	err := json.Unmarshal([]byte(input), &paragraph)
	if err != nil {
		return nil, fmt.Errorf("paragraph json Unmarshal failed, %+v", err)
	}

	if len(paragraph) < 3 {
		return nil, fmt.Errorf("paragraph length < 3")
	}

	innerJson, ok := paragraph[2].(string)
	if !ok {
		return nil, fmt.Errorf("inner json not string")
	}

	section := make([]interface{}, 0)
	err = json.Unmarshal([]byte(innerJson), &section)
	if err != nil {
		return nil, fmt.Errorf("section json Unmarshal failed, %+v", err)
	}

	var tmp interface{} = section
	step := 1

	tmp, err = getByIndex(tmp, 1, step)
	if err != nil {
		return nil, err
	}

	step++
	tmp, err = getByIndex(tmp, 0, step)
	if err != nil {
		return nil, err
	}

	step++
	tmp, err = getByIndex(tmp, 0, step)
	if err != nil {
		return nil, err
	}

	step++
	tmp, err = getByIndex(tmp, 5, step)
	if err != nil {
		return nil, err
	}

	step++
	list, ok := tmp.([]interface{})
	if !ok {
		return nil, fmt.Errorf("step%d type is incorrect", step)
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("translation is empty")
	}

	result := make([]string, 0)
	for _, sentences := range list {
		_sentences, ok := sentences.([]interface{})
		if !ok {
			return nil, fmt.Errorf("step%d type is incorrect", step)
		}

		first, ok := _sentences[0].(string)

		if !ok {
			return nil, fmt.Errorf("first translation is not string")
		}
		result = append(result, first)
	}
	return result, nil
}

func getByIndex(prevStepResult interface{}, index, step int) (interface{}, error) {
	list, ok := prevStepResult.([]interface{})
	if !ok {
		return nil, fmt.Errorf("step%d type is incorrect", step)
	}
	if len(list) < index+1 {
		return nil, fmt.Errorf("step%d length < %d", step, index+1)
	}

	return list[index], nil
}
