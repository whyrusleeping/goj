package main

import (
	"fmt"
)

//Todo, change arg to 'args' a slice of Equatables
type FunctionCall struct {
	fname string
	arg Equatable
}

func (f *FunctionCall) Value() float64 {
	return 0.0
}


func (f *FunctionCall) Print() string {
	var fmtstr string
	if _,ok := f.arg.(*Term); ok {
		fmtstr = "%s%s"
	} else {
		fmtstr = "%s(%s)"
	}
	return fmt.Sprintf(fmtstr, f.fname, f.arg.Print())
}

func (f *FunctionCall) simple() bool {
	return f.arg.simple()
}

func (f *FunctionCall) ContainsVar(v string) bool {
	return f.arg.ContainsVar(v)
}

func (f *FunctionCall) PrintOrderOfOps() {
	f.arg.PrintOrderOfOps()
	fmt.Printf("%s (X)\n", f.fname)
}
