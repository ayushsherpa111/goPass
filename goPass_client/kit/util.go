package kit

func ZipMap(keys []string, data []interface{}, basePatr string) map[string]interface{} {
	mergedMap := make(map[string]interface{})
	for i, v := range keys {
		switch data[i].(type) {
		case *int:
			data[i] = *(data[i].(*int))
			if data[i] == 0 {
				continue
			}
		case *string:
			data[i] = *(data[i].(*string))
			if data[i] == "" {
				if basePatr == "" {
					continue
				}
				data[i] = basePatr
			}
		}
		mergedMap[v] = data[i]
	}
	return mergedMap
}
