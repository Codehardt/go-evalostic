[![GoDoc](https://godoc.org/github.com/Codehardt/go-evalostic?status.svg)](https://godoc.org/github.com/Codehardt/go-evalostic)

## go-evalostic

`go-evalostic` can be used to evaluate logical string conditions with Golang.

## Usage

The usage of `go-evalotic` is very simple. Just define a list of *conditions* and pass it to the constructor. You will then have 
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
