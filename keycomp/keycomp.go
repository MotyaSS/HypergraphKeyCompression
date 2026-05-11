package keycomp

import (
	"bytes"
	"math/big"
)

type Rational struct {
	Num   *big.Int
	Denom *big.Int
}

func newRational(n, d *big.Int) Rational {
	return Rational{
		Num:   new(big.Int).Set(n),
		Denom: new(big.Int).Set(d),
	}
}

// Compress compresses hypergraph encryption algorithm key
func Compress(key []int, n uint16, k uint16) []byte {
	num, denom := getRationalBig(transform(key))
	path := RationalToPath(num, denom)
	compBytes, size := compressPath(path)

	res := make([]byte, 5+len(compBytes))
	res[0] = byte(n >> 8)
	res[1] = byte(n)
	res[2] = byte(k >> 8)
	res[3] = byte(k)
	lastBits := size % 8
	if lastBits == 0 && size > 0 {
		lastBits = 8
	}
	res[4] = byte(lastBits)
	copy(res[5:], compBytes)

	return res
}

// Decompress decompresses hypergraph encryption algorithm key
func Decompress(compressed []byte) ([]int, uint16, uint16) {
	n := uint16(compressed[0])<<8 | uint16(compressed[1])
	k := uint16(compressed[2])<<8 | uint16(compressed[3])
	lastBits := int(compressed[4])

	compBytes := compressed[5:]
	size := len(compBytes) * 8
	if size > 0 {
		size -= (8 - lastBits) % 8
	}

	path := decompressPath(compBytes, size)

	num, denom := PathToRational(path)

	cf := getContinuedFractionBig(num, denom)

	return reverseTransform(cf), n, k
}

// getRationalBig returns rational representation of finite continued fraction
func getRationalBig(nums []int) (*big.Int, *big.Int) {
	numPrev2 := big.NewInt(0)
	denPrev2 := big.NewInt(1)
	numPrev1 := big.NewInt(1)
	denPrev1 := big.NewInt(0)

	for _, a := range nums {
		aBig := big.NewInt(int64(a))

		numK := new(big.Int).Add(
			new(big.Int).Mul(aBig, numPrev1),
			numPrev2,
		)

		denK := new(big.Int).Add(
			new(big.Int).Mul(aBig, denPrev1),
			denPrev2,
		)

		numPrev2, numPrev1 = numPrev1, numK
		denPrev2, denPrev1 = denPrev1, denK
	}

	return numPrev1, denPrev1
}

// RationalToPath takes a rational number and returns path to that rational consisting of 'L' and 'R'
func RationalToPath(num, denom *big.Int) []byte {
	left := newRational(big.NewInt(0), big.NewInt(1))
	right := newRational(big.NewInt(1), big.NewInt(0))
	cur := newRational(big.NewInt(1), big.NewInt(1))
	target := newRational(num, denom)

	var path bytes.Buffer

	for cmp(target, cur) != 0 {
		if cmp(target, cur) < 0 {
			right = cur
			cur = mediant(left, cur)
			path.WriteByte('L')
		} else {
			left = cur
			cur = mediant(cur, right)
			path.WriteByte('R')
		}
	}

	return path.Bytes()
}

// PathToRational takes a bytes sequence consisting of 'L' and 'R' and returns corresponding rational
func PathToRational(path []byte) (*big.Int, *big.Int) {
	left := newRational(big.NewInt(0), big.NewInt(1))
	right := newRational(big.NewInt(1), big.NewInt(0))
	cur := newRational(big.NewInt(1), big.NewInt(1))

	for _, p := range path {
		if p == 'L' {
			right = cur
			cur = mediant(left, right)
		} else {
			left = cur
			cur = mediant(left, right)
		}
	}

	return cur.Num, cur.Denom
}

// cmp compares two Rational
// return value is
// >0 if a>b,
// <0 if a<b,
// 0 if a == b
func cmp(a, b Rational) int {
	left := new(big.Int).Mul(a.Num, b.Denom)
	right := new(big.Int).Mul(b.Num, a.Denom)
	return left.Cmp(right)
}

// mediant returns mediant of two Rational numbers
func mediant(left, right Rational) Rational {
	return Rational{
		Num:   new(big.Int).Add(left.Num, right.Num),
		Denom: new(big.Int).Add(left.Denom, right.Denom),
	}
}

// getContinuedFractionBig returns representation of rational as continued fraction
func getContinuedFractionBig(num, denom *big.Int) []int {
	if denom.Sign() == 0 {
		panic("denom cannot be 0")
	}

	n := new(big.Int).Set(num)
	d := new(big.Int).Set(denom)

	res := make([]int, 0)

	for d.Sign() != 0 {
		quot := new(big.Int).Div(n, d)

		res = append(res, int(quot.Int64())) // guaranteed safe

		tmp := new(big.Int).Mod(n, d)
		n, d = d, tmp
	}

	return res
}

// compressPath takes a bytes sequence consisting of 'L' and 'R' bytes and returns
// shortened bytes sequence of 0's and 1's (0 representing 'L' and 1 representing 'R')
func compressPath(path []byte) (bytes []byte, size int) {
	bytes = make([]byte, (len(path)+7)/8)
	for i := range path {
		a := byte(0)
		if path[i] == 'R' {
			a = 1
		}
		bytes[i/8] ^= a << (7 - i%8)
	}
	return bytes, len(path)
}

// decompressPath takes a bytes sequence of 0's and 1's and
// returns corresponding bytes sequence of 'L' and 'R' where 'L' stands for 0 and 'R' stands for 1
func decompressPath(path []byte, size int) []byte {
	res := make([]byte, size)
	for i := 0; i < size; i++ {
		bit := (path[i/8] >> uint(7-i%8)) & 1
		if bit == 0 {
			res[i] = 'L'
		} else {
			res[i] = 'R'
		}
	}
	return res
}

// transform performs conversion for nums
// Example
// [8 6 4 2 1] ->
// [1 (2-1) (4-2) (6-4) (8-6)] ->
// [1 2 3 3 4]
func transform(nums []int) []int {
	res := make([]int, len(nums))
	res[0] = nums[len(nums)-1]
	for i := 1; i < len(nums); i++ {
		res[i] = nums[len(nums)-1-i] - nums[len(nums)-i] + 1
	}
	res[len(res)-1] += 1
	return res
}

// reverseTransform performs reversed conversion for nums
// Example
// [1 2 3 3 4] ->
// [1 1 2 2 2] ->
// [1 (1+1) (1+1+2) (1+1+2+2) (1+1+2+2+2)] ->
// [1 2 4 6 8]
func reverseTransform(nums []int) []int {
	res := make([]int, len(nums))
	res[len(nums)-1] = nums[0]
	for i := 1; i < len(res); i++ {
		res[len(nums)-1-i] = nums[i] + res[len(nums)-i] - 1
	}

	res[0] -= 1
	return res
}
