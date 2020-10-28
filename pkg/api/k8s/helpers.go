package k8s

// stringIn returns true if string s i in the slice sx
func stringIn(s string, sx []string) bool {
	for _, v := range sx {
		if v == s {
			return true
		}
	}
	return false
}
