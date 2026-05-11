package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/MotyaSS/HypergraphKeyCompression/keycomp"
)

func main() {
	const n = 256
	const k = 4
	str, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	tests := strings.Split(string(str), "\n")

	compressions := make([]float64, 0)
	for _, t := range tests {
		testStr := strings.Fields(t)
		key := []int{}
		for _, n := range testStr {
			num, e := strconv.Atoi(n)
			if e != nil {
				log.Fatal(e)
			}

			key = append(key, num)
		}

		slices.SortFunc(key, func(a, b int) int {
			return b - a
		})
		startTime := time.Now()
		comp := keycomp.Compress(key, n, k)
		compTime := time.Since(startTime)
		uncomp_len, comp_len := len(comp), len(key)*4

		fmt.Printf("original size, key: %v %v\n", uncomp_len, key)
		fmt.Printf("compressed size, key: %v %v\n", comp_len, comp)
		decomp, _, _ := keycomp.Decompress(comp)
		decompTime := time.Since(startTime) - compTime
		fmt.Printf("decompressed key: %v\n", decomp)

		compression := float64(comp_len) / float64(uncomp_len)
		compressions = append(compressions, compression)

		fmt.Printf("compression multiplier: %v\n", compression)
		fmt.Printf("decompressed == original: %v\n", reflect.DeepEqual(key, decomp))
		fmt.Printf("compress/decompress time %v %v\n", compTime, decompTime)
		fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	}

	avgCompression := float64(0)
	for _, comp := range compressions {
		avgCompression += comp
	}
	avgCompression /= float64(len(compressions))
	fmt.Printf("avg_compression: %v\n", avgCompression)

}
