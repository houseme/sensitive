package main

import (
	"fmt"

	"github.com/houseme/sensitive"
)

func main() {
	filter := sensitive.New()
	_ = filter.LoadWordDict("../../dict/dict.txt")
	fmt.Println(filter.Replace("xC4x", '*'))
}
