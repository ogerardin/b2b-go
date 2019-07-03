package util

func Map(s []string, f func(string) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func MergeMaps(m1 map[string]interface{}, m2 map[string]interface{}) interface{} {
	for k, v := range m2 {
		m1[k] = v
	}
}
