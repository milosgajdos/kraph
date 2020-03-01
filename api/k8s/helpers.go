package k8s

func stringIn(s string, sx []string) bool {
	for _, v := range sx {
		if v == s {
			return true
		}
	}
	return false
}
