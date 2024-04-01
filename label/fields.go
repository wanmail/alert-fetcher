package label

import (
	"fmt"
	"strings"
)

const (
	escapeChar    = "\""
	delimiterChar = "."
)

type Index []string

func AppendIndex(i Index, chars ...string) Index {
	for _, c := range chars {
		if c == "" {
			continue
		}
		i = append(i, c)
	}

	return i
}

func ParseIndex(raw string) Index {
	if raw == "" {
		panic("empty index")
	}

	index := Index{}

	for {
		if raw == "" {
			break
		}

		s := strings.SplitN(raw, escapeChar, 3)

		switch len(s) {
		case 0:
			return index
		case 1:
			index = AppendIndex(index, strings.Split(s[0], delimiterChar)...)
			return index
		case 2:
			index = AppendIndex(index, strings.Split(s[0], delimiterChar)...)
			index = AppendIndex(index, s[1])
			return index
		case 3:
			index = AppendIndex(index, strings.Split(s[0], delimiterChar)...)
			index = AppendIndex(index, s[1])
			raw = s[2]
			continue
		}

	}

	return index
}

func FindMap(i Index, input map[string]interface{}) interface{} {
	if len(i) == 0 || input == nil {
		return nil
	}

	val := input[i[0]]

	// last index
	if len(i) == 1 {
		return val
	}

	m, ok := val.(map[string]interface{})

	if !ok {
		// return val
		return nil
	}

	return FindMap(i[1:], m)
}

// func FindStruct(i Index, input interface{}) interface{} {
// 	return nil
// }

type FieldExtractor struct {
	mappings map[string]Index
}

func (f *FieldExtractor) Extract(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range f.mappings {
		result[k] = FindMap(v, input)
	}

	return result
}

func (f *FieldExtractor) ExtractString(input map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for k, v := range f.mappings {
		result[k] = fmt.Sprintf("%v", FindMap(v, input))
	}

	return result
}

func NewFieldExtractor(mappings map[string]string) *FieldExtractor {
	m := make(map[string]Index)

	for k, v := range mappings {
		m[k] = ParseIndex(v)
	}

	return &FieldExtractor{
		mappings: m,
	}
}
