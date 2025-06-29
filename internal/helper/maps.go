package helper

// Maps recursively merges multiple maps of type map[string]any.
func Maps(p map[string]any, s ...map[string]any) map[string]any {
	result := make(map[string]any)
	subs := make(map[string][]map[string]any)
	for _, entry := range append([]map[string]any{p}, s...) {
		for key, value := range entry {
			switch value := value.(type) {
			case map[string]any:
				subs[key] = append(subs[key], value)
			default:
				result[key] = value
			}
		}
	}
	for key, sub := range subs {
		if len(sub) == 0 {
			continue
		}
		result[key] = Maps(sub[0], sub[1:]...)
	}
	return result
}
