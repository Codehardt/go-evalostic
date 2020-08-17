[![GoDoc](https://godoc.org/github.com/Codehardt/go-evalostic?status.svg)](https://godoc.org/github.com/Codehardt/go-evalostic)
[![Build Status](https://travis-ci.org/Codehardt/go-evalostic.svg?branch=master)](https://travis-ci.org/Codehardt/go-evalostic)
[![Go Report Card](https://goreportcard.com/badge/github.com/Codehardt/go-evalostic)](https://goreportcard.com/report/github.com/Codehardt/go-evalostic)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## go-evalostic

`go-evalostic` can be used to evaluate logical string conditions with Golang.

## Usage

The usage of `go-evalotic` is very simple. Just define a list of **conditions** and pass it to the constructor. You will then have 
access to a matcher that can apply all conditions to a specific string and returns the indices of all matching conditions.

## Condition Syntax

All strings in a condition have to be quoted (they will be unquoted with `strconv.Unquote()`). Multiple strings can be concatenated with the keywords `AND`, `OR` and parentheses `(` `)`. Strings or subconditions can be negated with `NOT`. Per default, strings are case sensitive. To make a string case insensitive, add a `i` after quotes. 

**Example (Case Sensitive)**: `"foo" AND NOT ("bar" OR "baz")`

**Example (Case Insensitive)**: `"foo"i AND NOT ("bar"i OR "baz"i)`

## Code Example

```golang
e, err := evalostic.New([]string{
    `"foo" OR "bar"`,
    `NOT "foo" AND ("bar" OR "baz")`,
    // add more conditions here
})
if err != nil {
    panic(err)
}
e.Match("foo") // returns [0]
e.Match("bar") // returns [0, 1]
e.Match("foobar") // returns [0]
e.Match("baz") // returns [1]
e.Match("qux") // returns nil
```
