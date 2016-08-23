package main

import "fmt"

type BInterface interface {
	call(int)
}

type A struct {
}

type B struct {
}

type C struct {
	A
}

func (this A) call(a int) {
	fmt.Println("Implementation of A: ", a)
}

func (this B) call(a int) {
	fmt.Println("Implememtation of B: ", a)
}

func (this C) call(a int) {
	fmt.Println("Implementation of C: ", a)
}

func inter_func(b BInterface) {
	b.call(42)
}

func main() {
	var a A
	var b B

	inter_func(a)
	inter_func(b)
	fmt.Println("-------------")

	var c C
	inter_func(c)
}
