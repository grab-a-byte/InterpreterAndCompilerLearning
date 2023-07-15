package lox

class Environment(private val enclosing: Environment? = null) {

    private val values = mutableMapOf<String, Any?>()

    fun define(name: String, value: Any?) {
        values[name] = value
    }

    fun get(name: Token) {
        val x = values[name.lexeme]
        if (enclosing != null && x == null) {
            return enclosing.get(name)
        }
        throw RuntimeError(name, "undefined variable")
    }

    fun assign(name: Token, value: Any?) {
        if (values.containsKey(name.lexeme)) values[name.lexeme] = value
        else if (enclosing != null) enclosing.assign(name, value)
        else throw RuntimeError(name, "variable ${name.lexeme} does not exist in the current context")
    }
}