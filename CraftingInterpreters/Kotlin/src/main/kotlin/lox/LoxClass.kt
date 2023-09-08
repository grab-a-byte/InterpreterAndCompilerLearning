package lox

class LoxClass(val name : String, private val methods: Map<String, LoxFunction>) : LoxCallable {
    override fun arity(): Int = 0

    override fun call(interpreter: Interpreter, args: List<Any?>): Any? {
        val instance = LoxInstance(this)
        return instance
    }

    override fun toString() = name

    fun findMethod(name: String): Any? {
       if (methods.containsKey(name)) {
           return methods[name]
       }

       return null
    }
}