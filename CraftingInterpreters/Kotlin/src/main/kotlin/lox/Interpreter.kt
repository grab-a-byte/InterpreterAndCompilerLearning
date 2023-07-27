package lox

import expressions.Expr
import expressions.Stmt

class Interpreter : Expr.Visitor<Any?>, Stmt.Visitor<Any?> {

    private var environment = Environment()

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
        environment.assign(expr.name, value)
        return value
    }

    override fun visitBinaryExpr(expr: Expr.Binary): Any? {
        val left = evaluate(expr.left)
        val right = evaluate(expr.right)
        return when(expr.operator.type) {
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

    override fun visitVariableExpr(expr: Expr.Variable): Any? = environment.get(expr.name)

    override fun visitUnaryExpr(expr: Expr.Unary): Any? {
        val right = evaluate(expr.right)
        return when(expr.operator.type) {
            TokenType.MINUS -> {
                checkNumberOperand(expr.operator, right)
                -(right as Double)
            }
            TokenType.BANG -> !isTruthy(right)
            else -> null
        }
    }

    private fun isTruthy(obj : Any?) : Boolean {
        return when(obj) {
            null -> false
            is Boolean -> obj
            else -> true
        }
    }

    private fun checkNumberOperand(operator: Token, obj : Any?) {
        if (obj !is Double) {
            throw RuntimeError(operator, "Operand must be a number")
        }
    }

    private fun checkNumberOperands(operator: Token, left: Any?, right: Any?) {
        if (left !is Double && right !is Double) {
            throw RuntimeError(operator, "Operands must be a numbers")
        }
    }

    private fun evaluate(expr: Expr) : Any? {
        return expr.accept(this)
    }

    private fun stringify(obj: Any?) : String {
        return when (obj) {
            null -> "nil"
            is Double -> {
                var text = obj.toString()
                if (text.endsWith(".0")) {
                    text = text.substring(0, text.length -2)
                }
                text
            }

            else -> obj.toString()
        }
    }

    override fun visitBlockStmt(stmt: Stmt.Block): Any? {
        val previous = this.environment
        try {
            this.environment = Environment(previous)
            for (stmt in stmt.statements) {
                if (stmt == null) throw RuntimeError(Token(TokenType.LEFT_BRACE, "{", "{", -1) ,"Null statement found in block")
                execute(stmt)
            }
        } finally {
            this.environment = previous
        }
        return null
    }

    override fun visitExpressionStmt(stmt: Stmt.Expression): Any? {
        evaluate(stmt.expression)
        return null
    }

    override fun visitIfStmt(stmt: Stmt.If): Any? {
        val condition = evaluate(stmt.condition)

        if (isTruthy(condition)) {
            execute(stmt.branch)
        } else if(stmt.elseBranch != null) {
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

    override fun visitWhileStmt(stmt: Stmt.While): Any? {
        while (isTruthy(evaluate(stmt.condition))) {
            execute(stmt.body)
        }

        return null
    }
}