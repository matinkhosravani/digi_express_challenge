package openweather

import (
	"encoding/xml"
	"log"
	"strconv"
	"strings"
)

func unmarshalXMLToMap(data []byte) (map[string]interface{}, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	stack := make([]map[string]interface{}, 0)
	result := make(map[string]interface{})

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			elementMap := make(map[string]interface{})
			for _, attr := range se.Attr {
				elementMap[attr.Name.Local] = attr.Value
			}

			stack = append(stack, result)
			result[se.Name.Local] = elementMap
			result = elementMap
		case xml.CharData:
			if len(strings.TrimSpace(string(se))) > 0 {
				result["_value"] = strings.TrimSpace(string(se))
			}
		case xml.EndElement:
			if len(stack) > 0 {
				result = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		}
	}

	normalizeMap(result)

	return result, nil
}

func normalizeMap(data map[string]interface{}) {
	if len(data) == 1 && data["_value"] != nil {
		data["_value"] = strings.TrimSpace(data["_value"].(string))
		delete(data, "_value")
	}

	for key, value := range data {
		if child, ok := value.(map[string]interface{}); ok {
			if len(child) == 1 && child["_value"] != nil {
				data[key] = child["_value"]
			} else {
				normalizeMap(child)
			}
		}
	}
}

func getPropertyValue(data map[string]interface{}, property string) (interface{}, bool) {
	keys := strings.Split(property, ".")
	if len(keys) == 0 {
		return nil, false
	}

	value, ok := data[keys[0]]
	if !ok {
		return nil, false
	}

	if len(keys) == 1 {
		return value, true
	}

	if nestedSlice, ok := value.([]interface{}); ok {

		index, err := strconv.Atoi(keys[1])
		if err != nil {
			log.Fatal(err.Error())
		}
		return getNestedValue(nestedSlice[index].(map[string]interface{}), strings.Join(keys[2:], "."))
	}
	if nestedMap, ok := value.(map[string]interface{}); ok {
		return getNestedValue(nestedMap, strings.Join(keys[1:], "."))
	}

	return nil, false
}

func getNestedValue(m map[string]interface{}, key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return nil, false
	}

	value, ok := m[keys[0]]
	if !ok {
		return nil, false
	}

	if len(keys) == 1 {
		return value, true
	}

	if nestedMap, ok := value.(map[string]interface{}); ok {
		return getNestedValue(nestedMap, strings.Join(keys[1:], "."))
	}

	return nil, false
}
