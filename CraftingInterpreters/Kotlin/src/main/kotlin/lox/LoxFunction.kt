package lox

import expressions.Stmt

class LoxFunction(private val declaration: Stmt.Function, private val closure: Environment, private val isInit: Boolean) : LoxCallable {
    override fun arity(): Int = declaration.args.size

    override fun call(interpreter: Interpreter, args: List<Any?>): Any? {
        val env = Environment(closure)
        for (i in 0 until declaration.args.size) {
            env.define(declaration.args[i].lexeme, args[i])
        }

        try {
            interpreter.executeBlock(declaration.body, env)
        } catch (returnValue: Return) {
            if (isInit) return closure.getAt(0, "this")
            returnValue.value
        }

        if (isInit) return closure.getAt(0, "this")
        return null
    }

    override fun toString(): String = "<fn ${declaration.name.lexeme}>"

    fun bind(loxInstance: LoxInstance): LoxFunction {
       val env = Environment(closure)
       env.define("this", loxInstance)
       return LoxFunction(declaration, env, isInit)
    }
}