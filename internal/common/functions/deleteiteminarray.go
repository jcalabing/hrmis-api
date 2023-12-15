package functions

import (
	"fmt"
	"reflect"
)

func DeleteItemInArray(s []interface{}, itemToDelete interface{}) []interface{} {
	var result []interface{}

	for _, item := range s {
		// Use reflection to compare items
		if !reflect.DeepEqual(item, itemToDelete) {
			// Append only if the items are not equal
			result = append(result, item)
			fmt.Println("Items to Append:", itemToDelete)
		} else {
			fmt.Println("Items are equal:", item, " vs ", itemToDelete)
		}
	}

	return result
}
