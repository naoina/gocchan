package main

import (
	"fmt"

	"github.com/naoina/gocchan"
)

type TypesFeature struct{}

func (f *TypesFeature) ActiveIf(context interface{}, options ...interface{}) bool {
	return true
}

func (f *TypesFeature) CallByString(context string) {
	fmt.Println(context)
}

func (f *TypesFeature) CallByInt(context int) {
	fmt.Println(context)
}

type TestStruct struct {
	name string
}

func (f *TypesFeature) CallByStruct(context *TestStruct) {
	fmt.Println(context)
}

func main() {
	gocchan.AddFeature("types", &TypesFeature{})
	gocchan.Invoke("string", "types", "CallByString", nil)              // print "string"
	gocchan.Invoke(1, "types", "CallByInt", nil)                        // print "1"
	gocchan.Invoke(&TestStruct{"struct"}, "types", "CallByStruct", nil) // print "&{struct}"
}
