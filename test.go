package main

import "fmt"

type name struct {
	Age int
}

func (n *name)getName(age  int) name {
	n.Age = age
	return *n
}

func main() {
	fmt.Println("working....")
	n := name{31}
	n.getName(35)
	fmt.Print(n)
}