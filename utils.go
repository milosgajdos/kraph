package kraph

import "strings"

func provides(verbs []string, verb string) bool {
	for _, v := range verbs {
		if v == verb {
			return true
		}
	}
	return false
}

func nodeName(kind, name string) string {
	return strings.ToLower(kind) + "-" + strings.ToLower(name)
}
