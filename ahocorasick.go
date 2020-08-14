package evalostic

import (
	"github.com/cloudflare/ahocorasick"
)

// stringMatcher is a wrapper for an Aho-Corasick matcher.
type ahoCorasick struct {
	matcher *ahocorasick.Matcher
}

func (a *ahoCorasick) match(s string) (indices []int) {
	return a.matcher.Match([]byte(s))
}

func newAhoCorasick(strings []string) *ahoCorasick {
	return &ahoCorasick{
		matcher: ahocorasick.NewStringMatcher(strings),
	}
}
