package util

// DeDupe returns a list with duplicated items removed
func DeDupe(original []string) []string {
	items := []string{}
	unique := map[string]bool{}
	for _, val := range original {
		if _, found := unique[val]; !found {
			items = append(items, val)
			unique[val] = true
		}
	}
	return items
}

// Difference returns any items in a that aren't in b
func Difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}
