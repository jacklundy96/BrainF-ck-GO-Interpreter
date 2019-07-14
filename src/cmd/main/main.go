package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"io"
	"src/pkg/apps/Compiler"
	"src/pkg/apps/VM"
)

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
