# BrainF*ck-GO-Compiler
An interpreter written in go for the BrainF*ck esoteric programming language, in this project i have two main goals: learn how to write a  simple interpreter an gain some expereince with Go.


# Rules of the Language
Language Rules:
> => Add 1 to the data pointer
< => Minus 1 on the data pointer
+ => Add 1 to the cell which the data pointer is pointing at
- => Minus 1 on the cell which the data pointer is pointing at
. => Take the integer stored in the current cell and convert it to ASCII then put it on the output stream
, => Take a character from the input stream convert it to an integer and  write it to the current cell
[ => Always used with the closing square bracket, if the current cell contains 0 then move the instruction pointer to the index after the matching closing bracket
] => Always used with the opening square bracket, if the cell doesn't contrain 0 then set the instruction pointer to the postion after the matching bracket