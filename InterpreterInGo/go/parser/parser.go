package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := Parser{
		l:              lexer,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}
	//Register Prefix Functions
	p.addPrefixFunc(token.IDENT, p.parseIdentifier)
	p.addPrefixFunc(token.INT, p.parseIntegerLiteral)
	p.addPrefixFunc(token.TRUE, p.parseBooleanLiteral)
	p.addPrefixFunc(token.FALSE, p.parseBooleanLiteral)
	p.addPrefixFunc(token.BANG, p.parsePrefixOperation)
	p.addPrefixFunc(token.MINUS, p.parsePrefixOperation)
	p.addPrefixFunc(token.LPAREN, p.parseGroupedExpression)
	p.addPrefixFunc(token.IF, p.parseIfExpression)
	p.addPrefixFunc(token.FUNCTION, p.parseFunctionLiteral)
	p.addPrefixFunc(token.STRING, p.parseStringLiteral)
	p.addPrefixFunc(token.LBRACKET, p.parseArrayLiteral)
	p.addPrefixFunc(token.LBRACE, p.parseHashLiteral)

	//RegisterinfixFunctions
	p.addInfixFunc(token.MINUS, p.parseInfixExpression)
	p.addInfixFunc(token.PLUS, p.parseInfixExpression)
	p.addInfixFunc(token.STAR, p.parseInfixExpression)
	p.addInfixFunc(token.SLASH, p.parseInfixExpression)
	p.addInfixFunc(token.EQ, p.parseInfixExpression)
	p.addInfixFunc(token.NOT_EQ, p.parseInfixExpression)
	p.addInfixFunc(token.LT, p.parseInfixExpression)
	p.addInfixFunc(token.GT, p.parseInfixExpression)
	p.addInfixFunc(token.LPAREN, p.parseCall)
	p.addInfixFunc(token.LBRACKET, p.parseIndexExpression)

	//Read 2 token to initialise curent and peek
	p.nextToken()
	p.nextToken()

	return &p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addPrefixFunc(tokenType token.TokenType, fn prefixParseFn) {

	_, ok := p.prefixParseFns[tokenType]
	if ok {
		panic(fmt.Sprintf("already found prefix function for token type: %q", tokenType))
	}

	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) addInfixFunc(tokenType token.TokenType, fn infixParseFn) {

	_, ok := p.infixParseFns[tokenType]
	if ok {
		panic(fmt.Sprintf("already found prefix function for token type: %q", tokenType))
	}

	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}
