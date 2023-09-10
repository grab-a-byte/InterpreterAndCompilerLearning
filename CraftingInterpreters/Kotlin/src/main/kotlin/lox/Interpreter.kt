package lox

import expressions.Expr
import expressions.Stmt

class Interpreter : Expr.Visitor<Any?>, Stmt.Visitor<Any?> {

    private val globals = Environment()
    private var environment = globals
    private val locals: MutableMap<Expr, Int> = mutableMapOf()

    init {
        globals.define("clock", object : LoxCallable {
            override fun arity(): Int = 0
            override fun call(interpreter: Interpreter, args: List<Any?>): Any = System.currentTimeMillis()
            override fun toString(): String = "<native fun>"
        })
    }

    fun resolve(expr: Expr, depth: Int) {
        locals[expr] = depth
    }

    fun interpret(stmts: List<Stmt?>) {
        try {
            for (stmt in stmts) {
                if (stmt == null) continue
                execute(stmt)
            }
        } catch (e: RuntimeError) {
            Lox.runtimeError(e)
        }
    }

    private fun execute(stmt: Stmt) {
        stmt.accept(this)
    }

    override fun visitAssignExpr(expr: Expr.Assign): Any? {
        val value = evaluate(expr.value)
        val distance = locals[expr]
        if (distance != null) environment.assignAt(distance, expr.name, value)
        else environment.assign(expr.name, value)
        return value
    }

    override fun visitBinaryExpr(expr: Expr.Binary): Any? {
        val left = evaluate(expr.left)
        val right = evaluate(expr.right)
        return when (expr.operator.type) {
            TokenType.MINUS -> {
                checkNumberOperands(expr.operator, left, right)
                (left as Double) - (right as Double)
            }

            TokenType.STAR -> {
                checkNumberOperands(expr.operator, left, right)
                (left as Double) * (right as Double)
            }

            TokenType.SLASH -> {
                checkNumberOperands(expr.operator, left, right)
                (left as Double) / (right as Double)
            }

            TokenType.GREATER -> {
                checkNumberOperands(expr.operator, left, right)
                (left as Double) > (right as Double)
            }

            TokenType.GREATER_EQUAL -> {
                checkNumberOperands(expr.operator, left, right)
                (left as Double) >= (right as Double)
            }

            TokenType.LESS -> {
                checkNumberOperands(expr.operator, left, right)
                (left as Double) < (right as Double)
            }

            TokenType.LESS_EQUAL -> {
                checkNumberOperands(expr.operator, left, right)
                (left as Double) <= (right as Double)
            }

            TokenType.BANG_EQUAL -> return isEqual(left, right).not()
            TokenType.EQUAL_EQUAL -> return isEqual(left, right)
            TokenType.PLUS -> {
                if (left is Double && right is Double) {
                    return left + right
                } else if (left is String && right is String) {
                    return left + right
                } else {
                    throw RuntimeError(expr.operator, "operands must be Numbers or Strings")
                }
            }

            else -> null
        }
    }

    override fun visitCallExpr(expr: Expr.Call): Any? {
        val callee = evaluate(expr.callee)
        val args: List<Any?> = expr.params.map { evaluate(it) }

        if (callee !is LoxCallable) {
            throw RuntimeError(expr.paren, "Can only call functions and classes")
        }

        val function = callee

        if (args.size != function.arity()) {
            throw RuntimeError(expr.paren, "Expected ${function.arity()} argument but got ${args.size}")
        }

        return function.call(this, args)
    }

    override fun visitGetExpr(expr: Expr.Get): Any? {
        val obj = evaluate(expr.obj)
        if (obj is LoxInstance) {
            return obj.get(expr.name)
        }

        throw RuntimeError(expr.name, "Only instances have properties")
    }

    private fun isEqual(left: Any?, right: Any?): Boolean {
        return if (left == null && right == null) true
        else left?.equals(right) ?: false
    }

    override fun visitGroupingExpr(expr: Expr.Grouping): Any? = evaluate(expr.expression)
    override fun visitLiteralExpr(expr: Expr.Literal): Any? = expr.value
    override fun visitLogicalExpr(expr: Expr.Logical): Any? {
        val left = evaluate(expr.left)

        if (expr.operator.type == TokenType.OR) {
            if (isTruthy(left)) return left
        } else if (!isTruthy(left)) {
            return left
        }

        return evaluate(expr.right)
    }

    override fun visitSetExpr(expr: Expr.Set): Any? {
        val obj = evaluate(expr.obj)
        if (obj !is LoxInstance) {
            throw RuntimeError(expr.name, "Only instances have fields")
        }
        val value = evaluate(expr.value)
        obj.set(expr.name, value)
        return value
    }

    override fun visitSuperExpr(expr: Expr.Super): Any? {
        val distance = locals[expr]
        val superclass : LoxClass = environment.getAt(distance as Int, "super") as LoxClass

        val obj : LoxInstance = environment.getAt(distance - 1, "this") as LoxInstance

        val method = superclass.findMethod(expr.method.lexeme)
        return method?.bind(obj)
    }

    override fun visitThisExpr(expr: Expr.This): Any? {
        return lookupVariable(expr.keyword, expr)
    }

    override fun visitVariableExpr(expr: Expr.Variable): Any? = lookupVariable(expr.name, expr)

    override fun visitUnaryExpr(expr: Expr.Unary): Any? {
        val right = evaluate(expr.right)
        return when (expr.operator.type) {
            TokenType.MINUS -> {
                checkNumberOperand(expr.operator, right)
                -(right as Double)
            }

            TokenType.BANG -> !isTruthy(right)
            else -> null
        }
    }

    private fun isTruthy(obj: Any?): Boolean {
        return when (obj) {
            null -> false
            is Boolean -> obj
            else -> true
        }
    }

    private fun checkNumberOperand(operator: Token, obj: Any?) {
        if (obj !is Double) {
            throw RuntimeError(operator, "Operand must be a number")
        }
    }

    private fun checkNumberOperands(operator: Token, left: Any?, right: Any?) {
        if (left !is Double && right !is Double) {
            throw RuntimeError(operator, "Operands must be a numbers")
        }
    }

    private fun evaluate(expr: Expr): Any? {
        return expr.accept(this)
    }

    private fun stringify(obj: Any?): String {
        return when (obj) {
            null -> "nil"
            is Double -> {
                var text = obj.toString()
                if (text.endsWith(".0")) {
                    text = text.substring(0, text.length - 2)
                }
                text
            }

            else -> obj.toString()
        }
    }

    fun executeBlock(stmts: List<Stmt?>, env: Environment): Any? {
        val previous = this.environment
        try {
            this.environment = env
            for (stmt in stmts) {
                if (stmt == null) throw RuntimeError(
                    Token(TokenType.LEFT_BRACE, "{", "{", -1),
                    "Null statement found in block"
                )
                execute(stmt)
            }
        } finally {
            this.environment = previous
        }
        return null
    }

    override fun visitBlockStmt(stmt: Stmt.Block): Any? = executeBlock(stmt.statements, Environment(environment))

    override fun visitClassStmt(stmt: Stmt.Class): Any? {
        val superclass = if (stmt.superclass != null) {
            val sc = evaluate(stmt.superclass)
            if (sc !is LoxClass) {
                throw RuntimeError(stmt.superclass.name, "Superclass must be a class")
            } else {
                sc
            }
        } else {
            null
        }
        environment.define(stmt.name.lexeme, null)

        if (stmt.superclass != null) {
            environment = Environment(environment)
            environment.define("super", superclass)
        }

        val methods = mutableMapOf<String, LoxFunction>()
        for (method in stmt.methods) {
            val function = LoxFunction(method, environment, method.name.lexeme == "init")
            methods.set(method.name.lexeme, function)
        }

        val klass = LoxClass(stmt.name.lexeme, superclass, methods)
        environment.assign(stmt.name, klass)
        return null
    }

    override fun visitExpressionStmt(stmt: Stmt.Expression): Any? {
        evaluate(stmt.expression)
        return null
    }

    override fun visitFunctionStmt(stmt: Stmt.Function): Any? {
        val function = LoxFunction(stmt, environment, false)
        environment.define(stmt.name.lexeme, function)
        return null
    }

    override fun visitIfStmt(stmt: Stmt.If): Any? {
        val condition = evaluate(stmt.condition)

        if (isTruthy(condition)) {
            execute(stmt.branch)
        } else if (stmt.elseBranch != null) {
            execute(stmt.elseBranch)
        }

        return null
    }

    override fun visitVarStmt(stmt: Stmt.Var): Any? {
        val value = if (stmt.initializer == null) null else evaluate(stmt.initializer)
        environment.define(stmt.name.lexeme, value)
        return null
    }

    override fun visitPrintStmt(stmt: Stmt.Print): Any? {
        val value = this.evaluate(stmt.expression)
        println(value)
        return null
    }

    override fun visitReturnStmt(stmt: Stmt.Return): Any? {
        val value = if (stmt.value == null) null else evaluate(stmt.value)
        throw Return(value)
    }

    override fun visitWhileStmt(stmt: Stmt.While): Any? {
        while (isTruthy(evaluate(stmt.condition))) {
            execute(stmt.body)
        }

        return null
    }

    private fun lookupVariable(name: Token, expr: Expr): Any? {
        val distance = locals.get(expr)
        return if (distance != null) {
            environment.getAt(distance, name.lexeme)
        } else {
            globals.get(name)
        }

    }
}