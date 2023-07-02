package lox

class Environment {
    private val values = mutableMapOf<String, Any?>()

    fun define(name: String, value: Any?) {
        values[name] = value
    }

    fun get(name: Token) = values[name.lexeme] ?: throw RuntimeError(name, "undefined variable")

}