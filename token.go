package evalostic

import (
	"fmt"
	"regexp"
)

type tokenType int8

const (
	tokenTypeNONE tokenType = iota
	tokenTypeAND
	tokenTypeOR
	tokenTypeNOT
	tokenTypeID
	tokenTypeVAL
	tokenTypeLPAR
	tokenTypeRPAR
	tokenTypeEQ
	tokenTypeNE
	tokenTypeLT
	tokenTypeLTE
	tokenTypeGT
	tokenTypeGTE
)

var tokenTypeString = map[tokenType]string{
	tokenTypeNONE: "NONE",
	tokenTypeAND:  "AND",
	tokenTypeOR:   "OR",
	tokenTypeNOT:  "NOT",
	tokenTypeVAL:  "VAL",
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
	{tokenTypeID, regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.-]*`)},
	{tokenTypeVAL, regexp.MustCompile(`^"[^"]*"`)},
	{tokenTypeLPAR, regexp.MustCompile(`^\(`)},
	{tokenTypeRPAR, regexp.MustCompile(`^\)`)},
}

type token struct {
	tokenType tokenType
	matched   string
	pos       int
}

func (t token) String() string {
	return fmt.Sprintf("%s[%q]{%d}", tokenTypeString[t.tokenType], t.matched, t.pos)
}

func tokenize(condition string) (tokens []token, err error) {
	var pos = 1
recognize:
	for len(condition) > 0 {
		for _, tokenDef := range tokenDefs {
			match := tokenDef.definition.FindStringSubmatchIndex(condition)
			if match != nil {
				if tokenDef.tokenType != tokenTypeNONE {
					tokens = append(tokens, token{
						tokenType: tokenDef.tokenType,
						matched:   condition[match[0]:match[1]],
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
