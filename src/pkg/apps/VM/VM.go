package VM

import (
	"io"
)

type VM struct {
	code []*Instruction
	ip   int

	memory [30000]int
	dp     int

	input  io.Reader
	output io.Writer

	readBuf []byte
}

func SpawnVM(instructions []*Instruction, in io.Reader, out io.Writer) *Machine {
	return &Machine{
		code:    instructions,
		input:   in,
		output:  out,
		readBuf: make([]byte, 1),
	}
}

func (m *VM) Execute() {
	for m.ip < len(m.code) {
		ins := m.code[m.ip]

		switch ins.Type {
		case Plus:
			m.memory[m.dp] += ins.Argument
		case Minus:
			m.memory[m.dp] -= ins.Argument
		case Right:
			m.dp += ins.Argument
		case Left:
			m.dp -= ins.Argument
		case PutChar:
			for i := 0; i < ins.Argument; i++ {
				m.putChar()
			}
		case ReadChar:
			for i := 0; i < ins.Argument; i++ {
				m.readChar()
			}
		case JumpIfZero:
			if m.memory[m.dp] == 0 {
				m.ip = ins.Argument
				continue
			}
		case JumpIfNotZero:
			if m.memory[m.dp] != 0 {
				m.ip = ins.Argument
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
