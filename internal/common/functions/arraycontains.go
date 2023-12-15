package functions

import (
	"reflect"
)

func ArrayContains(s []interface{}, itemtoCompare interface{}) bool {
	for _, item := range s {
		// Use reflection to compare items
		if reflect.DeepEqual(item, itemtoCompare) {
			return true
		}
	}

	return false
}
