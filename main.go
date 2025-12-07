package main

import (
	"fmt"
	"github.com/MotyaSS/HypergraphKeyCompression/keycomp"
)

func main() {
	fmt.Println(keycomp.GetPath(13, 44))
	// LLLRRLRL:
	// 1/1
	// 1/2
	// 1/3
	// 1/4
	// 2/7
	// 3/10
	// 5/17
	// 8/27
	// 13/44
}
