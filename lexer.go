package evalostic

import (
	"fmt"
	"regexp"
	"strconv"
)

type tokenType int8

const (
	tokenTypeNONE tokenType = iota
	tokenTypeAND
	tokenTypeOR
	tokenTypeNOT
	tokenTypeVAL
	tokenTypeLPAR
	tokenTypeRPAR
)

var tokenTypeString = map[tokenType]string{
	tokenTypeNONE: "NONE",
	tokenTypeAND:  "nodeAND",
	tokenTypeOR:   "nodeOR",
	tokenTypeNOT:  "nodeNOT",
	tokenTypeVAL:  "nodeVAL",
	tokenTypeLPAR: "LPAR",
	tokenTypeRPAR: "RPAR",
}

type tokenDefinition struct {
	tokenType  tokenType
	definition *regexp.Regexp
}

var tokenDefs = []tokenDefinition{
	{tokenTypeNONE, regexp.MustCompile(`^[\s\r\n]+`)},
	{tokenTypeAND, regexp.MustCompile(`^(?i)and`)},
	{tokenTypeOR, regexp.MustCompile(`^(?i)or`)},
	{tokenTypeNOT, regexp.MustCompile(`^(?i)not`)},
	//{tokenTypeVAL, regexp.MustCompile(`^"[^"]*"`)},
	{tokenTypeVAL, regexp.MustCompile(`^"(?:[^"\\]|\\.)*"`)},
	{tokenTypeLPAR, regexp.MustCompile(`^\(`)},
	{tokenTypeRPAR, regexp.MustCompile(`^\)`)},
}

type token struct {
	tokenType tokenType
	matched   string
	pos       int
}

func (t token) String() string {
	return fmt.Sprintf("%s ( %s ) at pos %d", tokenTypeString[t.tokenType], t.matched, t.pos)
}

func tokenize(condition string) (tokens []token, err error) {
	var pos = 1
recognize:
	for len(condition) > 0 {
		for _, tokenDef := range tokenDefs {
			match := tokenDef.definition.FindStringSubmatchIndex(condition)
			if match != nil {
				if tokenDef.tokenType != tokenTypeNONE {
					matched := condition[match[0]:match[1]]
					if tokenDef.tokenType == tokenTypeVAL {
						unquote, err := strconv.Unquote(matched)
						if err != nil {
							return nil, fmt.Errorf("could not unquote %s: %s", matched, err)
						}
						matched = unquote
					}
					tokens = append(tokens, token{
						tokenType: tokenDef.tokenType,
						matched:   matched,
						pos:       pos + match[0],
					})
				}
				pos += match[1]
				condition = condition[match[1]:]
				goto recognize
			}
		}
		return nil, fmt.Errorf("unexpected token in condition at position %s", condition)
	}
	return
}

func findToken(tokens []token, tokenType tokenType) int {
	for i, token := range tokens {
		if token.tokenType == tokenType {
			return i
		}
	}
	return -1
}
