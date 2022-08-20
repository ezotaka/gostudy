package main

import (
	"fmt"
	"reflect"
	"runtime"
)

func main() {
	run(
		test,
	)
}

func run(fns ...func()) {
	for _, fn := range fns {
		val := reflect.ValueOf(fn)
		name := runtime.FuncForPC(val.Pointer()).Name()
		fmt.Printf("[%s]\n", name)
		fn()
	}
}

func test() {

}
