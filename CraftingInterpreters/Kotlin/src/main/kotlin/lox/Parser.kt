package lox

import expressions.Expr
import expressions.Stmt

class Parser(private val tokens: List<Token>) {

    fun parse(): List<Stmt?> {
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

    private fun varDeclaration(): Stmt {
        val name = consume(TokenType.IDENTIFIER, "Expect variable name")
        val initializer = if (match(TokenType.EQUAL)) expression() else null
        consume(TokenType.SEMICOLON, "Expect ';' after variable declaration")
        return Stmt.Var(name, initializer)
    }

    private fun statement(): Stmt {
        return if (match(TokenType.IF)) ifStatement()
        else if (match(TokenType.FOR)) forStmt()
        else if (match(TokenType.PRINT)) printStatement()
        else if (match(TokenType.LEFT_BRACE)) Stmt.Block(blockStatement())
        else if (match(TokenType.WHILE)) whileStmt()
        else if (match(TokenType.FUN)) functionStatement("function")
        else if (match(TokenType.CLASS)) classDeclaration()
        else if (match(TokenType.RETURN)) returnStatement()
        else expressionStatement()
    }

    //Statement parsers
    private fun classDeclaration(): Stmt {
        val name = consume(TokenType.IDENTIFIER, "expected class name")
        val superclass = if (match(TokenType.LESS)) {
            consume(TokenType.IDENTIFIER, "expect superclass name")
            Expr.Variable(previous())
        } else {
            null
        }
        consume(TokenType.LEFT_BRACE, "expect left brace after class name")
        val methods: MutableList<Stmt.Function> = mutableListOf()
        while (!check(TokenType.RIGHT_BRACE) && !isAtEnd()) {
            methods.add(functionStatement("method"))
        }
        consume(TokenType.RIGHT_BRACE, "expect right brace at end of class")
        return Stmt.Class(name, superclass, methods.toList())
    }

    private fun returnStatement(): Stmt {
        val keyword = previous()
        val value = if (!check(TokenType.SEMICOLON)) expression() else null
        consume(TokenType.SEMICOLON, "expect ';' after return")
        return Stmt.Return(keyword, value)
    }

    private fun functionStatement(kind: String): Stmt.Function {
        val name = consume(TokenType.IDENTIFIER, "expected $kind name.")
        consume(TokenType.LEFT_PAREN, "expected '(' after $kind name.")

        val args = mutableListOf<Token>()
        if (!check(TokenType.RIGHT_PAREN)) {
            do {
                if (args.size >= 255) {
                    error(peek(), "Function call cannot have more than 255 arguments")
                }
                args.add(consume(TokenType.IDENTIFIER, "expected identifier name for argument"))
            } while (match(TokenType.COMMA))
        }

        consume(TokenType.RIGHT_PAREN, "Expected ')' after arguments")
        consume(TokenType.LEFT_BRACE, "Expected '{' after arguments")

        val body = blockStatement()

        return Stmt.Function(name, args, body)
    }

    private fun ifStatement(): Stmt {
        consume(TokenType.LEFT_PAREN, "Expected Left Paren after 'If'")
        val condition = expression()
        consume(TokenType.RIGHT_PAREN, "expected Right paren after if condition")
        val ifBranch = statement()
        val elseBranch = if (match(TokenType.ELSE)) statement() else null

        return Stmt.If(condition, ifBranch, elseBranch)
    }

    private fun printStatement(): Stmt {
        val expr = expression()
        consume(TokenType.SEMICOLON, "Expect ';' after value.")
        return Stmt.Print(expr)
    }

    private fun expressionStatement(): Stmt {
        val expr = expression()
        consume(TokenType.SEMICOLON, "Expect ';' after expression.")
        return Stmt.Expression(expr)
    }

    private fun blockStatement(): List<Stmt?> {
        val statements: MutableList<Stmt?> = mutableListOf()

        while (!check(TokenType.RIGHT_BRACE) && !isAtEnd()) {
            statements.add(declaration())
        }

        consume(TokenType.RIGHT_BRACE, "Expect '}' after block")
        return statements.toList()
    }

    private fun whileStmt(): Stmt {
        consume(TokenType.LEFT_PAREN, "Expected '(' after 'While'")
        val condition = expression()
        consume(TokenType.RIGHT_PAREN, "expected ')' after while condition")
        val body = statement()

        return Stmt.While(condition, body)
    }

    private fun forStmt(): Stmt {
        consume(TokenType.LEFT_PAREN, "expect '(' after 'for'")

        val initializer: Stmt? = if (match(TokenType.SEMICOLON)) null
        else if (match(TokenType.VAR)) varDeclaration()
        else expressionStatement()

        val condition = if (!check(TokenType.SEMICOLON)) expression() else null
        consume(TokenType.SEMICOLON, "expected ';' after loop condition")
        val increment = if (!check(TokenType.RIGHT_PAREN)) expression() else null

        consume(TokenType.RIGHT_PAREN, "expected ')' after 'if' conditions")

        var body = statement()
        if (increment != null) {
            body = Stmt.Block(listOf(body, Stmt.Expression(increment)))
        }

        if (condition != null) {
            body = Stmt.While(condition, body)
        }

        if (initializer != null) {
            body = Stmt.Block(listOf(initializer, body))
        }

        return body
    }

    //Expression parsers
    private fun call(): Expr {
        var expr = primary()
        while (true) {
            expr = if (match(TokenType.LEFT_PAREN)) {
                finishCall(expr)
            } else if (match(TokenType.DOT)) {
                val name = consume(TokenType.IDENTIFIER, "expected property name after .")
                Expr.Get(expr, name)
            } else {
                break
            }
        }

        return expr
    }

    private fun finishCall(inner: Expr): Expr {
        val args = mutableListOf<Expr>()

        if (!check(TokenType.RIGHT_PAREN)) {
            do {
                if (args.size >= 255) {
                    error(peek(), "Function call cannot have more than 255 arguments")
                }
                args.add(expression())
            } while (match(TokenType.COMMA))
        }
        val paren = consume(TokenType.RIGHT_PAREN, "Expected ')' after arguments")
        return Expr.Call(inner, paren, args)
    }

    private fun assignment(): Expr {
        val expr = or()

        if (match(TokenType.EQUAL)) {
            val equals = previous()
            val value = assignment()
            if (expr is Expr.Variable) {
                val name = expr.name
                return Expr.Assign(name, value)
            } else if (expr is Expr.Get) {
                return Expr.Set(expr.obj, expr.name, value)
            }

            error(equals, "Invalid assignment target")
        }

        return expr
    }

    private fun or(): Expr {
        val left = and()

        if (match(TokenType.OR)) {
            val operator = previous()
            val right = and()

            return Expr.Logical(left, operator, right)
        }

        return left
    }

    private fun and(): Expr {
        val left = equality()
        if (match(TokenType.AND)) {
            val operator = previous()
            val right = equality()

            return Expr.Logical(left, operator, right)
        }

        return left
    }

    private fun equality() = parseUntil(::comparison, TokenType.BANG_EQUAL, TokenType.EQUAL_EQUAL)
    private fun comparison() =
        parseUntil(::term, TokenType.GREATER, TokenType.GREATER_EQUAL, TokenType.LESS, TokenType.LESS_EQUAL)

    private fun term() = parseUntil(::factor, TokenType.MINUS, TokenType.PLUS)
    private fun factor() = parseUntil(::unary, TokenType.SLASH, TokenType.STAR)
    private fun unary() = parseUntil(::call, TokenType.BANG, TokenType.MINUS)
    private fun primary(): Expr {
        return if (match(TokenType.FALSE)) Expr.Literal(false)
        else if (match(TokenType.TRUE)) Expr.Literal(true)
        else if (match(TokenType.NIL)) Expr.Literal(null)
        else if (match(TokenType.NUMBER, TokenType.STRING)) Expr.Literal(previous().literal)
        else if (match(TokenType.THIS)) Expr.This(previous())
        else if (match(TokenType.IDENTIFIER)) Expr.Variable(previous())
        else if (match(TokenType.LEFT_PAREN)) {
            val expr = expression()
            consume(TokenType.RIGHT_PAREN, "Expect ')' after expressions.")
            Expr.Grouping(expr)
        } else if (match(TokenType.SUPER)) {
            val keyword = previous()
            consume(TokenType.DOT, "expected '.' after 'super'")
            val method = consume(TokenType.IDENTIFIER, "expect 'super.method' syntax, not method name found")
            Expr.Super(keyword, method)
        } else {
            throw error(peek(), "expected expression")
        }
    }

    private fun parseUntil(func: () -> Expr, vararg types: TokenType): Expr {
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

    private fun consume(type: TokenType, message: String): Token {
        if (check(type)) return advance()
        throw error(peek(), message)
    }

    private fun advance(): Token {
        if (!isAtEnd()) current++
        return previous()
    }

    private fun error(tok: Token, message: String): ParseError {
        Lox.error(tok, message)
        return ParseError()
    }

    private fun synchronize() {
        advance()
        val returnableTypes = listOf(
            TokenType.CLASS,
            TokenType.FOR,
            TokenType.FUN,
            TokenType.IF,
            TokenType.PRINT,
            TokenType.RETURN,
            TokenType.VAR,
            TokenType.WHILE
        )
        while (!isAtEnd()) {
            if (previous().type == TokenType.SEMICOLON) return
            if (returnableTypes.contains(peek().type)) return
            advance()
        }
    }

    private fun isAtEnd(): Boolean = peek().type == TokenType.EOF
    private fun previous() = tokens[current - 1]
    private fun peek() = tokens[current]

    private class ParseError : RuntimeException()
}