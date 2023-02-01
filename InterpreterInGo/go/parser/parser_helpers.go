package parser

import (
	"fmt"
	"monkey/token"
)

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if !p.peekTokenIs(t) {
		p.peekError(t)
		return false
	}
	p.nextToken()
	return true
}
func (p *Parser) peekError(t token.TokenType) {
	err := fmt.Sprintf("expected to find token %s but found %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) prefixNotFoundError(t token.Token) {
	err := fmt.Sprintf("Unable to find prefix for TokenType %q with Literal %q", t.Type, t.Literal)
	p.errors = append(p.errors, err)
}

func (p *Parser) infixNotFoundError(t token.Token) {
	err := fmt.Sprintf("Unable to find infix for TokenType %q with Literal %q", t.Type, t.Literal)
	p.errors = append(p.errors, err)
}

func (p *Parser) peekPrecedence() int {
	precedence, ok := precedences[p.peekToken.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}

func (p *Parser) curPrecedence() int {
	precedence, ok := precedences[p.curToken.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}
