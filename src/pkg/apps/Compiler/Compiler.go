package Compiler 

type InstructionType byte 

type Compiler struct {
	code string 
	codeLength int 
	position int 
	instructions []*Instruction
}

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