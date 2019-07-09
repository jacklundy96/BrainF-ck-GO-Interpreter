package main

import (
  "fmt"
  "bufio"
  "strings"
  "os"
) 

func main() {

  count := 0
	scanner := bufio.NewScanner(strings.NewReader(os.Stdin))
  
	for scanner.Scan() {
		tokenize(scanner.Text(), count)
    count ++
  }
}

func initTokens() ([8]string, [8]string) {
	TokenTags := [8]string{"INCREMENT","DECREMENT","SHIFT_LEFT","SHIFT_RIGHT","OUTPUT","INPUT","OPEN_LOOP","CLOSE_LOOP"}
	Tokens:= [8]string{"+","-","<",">",".","'","[","]"}

	return TokenTags, Tokens
}

func tokenize(line string, lineNumber int) {
	TokenTags, Tokens := initTokens()

  var err error
	for pos, char := range line {
		err := Tokens[char]
	}
	
	if err != nil {
		 fmt.Errorf("Error: %s", err)
	  }
}