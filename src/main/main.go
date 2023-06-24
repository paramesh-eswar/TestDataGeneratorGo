package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	fmt.Println("=======================")
	fmt.Println("Test Data Generator")
	fmt.Println("=======================")
}

func main() {
	fmt.Println("Test Data Generator tool started ...")
	fmt.Println(gofakeit.Name())
	fmt.Println(gofakeit.Name())
}
