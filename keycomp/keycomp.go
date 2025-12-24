package keycomp

import (
	"bytes"
)

type Rational struct {
	Num   int
	Denom int
}

// Compress compresses hypergraph encryption algorithm key
func Compress(key []int) []byte {
	// TODO
	return []byte{}
}

// Decompress decompresses hypergraph encryption algorithm key
func Decompress(compressed []byte) []int {
	// TODO
	return []int{}
}

// GetPath takes a rational number and returns path to that rational consisting of 'L' and 'R'
func GetPath(num, denom int) []byte {
	left := Rational{0, 1}
	right := Rational{1, 0}
	cur := Rational{1, 1}
	target := Rational{num, denom}
	path := bytes.Buffer{}
	for cur != target {
		val := cmp(target, cur)
		if val < 0 {
			right = cur
			cur = formFraction(left, cur)
			path.WriteRune('L')
		}
		if val > 0 {
			left = cur
			cur = formFraction(cur, right)
			path.WriteRune('R')
		}
	}
	return path.Bytes() // TODO compress to 1 and 0 before return
}

// RestoreFromPath takes a bytes sequence consisting of 'L' and 'R' and returns corresponding rational
func RestoreFromPath(path []byte) (a, b int) {
	left := Rational{0, 1}
	right := Rational{1, 0}
	cur := Rational{1, 1}

	// TODO decompress path from (1/0) => (L/R)

	for i := range path {
		if path[i] == 'L' {
			right = cur
			cur = formFraction(left, right)
		} else if path[i] == 'R' {
			left = cur
			cur = formFraction(left, right)
		} else {
			panic("smthing went wrong while restoring path")
		}
	}

	return cur.Num, cur.Denom
}

// cmp compares two Rational
// return value is
// >0 if a>b,
// <0 if a<b,
// 0 if a == b
func cmp(n, m Rational) int {
	a, b := n.Num*m.Denom, n.Denom*m.Num // TODO might overflow here
	return a - b

}

func formFraction(left, right Rational) Rational {
	return Rational{
		Num:   left.Num + right.Num,
		Denom: left.Denom + right.Denom,
	}
}

// compressPath takes a bytes sequence consisting of 'L' and 'R' bytes and returns
// shortened bytes sequence of 0's and 1's (0 representing 'L' and 1 representing 'R')
func compressPath(s []byte) []byte {
	return []byte{} // TODO
}

// getRational returns rational representation of finite continued fraction
func getRational(nums []int) Rational {
	numPrev2, denPrev2 := 0, 1
	numPrev1, denPrev1 := 1, 0
	numK, denK := 0, 0
	for _, num := range nums {
		numK = num*numPrev1 + numPrev2
		denK = num*denPrev1 + denPrev2

		numPrev2, numPrev1 = numPrev1, numK
		denPrev2, denPrev1 = denPrev1, denK
	}

	return Rational{numPrev1, denPrev1}
}

// decompressPath takes a bytes sequence of 0's and 1's and
// returns corresponding bytes sequence of 'L' and 'R' where 'L' stands for 0 and 'R' stands for 1
func decompressPath(path []byte) []byte {
	// TODO
	return []byte{}
}
