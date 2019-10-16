package matcher

import (
	"regexp"
)

type Matcher struct {
	name    string
	rxp     *regexp.Regexp
	onMatch func()
}

func (m Matcher) Count(line []byte) {
	if m.rxp.Match(line) {
		m.onMatch()
	}
}

func NumericMatcher(onMatch func()) Matcher {
	return Matcher{
		name:    "Match numeric",
		rxp:     regexp.MustCompile(`^\d+$`),
		onMatch: onMatch,
	}
}

func AlphaMatcher(onMatch func()) Matcher {
	return Matcher{
		name:    "Match alpha",
		rxp:     regexp.MustCompile(`^[a-zA-Z]+$`),
		onMatch: onMatch,
	}
}

func MixedMatcher(onMatch func()) Matcher {
	return Matcher{
		name:    "Match mixed character",
		rxp:     regexp.MustCompile(`^(\[a-zA-Z]+\d|\d+[a-zA-Z])`),
		onMatch: onMatch,
	}
}
