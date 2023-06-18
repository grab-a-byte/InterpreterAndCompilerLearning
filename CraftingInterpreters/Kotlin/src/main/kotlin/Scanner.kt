class Scanner(private val input: String) {

    private val source = input.toCharArray()
    private val tokens: MutableList<Token> = mutableListOf()
    private var start: Int = 0
    private var current: Int = 0
    private var line: Int = 1

    fun scanTokens() : List<Token> {
        while(!isAtEnd()) {
           start = current;
           scanToken()
        }
        tokens.add(Token(TokenType.EOF, "", null, line))
        return tokens
    }

    private fun isAtEnd() : Boolean {
        return current >= source.size
    }

    private fun addToken(type: TokenType) {
        addToken(type, null)
    }

    private fun addToken(type: TokenType, literal: Any?){
        val text = source.slice(start..current).joinToString()
        tokens.add(Token(type, text, literal, line))
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
                if (isDigit(c)) number() else Lox.error(line, "Unexpected Character")
            }
        }
    }

    private fun isDigit(c: Char) : Boolean {
        return ('0'..'9').contains(c)
    }

    private fun number() {
        while (isDigit(peek())) advance()

        //Look for fractional part
        if (peek() == '.' && isDigit(peekNext())) {
            //consume '.'
            advance()
            while (isDigit(peek())) advance()
        }
        val value = source.slice(start..current).joinToString()
        addToken(TokenType.NUMBER, value.toDouble())
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

    private fun peek() : Char {
        return if (isAtEnd()) '\n' else source[current]
    }

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
        current += 1
        return source[current]
    }
}