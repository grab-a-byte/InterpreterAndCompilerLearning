package lox

import expressions.Expr

class Interpreter : Expr.Visitor<Any?> {
    override fun visitBinaryExpr(expr: Expr.Binary): Any? {
        val left = evalutate(expr.left)
        val right = evalutate(expr.right)
        return when(expr.operator.type) {
            TokenType.MINUS -> (left as Double) - (right as Double)
            TokenType.STAR -> (left as Double) * (right as Double)
            TokenType.SLASH -> (left as Double) / (right as Double)
            TokenType.GREATER -> (left as Double) > (right as Double)
            TokenType.GREATER_EQUAL -> (left as Double) >= (right as Double)
            TokenType.LESS -> (left as Double) < (right as Double)
            TokenType.LESS_EQUAL -> (left as Double) <= (right as Double)
            TokenType.BANG_EQUAL -> return !isEqual(left, right)
            TokenType.EQUAL_EQUAL -> return isEqual(left, right)
            TokenType.PLUS -> {
                if (left is Double && right is Double) {
                    return left + right
                } else if (left is String && right is String) {
                    return left + right
                } else {
                    return null
                }
            }
            else -> null
        }
    }

    private fun isEqual(left: Any?, right: Any?): Any? {
        return if (left == null && right == null) true
        else left?.equals(right) ?: false
    }

    override fun visitGroupingExpr(expr: Expr.Grouping): Any? = evalutate(expr.expression)
    override fun visitLiteralExpr(expr: Expr.Literal): Any? = expr.value

    override fun visitUnaryExpr(expr: Expr.Unary): Any? {
        val right = evalutate(expr.right)
        return when(expr.operator.type) {
            TokenType.MINUS -> -(right as Double)
            TokenType.BANG -> !isTruthy(right)
            else -> null
        }
    }

    private fun isTruthy(obj : Any?) : Boolean {
        return when(obj) {
            null -> false
            is Boolean -> !obj
            else -> true
        }
    }

    fun evalutate(expr: Expr) : Any? {
        return expr.accept(this)
    }
}