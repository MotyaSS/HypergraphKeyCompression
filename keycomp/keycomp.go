package keycomp

import (
	"fmt"
	"strings"
)

func GetRational(nums []int) Rational {

	numPrev2, denPrev2 := 0, 1
	numPrev1, denPrev1 := 1, 0
	numK, denK := 0, 0
	for i, num := range nums {
		numK = num*numPrev1 + numPrev2
		denK = num*denPrev1 + denPrev2

		numPrev2, numPrev1 = numPrev1, numK
		denPrev2, denPrev1 = denPrev1, denK

		fmt.Printf("[%d] %d\n", i, numK)
		fmt.Printf("[%d] %d\n", i, denK)
		fmt.Println()
	}

	return Rational{numPrev1, denPrev1}
}

type Rational struct {
	Num   int
	Denom int
}

func cmp(n, m Rational) int {
	a, b := n.Num*m.Denom, n.Denom*m.Num // TODO need for overflow optimization minmax (use smth like big ints)
	return a - b

}

func formFraction(left, right Rational) Rational {
	return Rational{
		Num:   left.Num + right.Num,
		Denom: left.Denom + right.Denom,
	}
}

func GetPath(num, denom int) string {
	left := Rational{0, 1}
	right := Rational{1, 0}
	cur := Rational{1, 1}
	target := Rational{num, denom}
	path := strings.Builder{}
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
			// how
		}
		// val == 0 is invalid: it should've exited on loop condition
		if val == 0 {
			panic("smthing went wrong while getting path")
		}
	}
	return path.String()
}

func RestoreFromPath(path string) (a, b int) {
	left := Rational{0, 1}
	right := Rational{1, 0}
	cur := Rational{1, 1}

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
