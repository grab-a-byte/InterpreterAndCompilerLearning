package lox

import expressions.Stmt

class LoxFunction(private val declaration: Stmt.Function, private val closure: Environment) : LoxCallable {
    override fun arity(): Int = declaration.args.size

    override fun call(interpreter: Interpreter, args: List<Any?>): Any? {
        val env = Environment(closure)
        for (i in 0 until declaration.args.size) {
            env.define(declaration.args[i].lexeme, args[i])
        }

        return try {
            interpreter.executeBlock(declaration.body, env)
            null
        } catch (returnValue: Return) {
            returnValue.value
        }
    }

    override fun toString(): String = "<fn ${declaration.name.lexeme}>"
}