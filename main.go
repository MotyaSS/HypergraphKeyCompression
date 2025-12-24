package main

import (
	"fmt"
	"github.com/MotyaSS/HypergraphKeyCompression/keycomp"
)

func main() {
	//keycomp.GetRational([]int{4, 2, 6, 7})
	fmt.Println(string(keycomp.GetPath(13, 44)))
	fmt.Println(keycomp.RestoreFromPath([]byte("LLLRRLRL")))

	fmt.Println(keycomp.RestoreFromPath([]byte("LRLRLLR")))

	compressed := "LRLLLRLRRRLR"
	num, den := keycomp.RestoreFromPath([]byte(compressed))
	fmt.Println(string(keycomp.GetPath(num, den)))
	fmt.Println(compressed)

}
