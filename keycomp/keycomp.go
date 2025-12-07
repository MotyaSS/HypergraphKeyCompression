package keycomp

import (
	"strings"
)

type rational struct {
	Num   int
	Denom int
}

func cmp(n, m rational) int {
	a, b := n.Num*m.Denom, n.Denom*m.Num // TODO need for overflow optimization minmax (use smth like big ints)
	return a - b

}

func formFraction(left, right rational) rational {
	return rational{
		Num:   left.Num + right.Num,
		Denom: left.Denom + right.Denom,
	}
}

func GetPath(num, denom int) string {
	left := rational{0, 1}
	right := rational{1, 0}
	cur := rational{1, 1}
	target := rational{num, denom}
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
	return -1, -1 // TODO
}
