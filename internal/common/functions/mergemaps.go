package functions

func MergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Add elements from the first map
	for key, value := range map1 {
		result[key] = value
	}

	// Add or update elements from the second map
	for key, value := range map2 {
		result[key] = value
	}

	return result
}
