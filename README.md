# Gocchan  [![Build Status](https://travis-ci.org/naoina/gocchan.png?branch=master)](https://travis-ci.org/naoina/gocchan)

FeatureToggle library for [golang](http://golang.org/).

Gocchan is inspired by [Chanko](https://github.com/cookpad/chanko).

## Installation

    go get -u github.com/naoina/gocchan

## Usage

Implement the `Feature` interface:

```go
type MyFeature struct {}

func (f *MyFeature) ActiveIf(context interface{}, options ...interface{}) bool {
    return true
}
```

Define the function of Feature:

```go
func (f *MyFeature) ExecMyFeature(context interface{}) {
    // do something.
}
```

Add Feature:

```go
gocchan.AddFeature("name of feature", &MyFeature{})
```

Invoke:

```go
gocchan.Invoke("context", "name of feature", "ExecMyFeature", func() {
    // default processes.
})
```

The function literal passed to 4th argument of the `gocchan.Invoke` is called when any errors occurred in method of Feature.
And also when `ActiveIf` returns `false` is same as above.

See [Godoc](http://godoc.org/github.com/naoina/gocchan) for more docs.

## Example

See `hello.go` in `_example` and/or run it.

    go run hello.go

## License

Gocchan is licensed under the MIT
