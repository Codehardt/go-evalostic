package evalostic

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/cloudflare/ahocorasick"
)

// Evalostic is a matcher that can apply multiple conditions on a string.
type Evalostic struct {
	conditions  []node
	ahoCorasick *ahocorasick.Matcher
	strings     map[string]int
	mapping     map[int][]int // which string can be found in which condition
}

// New builds a new Evalostic matcher that compiles all conditions to one big rule set that can be applied to strings.
func New(conditions []string) (*Evalostic, error) {
	e := Evalostic{
		conditions: make([]node, len(conditions)),
		strings:    make(map[string]int),
		mapping:    make(map[int][]int),
	}
	var stringCounter int
	var allStrings []string
	for i, condition := range conditions {
		root, err := parseCondition(condition)
		if err != nil {
			return nil, fmt.Errorf("condition %d: %s", i, err)
		}
		e.conditions[i] = root
		condStrings := extractStrings(root)
		for _, str := range condStrings {
			strI, ok := e.strings[str]
			if !ok {
				strI = stringCounter
				stringCounter++
				e.strings[str] = strI
				allStrings = append(allStrings, str)
			}
			e.mapping[strI] = append(e.mapping[strI], i)
		}
	}
	e.ahoCorasick = ahocorasick.NewStringMatcher(allStrings)
	return &e, nil
}

func (e *Evalostic) Match(s string) (matchingConditions []int) {
	stringIndices := e.ahoCorasick.Match([]byte(s))
	possibleConditions := make(map[int]struct{})
	for _, strI := range stringIndices {
		for _, conditionI := range e.mapping[strI] {
			possibleConditions[conditionI] = struct{}{}
		}
	}
	for possibleCondition := range possibleConditions {
		if conditionMatches(e.conditions[possibleCondition], s) {
			matchingConditions = append(matchingConditions, possibleCondition)
		}
	}
	sort.Ints(matchingConditions)
	return
}

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

type node interface {
	String() string
	Value() string
	Children() (node, node)
}

type (
	oneSubNode  struct{ node node }
	twoSubNodes struct{ node1, node2 node }
	valueNode   struct {
		nodeValue  string
		nodeValueI int
	}
)

func (_ oneSubNode) Value() string  { return "" }
func (_ twoSubNodes) Value() string { return "" }
func (n valueNode) Value() string   { return n.nodeValue }

func (n oneSubNode) Children() (node, node)  { return n.node, nil }
func (n twoSubNodes) Children() (node, node) { return n.node1, n.node2 }
func (_ valueNode) Children() (node, node)   { return nil, nil }

type (
	nodeAND struct{ twoSubNodes }
	nodeOR  struct{ twoSubNodes }
	nodeNOT struct{ oneSubNode }
	nodeVAL struct{ valueNode }
)

func (n nodeAND) String() string { return fmt.Sprintf("nodeAND{%s,%s}", n.node1, n.node2) }
func (n nodeOR) String() string  { return fmt.Sprintf("nodeOR{%s,%s}", n.node1, n.node2) }
func (n nodeNOT) String() string { return fmt.Sprintf("nodeNOT{%s}", n.node) }
func (n nodeVAL) String() string { return fmt.Sprintf("nodeVAL{%q}", n.nodeValue) }

func parse(tokens []token) (node, error) {
	res := make([]interface{}, len(tokens))
	for i, token := range tokens {
		res[i] = token
	}
	// identify subexpressions with parentheses
	for {
		lPos := findToken(tokens, tokenTypeLPAR)
		if lPos < 0 {
			break
		}
		if lPos >= len(tokens)-1 {
			return nil, errors.New("condition ends with an opening parentheses")
		}
		var (
			rPos   = -1
			offset = lPos + 1
		)
		for {
			if offset >= len(tokens)-1 {
				return nil, errors.New("missing matching closing parentheses")
			}
			rPos = findToken(tokens[offset:], tokenTypeRPAR)
			if rPos < 0 {
				return nil, errors.New("missing matching closing parentheses")
			}
			rPos += offset
			if lPos+1 == rPos {
				return nil, errors.New("empty subexpression found")
			}
			if pos := findToken(tokens[offset:rPos], tokenTypeLPAR); pos < 0 {
				break
			}
			offset = rPos + 1
		}
		subNode, err := parse(tokens[lPos+1 : rPos])
		if err != nil {
			return nil, fmt.Errorf("could not parse subexpression: %s", err)
		}
		res = append(res[:lPos], append([]interface{}{subNode}, res[rPos+1:]...)...)
		tokens = append(tokens[:lPos], append([]token{{tokenType: tokenTypeNONE}}, tokens[rPos+1:]...)...)
	}
	for _, tokenType := range []tokenType{
		tokenTypeVAL,
		tokenTypeNOT,
		tokenTypeAND,
		tokenTypeOR,
	} {
		var tokenFound = true
		for tokenFound {
			tokenFound = false
			for i, elem := range res {
				token, _ := elem.(token)
				if token.tokenType != tokenType {
					continue
				}
				tokenFound = true
				switch tokenType {
				case tokenTypeVAL:
					res[i] = nodeVAL{valueNode{nodeValue: token.matched[1 : len(token.matched)-1]}}
				default:
					if i+1 >= len(res) {
						return nil, fmt.Errorf("missing parameter for %s operator", tokenTypeString[tokenType])
					}
					subNode1, ok := res[i+1].(node)
					if !ok {
						return nil, fmt.Errorf("parameter for %s operator is not a node (1), got: %s", tokenTypeString[tokenType], res[i+1])
					}
					res = append(res[:i+1], res[i+2:]...) // remove the (i+1)th element because it has become a sub node
					switch tokenType {
					case tokenTypeNOT:
						res[i] = nodeNOT{oneSubNode{node: subNode1}}
					default:
						if i == 0 {
							return nil, fmt.Errorf("missing parameter for %s operator", tokenTypeString[tokenType])
						}
						subNode2, ok := res[i-1].(node)
						if !ok {
							return nil, fmt.Errorf("parameter for %s operator is not a node (2)", tokenTypeString[tokenType])
						}
						n := twoSubNodes{subNode2, subNode1}
						switch tokenType {
						case tokenTypeAND:
							res[i] = nodeAND{n}
						case tokenTypeOR:
							res[i] = nodeOR{n}
						default:
							return nil, fmt.Errorf("invalid token type: %s", tokenTypeString[tokenType])
						}
						res = append(res[:i-1], res[i:]...) // remove the (i-1)the element because it has become a sub node
					}
					break
				}
			}
		}
	}
	if len(res) != 1 {
		return nil, errors.New("parse tree must have exactly one start node")
	}
	startNode, ok := res[0].(node)
	if !ok {
		return nil, errors.New("start node is not a node")
	}
	return startNode, nil
}

func conditionMatches(n node, s string) bool {
	switch v := n.(type) {
	case nodeAND:
		return conditionMatches(v.node1, s) && conditionMatches(v.node2, s)
	case nodeOR:
		return conditionMatches(v.node1, s) || conditionMatches(v.node2, s)
	case nodeNOT:
		return !conditionMatches(v.node, s)
	case nodeVAL:
		return strings.Contains(s, v.nodeValue)
	default:
		return false
	}
}

func parseCondition(s string) (node, error) {
	t, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	return parse(t)
}

func extractStrings(n node) []string {
	switch v := n.(type) {
	case nodeVAL:
		return []string{v.nodeValue}
	case nodeAND:
		return append(extractStrings(v.node1), extractStrings(v.node2)...)
	case nodeOR:
		return append(extractStrings(v.node1), extractStrings(v.node2)...)
	default:
		return nil
	}
}
