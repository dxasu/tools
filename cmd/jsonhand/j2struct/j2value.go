package j2struct

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

func Extend(iresult any) any {
	switch iresult := iresult.(type) {
	case map[any]any:
		ret := make(map[string]any, len(iresult))
		for k, v := range iresult {
			ret[fmt.Sprint(k)] = Extend(v)
		}
		return ret
	case map[string]any:
		ret := make(map[string]any, len(iresult))
		for k, v := range iresult {
			ret[k] = Extend(v)
		}
		return ret
	case []any:
		var ret []any
		for _, v := range iresult {
			ret = append(ret, Extend(v))
		}
		return ret
	case string:
		var obj any
		err := json.Unmarshal([]byte(iresult), &obj)
		if err == nil {
			return Extend(obj)
		}
	default:
	}
	return iresult
}

// Generate a struct definition given a JSON string representation of an object and a name structName.
func ToStructWithValue(input io.Reader, structName string, convertFloats bool) ([]byte, error) {
	var subStructMap map[string]string = nil

	var result map[string]interface{}

	iresult, err := parsejson(input)
	if err != nil {
		return nil, err
	}

	switch iresult := iresult.(type) {
	case map[interface{}]interface{}:
		result = convertKeysToStrings(iresult)
	case map[string]interface{}:
		result = iresult
	case []interface{}:
		src := fmt.Sprintf("%s\n",
			typeForValue(iresult, structName, nil, subStructMap, convertFloats))
		formatted, err := format.Source([]byte(src))
		if err != nil {
			err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
		}
		return formatted, err
	default:
		return nil, fmt.Errorf("unexpected type: %T", iresult)
	}

	src := fmt.Sprintf("%s\n}",
		generateValues(result, structName, 0, subStructMap, convertFloats))

	return []byte(src), err
}

// Generate go struct entries for a map[string]interface{} structure
func generateValues(obj map[string]interface{}, structName string, depth int, subStructMap map[string]string, convertFloats bool) string {
	structure := structName + "{"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	indent := strings.Repeat("    ", depth+1)
	suffix := "\n" + indent + "}"
	for _, key := range keys {
		value := obj[key]
		structName = FmtFieldName(key)
		valueType := value
		switch value := value.(type) {
		case []interface{}:
			if len(value) > 0 {
				sub := ""
				if v, ok := value[0].(map[interface{}]interface{}); ok {
					sub = generateValues(convertKeysToStrings(v), structName, depth+1, subStructMap, convertFloats) + suffix
				} else if v, ok := value[0].(map[string]interface{}); ok {
					sub = generateValues(v, structName, depth+1, subStructMap, convertFloats) + suffix
				}

				if sub != "" {
					valueType = "[]" + sub
				}
			}
		case map[interface{}]interface{}:
			valueType = generateValues(convertKeysToStrings(value), structName, depth+1, subStructMap, convertFloats) + suffix
		case map[string]interface{}:
			valueType = generateValues(value, structName, depth+1, subStructMap, convertFloats) + suffix
		case string:
			valueType = fmt.Sprintf(`"%s"`, value)
		default:
			valueType = cast.ToString(value)
		}

		fieldName := FmtFieldName(key)
		structure += fmt.Sprintf("\n%s%s: %s,",
			indent,
			fieldName,
			valueType)
	}
	structure = structure[:len(structure)-1]
	return structure
}
