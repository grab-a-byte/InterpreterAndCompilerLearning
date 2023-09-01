package lox

class LoxClass(val name : String) : LoxCallable {
    override fun arity(): Int = 0

    override fun call(interpreter: Interpreter, args: List<Any?>): Any? {
        val instance = LoxInstance(this)
        return instance
    }

    override fun toString() = name
}