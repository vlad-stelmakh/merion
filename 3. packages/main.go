package main

import (
	"github.com/vostelmakh/hello"

	"mathops/ops"
)

func main() {
	println("Hello World")
	println("add: ", ops.Add(4, 5))
	println("div: ", ops.Div(15, 5))
	println(hello.Hello())
}
