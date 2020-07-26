package parser

import "fmt"

type Expr struct {
	Op    string
	Left  *Expr
	Right *Expr
	Val   string
}

type Parser struct {
	lexer Lexer
}

func NewParser(lexer Lexer) Parser {
	return Parser{
		lexer: lexer,
	}
}

func (p *Parser) Parse() {
	tokens := p.lexer.Tokens()

	for _, t := range tokens {
		p := Expr{
			Op:    t.Type,
			Left:  nil,
			Right: nil,
			Val:   t.Value,
		}

		fmt.Println(t)
		fmt.Println(p)
	}
}
