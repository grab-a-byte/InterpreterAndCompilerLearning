package lox

class Scanner(private val input: String) {

    private val source = input.toCharArray()
    private val tokens: MutableList<Token> = mutableListOf()
    private var start: Int = 0
    private var current: Int = 0
    private var line: Int = 1

    private val keywords = hashMapOf(
        "and" to TokenType.AND,
        "class" to TokenType.CLASS,
        "else" to TokenType.ELSE,
        "false" to TokenType.FALSE,
        "for" to TokenType.FOR,
        "fun" to TokenType.FUN,
        "if" to TokenType.IF,
        "nil" to TokenType.NIL,
        "or" to TokenType.OR,
        "print" to TokenType.PRINT,
        "return" to TokenType.RETURN,
        "super" to TokenType.SUPER,
        "this" to TokenType.THIS,
        "true" to TokenType.TRUE,
        "var" to TokenType.VAR,
        "while" to TokenType.WHILE
    )

    fun scanTokens() : List<Token> {
        while(!isAtEnd()) {
           start = current;
           scanToken()
        }
        tokens.add(Token(TokenType.EOF, "", null, line))
        return tokens
    }

    private fun scanToken(){
        val c = advance()
        when(c) {
            '(' -> addToken(TokenType.LEFT_PAREN)
            ')' -> addToken(TokenType.RIGHT_PAREN)
            '{' -> addToken(TokenType.LEFT_BRACE)
            '}' -> addToken(TokenType.RIGHT_BRACE)
            ',' -> addToken(TokenType.COMMA)
            '.' -> addToken(TokenType.DOT)
            '-' -> addToken(TokenType.MINUS)
            '+' -> addToken(TokenType.PLUS)
            ';' -> addToken(TokenType.SEMICOLON)
            '*' -> addToken(TokenType.STAR)
            '!' -> addToken(if(match('=')) TokenType.BANG_EQUAL else TokenType.BANG)
            '=' -> addToken(if(match('=')) TokenType.EQUAL_EQUAL else TokenType.EQUAL)
            '<' -> addToken(if(match('=')) TokenType.LESS_EQUAL else TokenType.LESS)
            '>' -> addToken(if(match('=')) TokenType.GREATER_EQUAL else TokenType.GREATER)
            '\n' -> line++
            '/' -> {
                if (match('/')) {
                    while(peek() != '\n' && !isAtEnd()) advance()
                } else {
                    addToken(TokenType.SLASH)
                }
            }
            '"' -> string()
            ' ', '\r', '\t', -> {}
            else -> {
                if (isDigit(c)) number()
                else if (isAlpha(c)) identifier()
                else Lox.error(line, "Unexpected Character")
            }
        }
    }

    private fun isDigit(c: Char) : Boolean = ('0'..'9').contains(c)
    private fun isAlpha(c: Char): Boolean = (('a'..'z') + ('A'..'Z') + '_').contains(c)
    private fun isAlphaNumeric(c: Char) = isAlpha(c) || isDigit(c)

    private fun number() {
        while (isDigit(peek())) advance()

        //Look for fractional part
        if (peek() == '.' && isDigit(peekNext())) {
            //consume '.'
            advance()
            while (isDigit(peek())) advance()
        }
        val value = String(source.slice(start until current).toCharArray())
        addToken(TokenType.NUMBER, value.toDouble())
    }

    private fun identifier() {
        while (isAlphaNumeric(peek())) advance()
        val text = String(source.slice(start until current).toCharArray())
        val type = keywords[text] ?: TokenType.IDENTIFIER
        addToken(type)
    }

    private fun string() {
        while(peek() != '"' && !isAtEnd()) {
            if (peek() == '\n') line ++
            advance()
        }

        if(isAtEnd()) {
            Lox.error(line, "Unterminated String")
            return
        }

        advance() //Closing " char
        val value = source.slice(start+1 until current).joinToString()
        addToken(TokenType.STRING, value)
    }

    private fun isAtEnd() : Boolean = current >= source.size

    private fun addToken(type: TokenType) = addToken(type, null)

    private fun addToken(type: TokenType, literal: Any?){
        val text = String(source.slice(start until current).toCharArray())
        tokens.add(Token(type, text, literal, line))
    }

    private fun peek() : Char = if (isAtEnd()) '\n' else source[current]

    private fun peekNext() : Char {
        if (current + 1 > source.size) return '\n'
        return source[current + 1]
    }

    private fun match(expected : Char) : Boolean {
        if (isAtEnd()) return false
        if(source[current] != expected) return false
        current += 1
        return true
    }

    private fun advance() : Char {
        return source[current++]
    }
}