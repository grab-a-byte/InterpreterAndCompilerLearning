package lox

class LoxInstance(val klass: LoxClass) {
    private val fields = mutableMapOf<String, Any?>()

    override fun toString(): String = "${klass.name} instance"

    fun get(name: Token): Any? {
        if (fields.containsKey(name.lexeme)) {
            fields[name.lexeme]
        }
        val method = klass.findMethod(name.lexeme)
        if (method != null) return method
        throw RuntimeError(name, "Undefined Property ${name.lexeme}.")
    }

    fun set(name: Token, value: Any?) {
       fields[name.lexeme] = value
    }
}