//import lox.AstPrinter
import lox.Lox
import kotlin.system.exitProcess

fun main(args: Array<String>) {
//    val expr : Expr = Expr.Binary(
//        Expr.Unary(
//            Token(TokenType.MINUS, "-", null, 1),
//            Expr.Literal(123),
//        ),
//        Token(TokenType.STAR, "*", null, 1),
//        Expr.Grouping(Expr.Literal(45.67))
//    )
//
//    val printer = AstPrinter()
//
//    println(printer.print(expr))

    if (args.size > 1) {
        println("Usage: klox [script]")
        exitProcess(64)
    } else if (args.size == 1) {
        Lox.runFile(args[0])
    } else {
        Lox.runPrompt()
    }
}