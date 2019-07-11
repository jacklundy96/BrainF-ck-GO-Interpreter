package main 

import (
	"fmt"
	"io/ioutil"
	"os"
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