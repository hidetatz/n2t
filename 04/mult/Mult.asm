// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

// Adds R0 to R2 for R1 times

	// initialize i (loop index) as 0
	@i
	M=0
	// initialize R2 (a register to store the result) as 0
	@2
	M=0
(LOOP)
	// store i to D
	@i
	D=M
	// now, D is i
	// calculate current loop time
	// i = i - R1 -> negative value if more loop is needed
	@1
	D=D-M

	// Jump to END if no more loop is needed
	@END
	D;JGE

	// load R0 to register
	@0
	D=M

	// add loaded R0 to R2
	@2
	M=D+M

	// increment i
	@i
	M=M+1

	// loop again
	@LOOP
	0;JMP
(END)
	@END
	0;JMP // finish program
