package main

import (
	"fmt"

	"github.com/naoina/gocchan"
)

type HelloFeature struct{}

func (f *HelloFeature) ActiveIf(context interface{}, options ...interface{}) bool {
	return context.(bool)
}

func (f *HelloFeature) Say1(context interface{}) {
	fmt.Println("Hello Gocchan!")
}

func (f *HelloFeature) Say2(context interface{}) {
	panic("panic")
}

func init() {
	gocchan.AddFeature("hello", &HelloFeature{})
}

func helloWorld1(ctx bool) {
	gocchan.Invoke(ctx, "hello", "Say1", func() {
		fmt.Println("Hello world!")
	})
}

func helloWorld2() {
	gocchan.Invoke(nil, "hello", "Say2", func() {
		fmt.Println("Hello world!")
	})
}

func main() {
	helloWorld1(true)  // print "Hello Gocchan!".
	helloWorld1(false) // print "Hello world!".
	helloWorld2()      // print "Hello world!" but doesn't panicked.
}
