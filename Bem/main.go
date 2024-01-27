package main

import (
	readfile "Bem/Bem/readFile"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TokenType int

const (
	EOF TokenType = iota
	Identifier
	Number
	Plus
	Minus
	Multiply
	Divide
	Assign
	Print
	Scan
	String
	MathPrint
	PrintString
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
}

type Lexer struct {
	input   string
	scanner *bufio.Scanner
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.scanner = bufio.NewScanner(strings.NewReader(input))
	l.scanner.Split(bufio.ScanWords)
	return l
}
func (l *Lexer) NextToken() Token {
	if !l.scanner.Scan() {
		return Token{Type: EOF, Lexeme: "", Literal: nil}
	}

	token := l.scanner.Text()

	switch token {
	case "+":
		return Token{Type: Plus, Lexeme: token, Literal: nil}
	case "-":
		return Token{Type: Minus, Lexeme: token, Literal: nil}
	case "*":
		return Token{Type: Multiply, Lexeme: token, Literal: nil}
	case "/":
		return Token{Type: Divide, Lexeme: token, Literal: nil}
	case "=":
		return Token{Type: Assign, Lexeme: token, Literal: nil}
	case "Yazdir":
		return Token{Type: Print, Lexeme: token, Literal: nil}
	case "okut":
		return Token{Type: Scan, Lexeme: token, Literal: nil}
	case "Yazdirisle":
		return Token{Type: MathPrint, Lexeme: token, Literal: nil}
	default:
		if regMatch, _ := regexp.MatchString("'+[a-zA-z0-9_-]+'", token); regMatch {
			return Token{Type: String, Lexeme: token, Literal: nil}
		}
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			return Token{Type: Number, Lexeme: token, Literal: num}
		}
		return Token{Type: Identifier, Lexeme: token, Literal: token}
	}
}

type Parser struct {
	lexer     *Lexer
	currToken Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.consumeToken()
	return p
}

func (p *Parser) consumeToken() {
	p.currToken = p.lexer.NextToken()
}
func (p *Parser) parseExpression() float64 {
	result := p.parseTerm()

	for p.currToken.Type == Plus || p.currToken.Type == Minus {
		op := p.currToken.Type
		p.consumeToken()
		term := p.parseTerm()

		switch op {
		case Plus:
			result += term
		case Minus:
			result -= term

		}
	}

	return result
}

func (p *Parser) parseTerm() float64 {
	result := p.parseFactor()

	for p.currToken.Type == Multiply || p.currToken.Type == Divide {
		op := p.currToken.Type
		p.consumeToken()
		factor := p.parseFactor()

		switch op {
		case Multiply:
			result *= factor
		case Divide:
			result /= factor
		}
	}

	return result
}
func (p *Parser) parseFactor() float64 {
	switch p.currToken.Type {
	case Number:
		value := p.currToken.Literal.(float64)
		p.consumeToken()
		return value
	case Identifier:
		identifier := p.currToken.Literal.(string)
		p.consumeToken()
		return variables[identifier]
	case Minus:
		p.consumeToken()
		return -p.parseFactor()
	case Plus:
		p.consumeToken()
		return p.parseFactor()
	default:
		fmt.Println("Unexpected token:", p.currToken.Lexeme)
		os.Exit(1)
		return 0
	}
}

var variables map[string]float64

func main() {
	fmt.Println("bemlang")

	var fileName string
	fmt.Scan(&fileName)

	if !readfile.MatchFileName(fileName) {
		fmt.Println("Geçersiz dosya adı veya uzantısı")
		return
	}

	file, err := readfile.FileRead(fileName)
	if err != nil {
		fmt.Println("Dosya okuma hatası:", err)
		return
	}

	input := string(file)
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	variables = make(map[string]float64)

	for {
		if parser.currToken.Type == EOF {
			break
		}

		if parser.currToken.Type == Print {
			parser.consumeToken()
			identifier := parser.currToken.Literal.(string)
			fmt.Println(identifier)
			parser.consumeToken()
		} else if parser.currToken.Type == Scan {
			parser.consumeToken()
			identifier := parser.currToken.Literal.(string)
			fmt.Printf("Enter a value for %v: ", identifier)
			var inputValue float64
			fmt.Scan(&inputValue)
			variables[identifier] = inputValue
			parser.consumeToken()
		} else if parser.currToken.Type == MathPrint {
			parser.consumeToken()
			identifier := parser.currToken.Literal.(string)
			value := parser.parseExpression()
			variables[identifier] = value
			fmt.Println(value)
			parser.consumeToken()
		} else {
			identifier := parser.currToken.Literal.(string)
			parser.consumeToken()
			parser.consumeToken()
			value := parser.parseExpression()
			variables[identifier] = value
		}
	}

}
