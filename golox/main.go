package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

)

type Golox struct {
	//We’ll use this to ensure we don’t try to execute code that has a known error.
	hadError bool
}

func NewGolox() *Golox {
	return &Golox{
		hadError: false,
	}
}

// Right now, it prints out the tokens our forthcoming scanner will emit so that we can see if we’re making progress.
func (g *Golox) run(source string) {
	s := NewScanner(g, source)
	for _, tok := range s.scanTokens() {
		fmt.Println(tok.toString())
	}
}

// tells the user some syntax error occurred on a given line.
func (g *Golox) goloxError(line int, message string) {
	report(line, "", message)
	g.hadError = true
}

// Our interpreter supports two ways of running code.
//
// If you start golox from the command line and give it a path to a file,
//
//	it reads the file and executes it.
func (g *Golox) runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	g.run(string(bytes))
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
func (g *Golox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		g.run(line)
	}
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

func (t TokenType) String() string {
	switch t {
	case LEFT_PAREN:
		return "("
	case RIGHT_PAREN:
		return ")"
	case LEFT_BRACE:
		return "{"
	case RIGHT_BRACE:
		return "}"
	case COMMA:
		return ","
	case DOT:
		return "."
	case MINUS:
		return "-"
	case PLUS:
		return "+"
	case SEMICOLON:
		return ";"
	case SLASH:
		return "/"
	case STAR:
		return "*"

	// One or two character tokens.
	case BANG:
		return "!"
	case BANG_EQUAL:
		return "!="
	case EQUAL:
		return "="
	case EQUAL_EQUAL:
		return "=="
	case GREATER:
		return ">"
	case GREATER_EQUAL:
		return ">="
	case LESS:
		return "<"
	case LESS_EQUAL:
		return "<="

	// Literals
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"

	// Keywords
	case AND:
		return "and"
	case CLASS:
		return "class"
	case ELSE:
		return "else"
	case FALSE:
		return "false"
	case FUN:
		return "fun"
	case FOR:
		return "for"
	case IF:
		return "if"
	case NIL:
		return "nil"
	case OR:
		return "or"
	case PRINT:
		return "print"
	case RETURN:
		return "return"
	case SUPER:
		return "super"
	case THIS:
		return "this"
	case TRUE:
		return "true"
	case VAR:
		return "var"
	case WHILE:
		return "while"

	case EOF:
		return "EOF"
	}

	return "UNKNOWN"
}

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
	literal any
	line    int
}

func (t Token) New(passed_t_Type TokenType, passed_lexeme string, passed_literal string, passed_line int) Token {
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
	return fmt.Sprintf("tokenType:%s, tokenValue:%v", t.t_Type.String(), t.literal)
}

func report(line int, where, message string) {
	fmt.Println("[line ]" + strconv.Itoa(line) + "] Error" + where + ": " + message)
}

type Scanner struct {
	golox *Golox

	// We store the raw source code as a simple array of runes
	source []rune

	// We then have an array ready to fill with tokens
	tokens []Token

	// We define the set of reserved words in a map.
	keywords map[string]TokenType

	//The start and current fields are offsets, that index into the string.

	// The start field points to the first character in the lexeme being scanned
	start int

	// Current points at the charecter currently being considered
	current int

	// The line field tracks what source line CURRENT is on so we can produce tokens that know their location
	line int
}

func NewScanner(g *Golox, passed_source string) Scanner {

	keywords := make(map[string]TokenType)
	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["fun"] = FUN
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE

	return Scanner{
		golox:    g,
		source:   []rune(passed_source),
		tokens:   []Token{},
		keywords: keywords,
		start:    0,
		current:  0,
		line:     1,
	}
}

// little helper function that tells us if we’ve consumed all the characters
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++

	return true
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addAbstractToken(LEFT_PAREN)
	case ')':
		s.addAbstractToken(RIGHT_PAREN)
	case '{':
		s.addAbstractToken(LEFT_BRACE)
	case '}':
		s.addAbstractToken(RIGHT_BRACE)
	case ',':
		s.addAbstractToken(COMMA)
	case '.':
		s.addAbstractToken(DOT)
	case '-':
		s.addAbstractToken(MINUS)
	case '+':
		s.addAbstractToken(PLUS)
	case ';':
		s.addAbstractToken(SEMICOLON)
	case '*':
		s.addAbstractToken(STAR)
	case '!':
		token := BANG

		if s.match('=') {
			token = BANG_EQUAL
		}

		s.addAbstractToken(token)
	case '=':
		token := EQUAL

		if s.match('=') {
			token = EQUAL_EQUAL
		}

		s.addAbstractToken(token)
	case '<':
		token := LESS

		if s.match('=') {
			token = LESS_EQUAL
		}

		s.addAbstractToken(token)
	case '>':
		token := GREATER

		if s.match('=') {
			token = GREATER_EQUAL
		}

		s.addAbstractToken(token)
	case '/':
		/*
			This is similar to the other two-character operators,
			except that when we find a second /, we don’t end the token yet.
			Instead, we keep consuming characters until we reach the end of the line.
		*/

		if s.match('/') {
			// a comment goes until the end of the line
			for {
				if s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			}

		} else {
			s.addAbstractToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.goloxString()
	case 'o':
		if s.match(('r')) {
			s.addAbstractToken(OR)
		}
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {

			s.golox.goloxError(s.line, fmt.Sprintf("Unexpected character: %q", c))
		}
	}
}

func (s *Scanner) identifier() {
	for {
		if s.isAlphaNumeric(s.peek()) {
			s.advance()
		} else {
			break
		}
	}
	text := s.source[s.start:s.current]
	t_Type, OK := s.keywords[string(text)]
	if OK {
		t_Type = IDENTIFIER
	} else {
		// Regular user-defined indentifier
		t_Type = IDENTIFIER
	}
	s.addAbstractToken(t_Type)

}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) number() {
	for {
		if s.isDigit(s.peek()) {
			s.advance()
		} else {
			break
		}
	}

	// Look for a fractional part
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// consume the "."
		s.advance()

		for {
			if s.isDigit(s.peek()) {
				s.advance()
			}
		}
	}

	value, _ := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	s.addLiteralToken(NUMBER, value)
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) goloxString() {
	for {
		if s.peek() != '"' && !s.isAtEnd() {
			if s.peek() == '\n' {
				s.line++
			}
			s.advance()
		} else {
			break
		}
	}
	if s.isAtEnd() {
		s.golox.goloxError(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addLiteralToken(STRING, string(value))
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return s.source[s.current]
}

func (s *Scanner) advance() rune {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) addAbstractToken(t_Type TokenType) {
	s.addLiteralToken(t_Type, "")
}

func (s *Scanner) addLiteralToken(passed_t_Type TokenType, passed_literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, Token{passed_t_Type, text, passed_literal, s.line})
}

/*
The scanner works its way through the source code,
adding tokens until it runs out of characters.

Then it appends one final "end of file" token.
*/
func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", "", s.line})
	return s.tokens
}

// TODO: Make this better by doing something like
/*
Error: Unexpected "," in argument list.

    15 | function(first, second,);
                               ^-- Here.
*/

func main() {
	goLox := NewGolox()
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
