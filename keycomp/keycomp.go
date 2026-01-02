package keycomp

import (
	"bytes"
	"math/big"
)

type Rational struct {
	Num   int
	Denom int
}

// Compress compresses hypergraph encryption algorithm key
func Compress(key []int) (compressed []byte, size int) {
	path := RationalToPath(getRational(shift(key))) // L/R bytes

	return compressPath(path) // 0/1 bits
}

// Decompress decompresses hypergraph encryption algorithm key
func Decompress(compressed []byte, size int) []int {
	path := decompressPath(compressed, size)
	key := getContinuedFraction(PathToRational(path))

	return reverseShift(key)
}

// RationalToPath takes a rational number and returns path to that rational consisting of 'L' and 'R'
func RationalToPath(num, denom int) []byte {
	left := Rational{0, 1}
	right := Rational{1, 0}
	cur := Rational{1, 1}
	target := Rational{num, denom}

	path := bytes.Buffer{}

	for target != cur {
		val := cmp(target, cur)
		if val < 0 {
			right = cur
			cur = mediant(left, cur)
			path.WriteRune('L')
		}
		if val > 0 {
			left = cur
			cur = mediant(cur, right)
			path.WriteRune('R')
		}
	}

	return path.Bytes()
}

// PathToRational takes a bytes sequence consisting of 'L' and 'R' and returns corresponding rational
func PathToRational(path []byte) (num, denom int) {
	left := Rational{0, 1}
	right := Rational{1, 0}
	cur := Rational{1, 1}
	for i := range path {
		if path[i] == 'L' {
			right = cur
			cur = mediant(left, right)
		} else if path[i] == 'R' {
			left = cur
			cur = mediant(left, right)
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
	n1 := big.NewInt(int64(n.Num))
	d1 := big.NewInt(int64(n.Denom))
	n2 := big.NewInt(int64(m.Num))
	d2 := big.NewInt(int64(m.Denom))

	left := new(big.Int).Mul(n1, d2)
	right := new(big.Int).Mul(n2, d1)

	return left.Cmp(right)

}

// mediant returns mediant of two Rational numbers
func mediant(left, right Rational) Rational {
	return Rational{
		Num:   left.Num + right.Num,
		Denom: left.Denom + right.Denom,
	}
}

// compressPath takes a bytes sequence consisting of 'L' and 'R' bytes and returns
// shortened bytes sequence of 0's and 1's (0 representing 'L' and 1 representing 'R')
func compressPath(path []byte) (bytes []byte, size int) {
	res := make([]byte, (len(path)+7)/8)
	for i := range path {
		a := byte(0)
		if path[i] == 'R' {
			a = 1
		}

		res[i/8] ^= a >> (uint8(i%8) + 1)
	}
	return res, len(path)
}

// shift performs +0, +1, +1, +1,..., +2 conversion for nums
func shift(nums []int) []int {
	res := make([]int, len(nums))
	res[0] = nums[0]
	for i := 1; i < len(nums); i++ {
		res[i] = nums[i] + 1
	}
	res[len(res)-1] += 1
	return res
}

// reverseShift performs -0, -1, -1, -1, ..., -2 conversion for nums
func reverseShift(nums []int) []int {
	res := make([]int, len(nums))
	res[0] = nums[0]
	for i := 1; i < len(res); i++ {
		res[i] = nums[i] - 1
	}
	res[len(res)-1] -= 1

	return res
}

// getRational returns rational representation of finite continued fraction
func getRational(nums []int) (num, denom int) {

	numPrev2, denPrev2 := 0, 1
	numPrev1, denPrev1 := 1, 0
	numK, denK := 0, 0
	for _, num := range nums {
		numK = num*numPrev1 + numPrev2
		denK = num*denPrev1 + denPrev2

		numPrev2, numPrev1 = numPrev1, numK
		denPrev2, denPrev1 = denPrev1, denK
	}

	return numPrev1, denPrev1
}

// getContinuedFraction returns representation of rational as continued fraction
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

// decompressPath takes a bytes sequence of 0's and 1's and
// returns corresponding bytes sequence of 'L' and 'R' where 'L' stands for 0 and 'R' stands for 1
func decompressPath(path []byte, size int) []byte {
	res := make([]byte, size)
	for i := 0; i < size; i++ {
		bit := path[i/8] ^ (1 >> uint(i%8))
		if bit == 1 {
			res[i] = 'R'
		} else {
			res[i] = 'L'
		}
	}

	return res
}
