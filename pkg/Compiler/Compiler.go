package Compiler

import (
	"io"
)

//============
//Type definitions and constants
//============

type InstructionType byte 

type Instruction struct {
	Type InstructionType
	Arguement int
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

type Compiler struct {
	code string 
	codeLength int 
	position int 
	instructions []*Instruction
}

type VM struct {
	code []*Instruction
	ip   int

	memory [30000]int
	dp     int

	input  io.Reader
	output io.Writer

	readBuf []byte
}

//============
//Functions
//============

func SpawnCompiler (code string) *Compiler {
	return &Compiler {
		code: code,
		codeLength: len(code),
		instructions: []*Instruction{},
	}
}

func SpawnVM(instructions []*Instruction, in io.Reader, out io.Writer) *VM {
	return &VM{
		code:    instructions,
		input:   in,
		output:  out,
		readBuf: make([]byte, 1),
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
//Internal IO Methods
//============

func (m *VM) readChar() {
	n, err := m.input.Read(m.readBuf)
	if err != nil {
		panic(err)
	}
	if n != 1 {
		panic("Read the wrong number of bytes")
	}

	m.memory[m.dp] = int(m.readBuf[0])
}

func (m *VM) putChar() {
	m.readBuf[0] = byte(m.memory[m.dp])

	n, err := m.output.Write(m.readBuf)
	if err != nil {	
		panic(err)
	}
	if n != 1 {
		panic("Wrote the wrong number of bytes")
	}
}
