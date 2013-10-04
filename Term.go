package main

import (
	"math"
	"fmt"
)

type Term struct {
	left, right Equatable
	operator int
}

func (t *Term) Value() float64 {
	switch t.operator {
		case OAdd: return t.left.Value() + t.right.Value()
		case OSub: return t.left.Value() - t.right.Value()
		case OMul: return t.left.Value() * t.right.Value()
		case ODiv: return t.left.Value() / t.right.Value()
		case OPow: return math.Pow(t.left.Value(), t.right.Value())
	}
	return 0
}

func (t *Term) simple() bool {
	return t.left.simple() && t.right.simple()
}

func (t *Term) Print() string{
	ops := GetOperatorString(t.operator)
	return fmt.Sprintf("(%s %s %s)",t.left.Print(), ops, t.right.Print())
}

func (t *Term) ContainsVar(v string) bool {
	return t.left.ContainsVar(v) || t.right.ContainsVar(v)
}

func (t *Term) PrintOrderOfOps() {
	t.left.PrintOrderOfOps()
	t.right.PrintOrderOfOps()
	fmt.Printf("((%s)) %s ((%s))\n", t.left.Print(), GetOperatorString(t.operator), t.right.Print())
}
