package main

import (
	"fmt"
	"strings"
)

func main() {

	file := "/web/*filepath"
	fmt.Println(strings.Split(file, "*")[0])
}
