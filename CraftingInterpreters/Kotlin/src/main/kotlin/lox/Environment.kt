package lox

class Environment(private val enclosing: Environment? = null) {

    private val values = mutableMapOf<String, Any?>()

    fun define(name: String, value: Any?) {
        values[name] = value
    }

    fun get(name: Token): Any? {
        val x = values[name.lexeme]
        if (enclosing != null && x == null) {
            return enclosing.get(name)
        } else if (x != null) {
            return x
        }
        throw RuntimeError(name, "undefined variable")
    }

    fun getAt(distance: Int, name: String) : Any? {
        return ancestor(distance).values[name];
    }

    fun assign(name: Token, value: Any?) {
        if (values.containsKey(name.lexeme)) values[name.lexeme] = value
        else if (enclosing != null) enclosing.assign(name, value)
        else throw RuntimeError(name, "variable ${name.lexeme} does not exist in the current context")
    }

    fun assignAt(distance: Int, name: Token, value: Any?) {
        ancestor(distance).values[name.lexeme] = value
    }

    private fun ancestor(distance : Int) : Environment {
        var environment: Environment? = this
        for ( i in 0 until distance) {
            if (environment?.enclosing != null) environment = environment.enclosing
        }

        return environment ?: this
    }
}