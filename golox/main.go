package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/scanner"

)

type golox struct {
	//We’ll use this to ensure we don’t try to execute code that has a known error.
	hadError bool
}

type TokenType int

const (

	// Single Character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	t_Type TokenType
	/*
		A lexeme is the lowest-level atomic unit of source code
		(e.g., a specific variable name sum, a keyword if, or an operator +)
		that is recognized by the lexer and passed to the parser.
		It is a raw sequence of characters matching a token's pattern.
		Lexemes are transformed into tokens, which are categorized structures used to build the syntax tree.
	*/
	lexeme  string
	literal struct{}
	line    int
}

func (t Token) New(passed_t_Type TokenType, passed_lexeme string, passed_literal struct{}, passed_line int) Token {
	return Token{
		t_Type:  passed_t_Type,
		lexeme:  passed_lexeme,
		literal: passed_literal,
		line:    passed_line,
	}
}

/*
Prints with the following formatting

	%d = integer
	%s = string
	%v = any value in a default format
*/

func (t Token) toString() string {
	return fmt.Sprintf("%d %s %v", t.t_Type, t.lexeme, t.literal)
}

// Our interpreter supports two ways of running code.
//
// If you start golox from the command line and give it a path to a file,
//
//	it reads the file and executes it.
func (g *golox) runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	run(string(bytes))
	if g.hadError {
		os.Exit(65)
	}
	g.hadError = false

}

// You can also run golox interactively.
// An interactive prompt is also called a “REPL”
//
//	Read, Evaluate, Print, Loop
//
// Fire up golox without any arguments
// and it drops you into a prompt where you can enter and execute code one line at a time.
func (g *golox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line)
	}
}

// Right now, it prints out the tokens our forthcoming scanner will emit so that we can see if we’re making progress.
func run(source string) {
	var s scanner.Scanner

	s.Init(strings.NewReader(source))
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s : %s\n", s.Position, s.TokenText())
	}
}

// tells the user some syntax error occurred on a given line.
func (g *golox) goloxError(line int, message string) {
	report(line, "", message)
	g.hadError = true
}

func report(line int, where, message string) {
	fmt.Println("[line ]" + strconv.Itoa(line) + "] Error" + where + ": " + message)
}

type Scanner struct {
	// We store the raw source code as a simple string
	source string

	// We then have an array ready to fill with tokens
	tokens []Token
}

func (s *Scanner) New(passed_source string) Scanner {
	return Scanner{
		source: passed_source,
	}
}

func (s *Scanner) scanTokens() {

}

// TODO: Make this better by doing something like
/*
Error: Unexpected "," in argument list.

    15 | function(first, second,);
                               ^-- Here.
*/

func main() {
	goLox := golox{hadError: false}
	argsSet := os.Args[1:]

	if len(argsSet) > 1 {
		println("Usage: golox [script]")
		os.Exit(64)
	} else if len(argsSet) == 1 {
		goLox.runFile(argsSet[0])
	} else {
		goLox.runPrompt()
	}
}
