package expressions

import lox.Token

abstract class Expr {

    abstract fun <R> accept(visitor: Visitor<R>) : R
    interface Visitor<R> {
        fun visitBinaryExpr(expr : Binary): R
        fun visitGroupingExpr(expr : Grouping): R
        fun visitLiteralExpr(expr : Literal): R
        fun visitVariableExpr(expr : Variable): R
        fun visitUnaryExpr(expr : Unary): R
    }

    class Binary (
        val left: Expr,
        val operator: Token,
        val right: Expr) : Expr() {
            override fun <R> accept(visitor: Visitor<R>): R {
                return visitor.visitBinaryExpr(this)
            }
    }
    class Grouping (
        val expression: Expr) : Expr() {
            override fun <R> accept(visitor: Visitor<R>): R {
                return visitor.visitGroupingExpr(this)
            }
    }
    class Literal (
        val value: Any?) : Expr() {
            override fun <R> accept(visitor: Visitor<R>): R {
                return visitor.visitLiteralExpr(this)
            }
    }
    class Variable (
        val name: Token) : Expr() {
            override fun <R> accept(visitor: Visitor<R>): R {
                return visitor.visitVariableExpr(this)
            }
    }
    class Unary (
        val operator: Token,
        val right: Expr) : Expr() {
            override fun <R> accept(visitor: Visitor<R>): R {
                return visitor.visitUnaryExpr(this)
            }
    }
}
