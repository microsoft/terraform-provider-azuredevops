package datahelper

import (
	"fmt"

	"github.com/ahmetb/go-linq"
)

// GetAttributeValues converts an array of items into an array of one of their properties
func GetAttributeValues(items []interface{}, attributeName string) ([]string, error) {
	var result []string
	for _, element := range items {
		result = append(result, element.(map[string]interface{})[attributeName].(string))
	}
	return result, nil
}

// JoinMap converts a map into a string by a give key/value separator
func JoinMap(permissions map[string]string, keyValueSeparator string, elementSeparator string) string {
	return linq.From(permissions).
		Select(func(i interface{}) interface{} {
			kv := i.(linq.KeyValue)
			return fmt.Sprintf(`%s %s "%s"`, kv.Key, keyValueSeparator, kv.Value)
		}).
		Aggregate(func(r interface{}, i interface{}) interface{} {
			if r.(string) == "" {
				return i
			}
			return r.(string) + elementSeparator + i.(string)
		}).(string)
}
