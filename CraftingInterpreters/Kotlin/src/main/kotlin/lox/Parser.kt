package lox

import expressions.Expr
import expressions.Stmt
import java.lang.RuntimeException

class Parser(private val tokens: List<Token>) {

    fun parse() : List<Stmt?>? {
        val statements = mutableListOf<Stmt?>()
        while (!isAtEnd()) {
            statements.add(declaration())
        }
        return statements
    }

    private var current: Int = 0

    private fun expression(): Expr {
        return assignment()
    }

    private fun declaration(): Stmt? {
        return try {
            if (match(TokenType.VAR)) varDeclaration()
            else statement()
        } catch (e: ParseError) {
            synchronize()
            null
        }
    }

    private fun varDeclaration() : Stmt? {
        val name = consume(TokenType.IDENTIFIER, "Expect variable name")
        val initializer = if (match(TokenType.EQUAL)) expression() else null
        consume(TokenType.SEMICOLON, "Expect ';' after variable declaration")
        return Stmt.Var(name, initializer)
    }

    private fun statement(): Stmt {
        return if (match(TokenType.PRINT)) printStatement()
        else if (match(TokenType.LEFT_BRACE)) Stmt.Block(blockStatement())
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

    private fun blockStatement() : List<Stmt?> {
        val statements: MutableList<Stmt?> = mutableListOf()

        while (!check(TokenType.RIGHT_BRACE) && !isAtEnd()) {
            statements.add(declaration())
        }

        consume(TokenType.RIGHT_BRACE, "Expect '}' after block")
        return statements.toList()
    }

    //Expression parsers
    private fun assignment() : Expr {
        val expr = equality()

        if(match(TokenType.EQUAL)) {
            val equals = previous()
            val value = assignment()
            if (expr is Expr.Variable) {
                val name = expr.name
                return Expr.Assign(name, value)
            }

            error(equals, "Invalid assingment target")
        }

        return expr
    }
    private fun equality() = parseUntil(::comparison, TokenType.BANG_EQUAL, TokenType.EQUAL_EQUAL)
    private fun comparison() = parseUntil(::term, TokenType.GREATER, TokenType.GREATER_EQUAL, TokenType.LESS, TokenType.LESS_EQUAL)
    private fun term() = parseUntil(::factor, TokenType.MINUS, TokenType.PLUS)
    private fun factor() = parseUntil(::unary, TokenType.SLASH, TokenType.STAR)
    private fun unary() = parseUntil(::primary, TokenType.BANG, TokenType.MINUS)
    private fun primary() : Expr {
        if (match(TokenType.FALSE)) return Expr.Literal(false)
        if (match(TokenType.TRUE)) return Expr.Literal(true)
        if (match(TokenType.NUMBER, TokenType.STRING)) return Expr.Literal(previous().literal)
        if (match(TokenType.IDENTIFIER)) return Expr.Variable(previous())
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

    private fun check(vararg types: TokenType): Boolean {
        return if (isAtEnd()) {
            false
        } else {
            types.contains(peek().type)
        }
    }

    private fun consume(type : TokenType, message: String): Token{
        if (check(type)) return advance()
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