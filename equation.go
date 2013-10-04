package main

import (
	"errors"
	"strings"
	"fmt"
	"bytes"
)

func IsOperator(c uint8) bool {
	if c == '+' ||
	c == '-' ||
	c == '*' ||
	c == '/' ||
	c == '^' {
		return true
	}
	return false
}

func IsAlpha(c uint8) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true
	}
	return false
}

func IsNum(c uint8) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

const (
	TLparen = byte(iota)
	TRparen
	TObrace
	TCbrace
	TNumber
	TComma
	TVariable
	TOperator
	TFunction
	TUnknown
)

type Token struct {
	kind byte
	val string
}

func GetToken(segment string) (typ byte, name string) {
	typ = TVariable
	for i := 0; i < len(segment); i++ {
		if segment[i] == '(' {
			typ = TFunction
			name = segment[:i]
			break
		} else if !(IsAlpha(segment[i]) || IsNum(segment[i])) {
			typ = TVariable
			name = segment[:i]
			break
		}
	}
	return
}

func GetOp(segment string) (typ byte, name string) {
	typ = TUnknown
	if IsOperator(segment[0]) {
		typ = TOperator
		name = segment[:1]
	} else if len(segment) > 1 && segment[:2] == ":=" {
		typ = TOperator
		name = ":="
	} else if segment[0] == '=' {
		typ = TOperator
		name = "="
	} else if segment[0] == '{' {
		typ = TObrace
		name = "{"
	} else if segment[0] == '}' {
		typ = TCbrace
		name = "}"
	} else if segment[0] == ',' {
		typ = TComma
		name = ","
	}

	return
}

func Tokenize(input string) []*Token {
	tokens := make([]*Token, 0, len(input))
	buf := new(bytes.Buffer)

	for i := 0; i < len(input); i++ {
		c := input[i]
		fmt.Println(string([]byte{c}))
		if IsNum(c) || c == '.' {
			buf.WriteByte(input[i])
		} else {
			if buf.Len() > 0 {
				tokens = append(tokens,&Token{TNumber, buf.String()})
				buf.Reset()
			}
			t := TUnknown
			if c == '(' {
				t = TLparen
			} else if c == ')' {
				t = TRparen
			} else if IsAlpha(c) {
				typ, name := GetToken(input[i:])
				if typ == TUnknown {
					fmt.Println("Syntax Error!")
					return nil
				}
				i += len(name)
				t = typ
				tokens = append(tokens,&Token{t, name})
				if typ == TFunction {
					tokens = append(tokens,&Token{TOperator, "F"})
					i--
				}
				continue
			} else {
				t, op := GetOp(input[i:])
				if t != TUnknown {
					i += len(op) - 1
					tokens = append(tokens,&Token{t, op})
					continue
				}
			}

			//This is wrong as it assumes that each
			//token is only one character long
			if t != TUnknown {
				tokens = append(tokens,&Token{t,input[i:i+1]})
			}
		}
	}
	if buf.Len() > 0 {
		tokens = append(tokens,&Token{TNumber, buf.String()})
	}
	return tokens
}

//ParseExpression and validate syntax, also expand any 'shortcuts'
func Validate(tokens []*Token) ([]*Token , error) {
	lt := TUnknown
	passtwo := make([]*Token, len(tokens)*2)
	tokc := 0
	for	i := 0; i < len(tokens); i++ {
		switch tokens[i].kind {
		case TLparen:
			passtwo[tokc] = tokens[i]
			tokc++
		case TRparen:
			if lt == TOperator {
				return nil, errors.New("Invalid syntax, Closing Paren cannot follow operator")
			}
			passtwo[tokc] = tokens[i]
			tokc++
		case TOperator:
			if lt == TOperator || lt == TLparen {
				return nil, errors.New("Invalid syntax, improper operator placement.")
			}
			passtwo[tokc] = tokens[i]
			tokc++
		case TVariable, TNumber, TFunction:
			passtwo[tokc] = tokens[i]
			tokc++
		}
		lt = tokens[i].kind
	}
	return passtwo[:tokc], nil
}

//returns true if operator a has higher precedence than b
func comparePrecedence(a, b int) bool {
	if a == b {
		return false
	}

	if a == OFun {
		return false
	}

	if a == OPow {
		return true
	}

	if (a == OMul || a == ODiv) && (b == OAdd || b == OSub) {
		return true
	}

	return false
}

func OpSignToConst(op string) (rop int) {
	switch op {
	case "+":
		rop = OAdd
	case "-":
		rop = OSub
	case "*":
		rop = OMul
	case "/":
		rop = ODiv
	case "^":
		rop = OPow
	case "F":
		rop = OFun
	}
	return rop
}

func GetOperatorString(oper int) string {
	ops := ""
	switch oper {
	case OAdd:
		ops = "+"
	case OSub:
		ops = "-"
	case OMul:
		ops = "*"
	case ODiv:
		ops = "/"
	case OPow:
		ops = "^"
	}
	return ops
}

func build(tokens []*Token) Equatable {
	stack := NewTokStack(len(tokens))
	postfix := NewTokStack(len(tokens))
	for _,t := range tokens {
		switch t.kind {
		case TNumber, TVariable, TFunction:
			postfix.Push(t)
		case TLparen:
			stack.Push(t)
		case TOperator:
			for stack.Size() > 0 && stack.Peek().kind != TLparen {
				if comparePrecedence(OpSignToConst(stack.Peek().val),
						OpSignToConst(t.val)) {
					postfix.Push(stack.Pop())
				} else {
					break
				}
			}
			stack.Push(t)
		case TRparen:
			for stack.Size() > 0 && stack.Peek().kind != TLparen {
				postfix.Push(stack.Pop())
			}
			if stack.Size() > 0 {
				stack.Pop()
			}
		}
	}
	for stack.Size() > 0 {
		postfix.Push(stack.Pop())
	}
	eqs := make([]Equatable, len(postfix.GetSlice()))
	eqsc := 0
	for _,t :=  range postfix.GetSlice() {
		if t.kind == TVariable {
			eqs[eqsc] = NewVariable(t.val)
			eqsc++
		} else if t.kind == TNumber {
			eqs[eqsc] = NewConstant(t.val)
			eqsc++
		} else if t.kind == TFunction {
			eqs[eqsc] = &FunctionCall{t.val,nil}
			eqsc++
		} else if t.kind == TOperator {
			op := OpSignToConst(t.val)
			if op == OFun {
				tpar := eqs[eqsc - 1]
				tfun,_ := eqs[eqsc - 2].(*FunctionCall)
				tfun.arg = tpar
				eqsc--
			} else {
				neq := &Term{eqs[eqsc - 2] ,eqs[eqsc - 1], op}
				eqsc--
				eqs[eqsc - 1] = neq
			}
		}
	}
	return eqs[0]
}

func ParseEquation(input string) (*Equality, error) {
	if !strings.Contains(input,"=") {
		return nil, errors.New("Not a valid equality, must contain '='.")
	}
	spl := strings.Split(input, "=")
	l,lerr := ParseExpression(spl[0])
	r,rerr := ParseExpression(spl[1])
	if lerr != nil {
		return nil, lerr
	}
	if rerr != nil {
		return nil, rerr
	}
	l = Simplify(l)
	r = Simplify(r)
	return &Equality{l,r}, nil
}

func SimilarOp(a, b int) bool {
	if (a == OAdd || a == OSub) && (b == OAdd || b == OSub) {
		return true
	} else if (a == OMul || a == ODiv) && (b == OMul || b == ODiv) {
		return true
	}
	return false
}

func ParseExpression(input string) (Equatable, error) {
	tokens := Tokenize(input)
	if len(tokens) == 1 {
		if tokens[0].kind == TVariable {
			return NewVariable(tokens[0].val),nil
		} else if tokens[0].kind == TNumber {
			return NewConstant(tokens[0].val),nil
		}
	}
	tokens,err := Validate(tokens)
	if err != nil {
		return nil, err
	}
	eq := build(tokens)
	return eq, nil
}

