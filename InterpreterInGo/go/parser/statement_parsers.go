package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	stmt.Name = name

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	exp := p.parseExpression(LOWEST)
	stmt.Value = exp

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	stmt.ReturnValue = exp

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.prefixNotFoundError(p.curToken)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseExpresisonStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpresisonStatement()
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("Unable to parse %q into integer", p.curToken.Literal)
		p.errors = append(p.errors, errMessage)
	}

	return &ast.IntegerLiteral{Token: p.curToken, Value: value}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanExpression{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixOperation() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken()
	rightExp := p.parseExpression(PREFIX)
	exp.Expression = rightExp

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token: p.curToken,
		Left:  left, Operator: p.curToken.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBlockExpression() ast.BlockStatement {
	exp := ast.BlockStatement{
		Token: p.curToken,
	}
	p.nextToken()

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			exp.Statements = append(exp.Statements, stmt)
		}
		p.nextToken()
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	cond := p.parseExpression(LOWEST)
	ifExpression.Condition = cond

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	consequence := p.parseBlockExpression()
	ifExpression.Consequence = &consequence

	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		alt := p.parseBlockExpression()
		ifExpression.Alternative = &alt
	}

	return ifExpression
}

func (p *Parser) parseFunctionParamters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, &ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, &ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	literal.Parameters = p.parseFunctionParamters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	body := p.parseBlockExpression()
	literal.Body = &body

	return literal
}

func (p *Parser) parseCallArguments() []ast.Expression {
	params := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()
	param := p.parseExpression(LOWEST)
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		param := p.parseExpression(LOWEST)
		params = append(params, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	call := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	call.Arguments = p.parseCallArguments()
	return call
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.Array{}

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return array
	}

	p.nextToken()
	first := p.parseExpression(LOWEST)
	array.Items = append(array.Items, first)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		exp := p.parseExpression(LOWEST)
		array.Items = append(array.Items, exp)
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	idxExp := ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	idx := p.parseExpression(LOWEST)

	idxExp.Index = idx
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return &idxExp
}
