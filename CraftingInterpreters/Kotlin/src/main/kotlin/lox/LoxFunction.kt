package lox

import expressions.Stmt

class LoxFunction(private val declaration: Stmt.Function) : LoxCallable {
    override fun arity(): Int = declaration.args.size

    override fun call(interpreter: Interpreter, args: List<Any?>): Any? {
        val env = Environment(interpreter.globals)
        for (i in 0 until declaration.args.size) {
            env.define(declaration.args[i].lexeme, args[i])
        }

        return interpreter.executeBlock(declaration.body, env)
    }

    override fun toString(): String = "<fn ${declaration.name.lexeme}>"
}