//package lox
//
//import expressions.Expr
//
//class AstPrinter : Expr.Visitor<String> {
//    override fun visitBinaryExpr(expr: Expr.Binary): String =parenthesize(expr.operator.lexeme, expr.left, expr.right)
//    override fun visitGroupingExpr(expr: Expr.Grouping): String = parenthesize("group", expr.expression)
//    override fun visitLiteralExpr(expr: Expr.Literal): String = expr.value?.toString() ?: "nil"
//    override fun visitVariableExpr(expr: Expr.Variable): String {
//        TODO("Not yet implemented")
//    }
//
//    override fun visitUnaryExpr(expr: Expr.Unary): String = parenthesize(expr.operator.lexeme, expr.right)
//
//    fun print(expr: Expr): String {
//        return expr.accept(this)
//    }
//
//    private fun parenthesize(name: String, vararg exprs: Expr) : String{
//        val builder = StringBuilder()
//        builder.append("($name")
//        for (expr in exprs) {
//            builder.append(" ")
//            builder.append(expr.accept(this))
//        }
//        builder.append(")")
//        return builder.toString()
//    }
//}