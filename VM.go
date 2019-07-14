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

	compiler := SpawnCompiler(string(code))
	instructions := compiler.Compile()

	m := SpawnVM(instructions, os.Stdin, os.Stdout)
	m.Execute()
}

type InstructionType byte 

//Type for the BrainFuck VM
type VM struct {

	//Code in and Instruction pointer
	code []*Instruction 
	ip int
	
	//Memory allocation, data pointer and buffer slice 
	memory [30000]int
	buf []byte
	dp int 
	
	//Input and output streams 
	input io.Reader
	output io.Writer
}

type Compiler struct {
	code string 
	codeLength int 
	position int 
	instructions []*Instruction
}

const (
	Plus InstructionType = '+'
	Minus InstructionType = '-'
	Right InstructionType = '>'
	Left InstructionType = '<'
	PutChar InstructionType = '.'
	ReadChar InstructionType = ','
	JumpIfZero InstructionType = '['
	JumpIfNotZero InstructionType = ']'
)

type Instruction struct {
	Type InstructionType
	Arguement int
}

//Creates a new brainfuck VM 
func SpawnVM (instructions []*Instruction, in io.Reader, out io.Writer) *VM {
	return &VM {
		code: instructions, 
		input: in, 
		output: out, 
		buf: make([]byte, 1),
	}
}

func SpawnCompiler (code string) *Compiler {
	return &Compiler {
		code: code,
		codeLength: len(code),
		instructions: []*Instruction{},
	}
}

func (c *Compiler) Compile() []*Instruction {
	loopStack := []int{}
	
	for c.position < c.codeLength {
		current := c.code[c.position]

		switch current {
		case '+':
			c.CompileFoldableInstruction('+',Plus)
		case '-':
			c.CompileFoldableInstruction('-',Minus)
		case '<':
			c.CompileFoldableInstruction('<',Left)
		case '>':
			c.CompileFoldableInstruction('>',Right)
		case '.':
			c.CompileFoldableInstruction('.',PutChar)
		case ',':
			c.CompileFoldableInstruction(',',ReadChar)
		case '[':
			insPos := c.EmitWithArg(JumpIfZero, 0)
			loopStack = append(loopStack, insPos)
		case ']':
			// Pop position of last JumpIfZero ("[") instruction off stack
			openInstruction := loopStack[len(loopStack)-1]
			loopStack = loopStack[:len(loopStack)-1]
	
			// Emit the new JumpIfNotZero ("]") instruction,
			// with correct position as argument
			closeInstructionPos := c.EmitWithArg(JumpIfNotZero, openInstruction)
	
			// Patch the old JumpIfZero ("[") instruction with new position
			c.instructions[openInstruction].Arguement = closeInstructionPos
		}
		c.position++
	}

	return c.instructions
}

func (c *Compiler) CompileFoldableInstruction(char byte, insType InstructionType) {
	count := 1

	for c.position < c.codeLength-1 && c.code[c.position+1] == char {
		count++
		c.position++
	}

	c.EmitWithArg(insType, count)
}

func (c *Compiler) EmitWithArg(insType InstructionType, arg int) int {
	ins := &Instruction{Type: insType, Arguement: arg}
	c.instructions = append(c.instructions, ins)
	return len(c.instructions) - 1
}

//Executes code for the VM 
func (m *VM) Execute() {
	for m.ip < len(m.code) {
		ins := m.code[m.ip]

		switch ins.Type {
		case Plus: 
				m.memory[m.dp] += ins.Arguement
		case Minus: 
				m.memory[m.dp] -= ins.Arguement
		case Right: 
				m.dp += ins.Arguement
		case Left:
				m.dp -= ins.Arguement	
		case PutChar:
				for i := 0; i < ins.Arguement; i++ {
					m.putChar()
				}
		case ReadChar:
			for i := 0; i < ins.Arguement; i++ {
				m.readChar()
			}
		case JumpIfZero:
				if m.memory[m.dp] == 0 {
					m.ip = ins.Arguement
					continue
				}
		case JumpIfNotZero:  
				if m.memory[m.dp] != 0 {
					m.ip = ins.Arguement
					continue
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
