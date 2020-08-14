[![GoDoc](https://godoc.org/github.com/Codehardt/go-evalostic?status.svg)](https://godoc.org/github.com/Codehardt/go-evalostic)

## go-evalostic

`go-evalostic` can be used to evaluate logical string conditions with Golang.

## Usage

The usage of `go-evalotic` is very simple. Just define a list of **conditions** and pass it to the constructor. You will then have 
access to a matcher that can apply all conditions to a specific string.

```golang
e, err := evalostic.New([]string{
    /* your conditions */
})
if err != nil {
    /* malformed condition */
}
matchingConditions := e.Match(/* your string */)
```

## Example

```golang
e, err := evalostic.New([]string{
    `"foo" OR "bar"`,
    `NOT "foo" AND ("bar" OR "baz")`,
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
