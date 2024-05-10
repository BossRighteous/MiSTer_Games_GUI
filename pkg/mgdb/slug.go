package mgdb

import (
	"regexp"
	"strings"
)

func SlugifyString(input string) string {
	r := regexp.MustCompile(`(\(.*\))|(\[.*\])|(, \w*)|(\.\w*$)|[^a-z0-9A-Z]`)
	rep := r.ReplaceAllStringFunc(input, func(m string) string {
		return ""
	})
	return strings.ToLower(rep)
}
