package main

import (
	"fmt"

	//+imports

	"not-edit"
	//+imports-end
)

var ()

//go:generate go-gen-import
//go:generate go fmt main.go
func main() {
	fmt.Println("this is sample")
	log.Println("this is sample by log")
}
