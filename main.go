package main

import (
	"fmt"
	"github.com/MotyaSS/HypergraphKeyCompression/keycomp"
)

func main() {
	key := []int{730, 705, 637, 593, 572, 571, 130}
	comp, size := keycomp.Compress(key)
	fmt.Printf("original key, size: %v %v\n", len(key)*8, key)
	fmt.Printf("compressed size, key: %v, %v\n", len(comp), comp)
	fmt.Printf("decompressed key: %v\n", keycomp.Decompress(comp, size))

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
