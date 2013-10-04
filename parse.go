package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	ADD_NUM = iota
	SUB_NUM
	MUL_NUM
	DIV_NUM
)

var Punct []byte = []byte{'(',')',',','{','}','[',']','*','^','+','-','/','=',':',';'}

type Instruction struct {
	Op int
	A *Variable
	B *Variable
	C *Variable
}

type FunctionDef struct {
	name string
	instructions []Instruction
	vars map[string]*Variable
	params []*Variable
	retval []string
}

func main() {
	fi, err := os.Open("test.go.sample")
	if err != nil {
		panic(err)
	}
	br := bufio.NewScanner(fi)
	for br.Scan() {
		if br.Err() != nil {
			break
		}
		toks := Tokenize(br.Text())
		for _,t := range toks {
			fmt.Printf("%s ",t.val)
		}
		fmt.Println()
	}
}

/*
func (f *Function) ParseStatement(toks []string) {
	for _,t := range toks {
		fmt.Printf("%s ",t)
	}
	fmt.Println()
	//First, check if its an assignment
	i := 0
	for ;i < len(toks) && toks[i] != "="; i++ { }
	if i == len(toks) {
		fmt.Println("Unsupported Operation!")
		return
	}

}

func (f *Function) HandleEquation(toks []string) {

}

func ParseFunc(scan *bufio.Scanner) *Function {
	reader := bytes.NewBuffer(scan.Bytes())
	buffer := new(bytes.Buffer)
	fc := new(Function)

	//skip 'func '
	reader.Next(5)
	for {
		b, _ := reader.ReadByte()
		if b == '(' {
			break
		}
		buffer.WriteByte(b)
	}
	fc.name = buffer.String()
	buffer.Reset()
	//Read parameters

	var b byte
	for b != ')' {
		par := new(Variable)
		b = ' '
		for b == ' ' {
			b, _ = reader.ReadByte()
		}
		for b != ' ' {
			buffer.WriteByte(b)
			b, _ = reader.ReadByte()
		}
		par.Name = buffer.String()
		buffer.Reset()
		for b == ' ' {
			b, _ = reader.ReadByte()
		}
		for b != ' ' && b != ',' && b != ')' {
			buffer.WriteByte(b)
			b, _ = reader.ReadByte()
		}
		par.Type = buffer.String()
		buffer.Reset()

		for b != ',' && b != ')' {
			b, _ = reader.ReadByte()
		}
		fc.params = append(fc.params, par)
	}

	b = ' '
	for b == ' ' {
		b, _ = reader.ReadByte()
	}
	for b != ' ' && b != '{' {
		buffer.WriteByte(b)
		b, _ = reader.ReadByte()
	}
	if len(buffer.String()) > 0 {
		fc.retval = append(fc.retval, buffer.String())
	}
	buffer.Reset()


	//Read content of function
	for scan.Scan() {
		
	}

	return fc
}

func ParseSimpleLine(line string) []*Instruction {
	return nil
}

func IsDelimChar(c byte) bool {
	return true
}

func IsInSet(b byte, set []byte) bool {
	for _,s := range set {
		if b == s {
			return true
		}
	}
	return false
}

func IsWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func Tokenize(src io.ByteReader) []string {
	toks := make([]string, 0, 256)
	buf := new(bytes.Buffer)
	//tokType := 0

	var err error
	var c byte
	for ; err == nil; c,err = src.ReadByte() {
		if IsWhitespace(c) {
			if buf.Len() > 0 {
				toks = append(toks, buf.String())
				buf.Reset()
			}
			if c == '\n' {
				toks = append(toks, ";")
			}
		} else if IsInSet(c, Punct) {
			if buf.Len() > 0 {
				toks = append(toks, buf.String())
				buf.Reset()
			}
			toks = append(toks, string([]byte{c}))
		} else {
			buf.WriteByte(c)
		}
	}
	if buf.Len() > 0 {
		toks = append(toks, buf.String())
	}

	return toks
}
*/
