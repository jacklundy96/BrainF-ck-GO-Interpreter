package main


/*
Language Rules:
> => Add 1 to the data pointer
< => Minus 1 on the data pointer
+ => Add 1 to the cell which the data pointer is pointing at
- => Minus 1 on the cell which the data pointer is pointing at
. => Take the integer stored in the current cell and convert it to ASCII then put it on the output stream
, => Take a character from the input stream convert it to an integer and  write it to the current cell
[ => Always used with the closing square bracket, if the current cell contains 0 then move the instruction pointer to the index after the matching closing bracket
] => Always used with the opening square bracket, if the cell doesn't contrain 0 then set the instruction pointer to the postion after the matching bracket
*/

//Useage: go build -o machine && ./machine ./[code].b

import (
	"fmt"
	"io/ioutil"
	"os"
	"io"
)



func main() {
	fileName := os.Args[1]
	code, err  := ioutil.ReadFile(fileName)

	if err != nil {
			fmt.Fprint(os.Stderr, "error: %s\n", err)
			os.Exit(-1)
	}

	m := SpawnVM(string(code), os.Stdin, os.Stdout)
	m.Execute()
}


//Type for the BrainFuck VM
type VM struct {

	//Code in and Instruction pointer
	code string 
	ip int
	
	//Memory allocation, data pointer and buffer slice 
	memory [30000]int
	buf []byte
	dp int 
	
	//Input and output streams 
	input io.Reader
	output io.Writer
}

//Creates a new brainfuck VM 
func SpawnVM (code string, in io.Reader, out io.Writer) *VM {
	return &VM {
		code: code, 
		input: in, 
		output: out, 
		buf: make([]byte, 1),
	}
}

//Executes code for the VM 
func (m *VM) Execute() {
	for m.ip < len(m.code) {
		ins := m.code[m.ip]

		switch ins {
		case '+': 
				ap := m.ip
				counter := 0
				for m.memory[ap] == '+' {
					counter ++ 
					ap ++
				}
				m.dp = ap
				m.memory[m.dp] += counter
		case '-': 
				ap := m.ip
				counter := 0
				for m.memory[ap] == '-' {
					counter ++ 
					ap ++
				}
				m.dp = ap
				m.memory[m.dp]  -= counter
		case '>': 
				m.dp++
		case '<':
				m.dp--	
		case '.':
				m.putChar()
		case ',':
				m.readChar()
		case '[':
				if m.memory[m.dp] == 0 {
					depth := 1 
					for depth != 0 {
						m.ip ++
						switch m.code[m.ip] {
						case '[':
							depth++
						case ']':
							depth--
						}
					} 
				}
		case ']':  
			if m.memory[m.dp] != 0 {
				depth := 1 
				for depth != 0 {
					m.ip --
					switch m.code[m.ip] {
					case ']':
						depth++
					case '[':
						depth--
					}
				} 
			}
		}

		m.ip++
	}
}


//============
//IO Methods
//============

func (m *VM) readChar() {
	n, err := m.input.Read(m.buf)
	if err != nil {
		panic(err)
	}
	if n != 1 {
		panic("Read the wrong number of bytes")
	}

	m.memory[m.dp] = int(m.buf[0])
}

func (m *VM) putChar() {
	m.buf[0] = byte(m.memory[m.dp])

	n, err := m.output.Write(m.buf)
	if err != nil {	
		panic(err)
	}
	if n != 1 {
		panic("Wrote the wrong number of bytes")
	}
}
