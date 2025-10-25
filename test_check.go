package main

import (
	"cineverse/utils"
	"fmt"
)

func main() {
	hash := "$2a$10$033r6FNfsKSj0z2HU8UfzuIRMLHPyyG7KH92d4sxu7lAQLTa1wY7q"
	result := utils.CheckPasswordHash(hash, "password")
	fmt.Println(result)
}
