package keyterm

import (
	"regexp"
	"strings"
	"sync"
)

var GetRegex = sync.OnceValue(getRgx)

func getRgx() *regexp.Regexp {
	connectorPattern := alternatives(connectors)
	prefixPattern := alternatives(prefixes)
	capitalizedWord := `\p{Lu}[\p{Ll}\p{M}]*`
	pattern := `(?:\b(?:` + prefixPattern + `) )?` +
		capitalizedWord +
		`(?: (?:(?:` + connectorPattern + `) )?` + capitalizedWord + `)*`

	return regexp.MustCompile(pattern)
}

func alternatives(words string) string {
	items := strings.Split(words, ",")

	for index := range items {
		items[index] = regexp.QuoteMeta(items[index])
	}

	return strings.Join(items, "|")
}

const connectors = "de,da,das,do,dos"
const prefixes = "a,as,o,os"
