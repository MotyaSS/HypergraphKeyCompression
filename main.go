package main

import (
	"fmt"
	"github.com/MotyaSS/HypergraphKeyCompression/keycomp"
)

func getContinuedFraction(num, denom int) []int {
	if num < 0 || denom == 0 {
		panic("denom cannot be 0")
	}

	res := make([]int, 0)
	for denom != 0 { //Euclidean algorithm
		res = append(res, num/denom)
		num, denom = denom, num%denom
	}
	return res
}

func main() {

	comp, size := keycomp.Compress(getContinuedFraction(59, 77))
	fmt.Println(getContinuedFraction(59, 77))
	fmt.Println(keycomp.Decompress(comp, size))

	//fmt.Println(string(keycomp.GetPath(13, 44)))
	//fmt.Println(keycomp.RestoreFromPath([]byte("LLLRRLRL")))
	//
	//fmt.Println(keycomp.RestoreFromPath([]byte("LRLRLLR")))
	//
	//compressed := "LRLLLRLRRRLR"
	//num, den := keycomp.RestoreFromPath([]byte(compressed))
	//fmt.Println(string(keycomp.GetPath(num, den)))
	//fmt.Println(compressed)

}
