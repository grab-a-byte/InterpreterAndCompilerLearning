package lox

class LoxClass(val name: String, private val superclass: LoxClass?, private val methods: Map<String, LoxFunction>) :
    LoxCallable {
    override fun arity(): Int {
        val init = findMethod("init")
        return init?.arity() ?: 0
    }

    override fun call(interpreter: Interpreter, args: List<Any?>): Any? {
        val instance = LoxInstance(this)
        val initializer = findMethod("init")
        initializer?.bind(instance)?.call(interpreter, args)
        return instance
    }

    override fun toString() = name

    fun findMethod(name: String): LoxFunction? {
        if (methods.containsKey(name)) {
            return methods[name]
        }

        if (superclass != null) {
            return superclass.findMethod(name)
        }

        return null
    }
}