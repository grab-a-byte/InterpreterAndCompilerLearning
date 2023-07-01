package lox

import expressions.Expr
import expressions.Stmt
import java.lang.RuntimeException

class Parser(private val tokens: List<Token>) {

    fun parse() : List<Stmt?>? {
        val statements = mutableListOf<Stmt?>()
        while (!isAtEnd()) {
            statements.add(statement())
        }
        return statements
    }

    private var current: Int = 0

    private fun expression(): Expr {
        return equality()
    }

    private fun statement(): Stmt? {
        return if (match(TokenType.PRINT)) printStatement()
        else expressionStatement()
    }

    //Statement parsers
    private fun printStatement() : Stmt {
        val expr = expression()
        consume(TokenType.SEMICOLON, "Expect ';' after value.")
        return Stmt.Print(expr)
    }

    private fun expressionStatement(): Stmt {
        val expr = expression()
        consume(TokenType.SEMICOLON, "Expect ';' after expression.")
        return Stmt.Expression(expr)
    }

    //Expression parsers
    private fun equality() = parseUntil(::comparison, TokenType.BANG_EQUAL, TokenType.EQUAL_EQUAL)
    private fun comparison() = parseUntil(::term, TokenType.GREATER, TokenType.GREATER_EQUAL, TokenType.LESS, TokenType.LESS_EQUAL)
    private fun term() = parseUntil(::factor, TokenType.MINUS, TokenType.PLUS)
    private fun factor() = parseUntil(::unary, TokenType.SLASH, TokenType.STAR)
    private fun unary() = parseUntil(::primary, TokenType.BANG, TokenType.MINUS)
    private fun primary() : Expr {
        if (match(TokenType.FALSE)) return Expr.Literal(false)
        if (match(TokenType.TRUE)) return Expr.Literal(true)
        if (match(TokenType.NUMBER, TokenType.STRING)) return Expr.Literal(previous().literal)
        if(match(TokenType.LEFT_PAREN)) {
            val expr = expression()
            consume(TokenType.RIGHT_PAREN, "Expect ')' after expressions.")
            return Expr.Grouping(expr)
        }

        throw error(peek(), "expected expression");
    }


    private fun parseUntil(func: () -> Expr, vararg types: TokenType)  : Expr {
        var expr = func()
        while (match(*types)) {
            val operator = previous()
            val right = func()
            expr = Expr.Binary(expr, operator, right)
        }
        return expr
    }

    private fun match(vararg types: TokenType): Boolean {
        if (isAtEnd()) {
            return false
        } else {
            if (types.contains(peek().type)) {
                advance()
                return true
            }
            return false
        }
    }

    private fun consume(type : TokenType, message: String): Token{
        if (match(type)) return advance()
        throw error(peek(), message)
    }

    private fun advance(): Token {
        if (!isAtEnd()) current++
        return previous()
    }

    private fun error(tok: Token, message: String) : ParseError {
        Lox.error(tok, message)
        return ParseError()
    }

    private fun synchronize() {
        advance()
        val returnableTypes = listOf(TokenType.CLASS, TokenType.FOR, TokenType.FUN, TokenType.IF, TokenType.PRINT, TokenType.RETURN, TokenType.VAR, TokenType.WHILE)
        while(!isAtEnd()) {
            if (previous().type == TokenType.SEMICOLON) return
            if(returnableTypes.contains(peek().type)) return
            advance()
        }
    }

    private fun isAtEnd() : Boolean = peek().type == TokenType.EOF
    private fun previous() = tokens[current -1]
    private fun peek() = tokens[current]

    private class ParseError : RuntimeException() {}
}