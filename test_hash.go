package main

import (
	"cineverse/utils"
	"fmt"
)

func main() {
	hash, _ := utils.HashPassword("password")
	fmt.Println(hash)
}
