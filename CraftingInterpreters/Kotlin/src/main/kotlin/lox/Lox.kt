package lox

import java.io.File
import kotlin.math.exp
import kotlin.system.exitProcess

object Lox {

    private var hadError: Boolean = false
    private var hadRuntimeError: Boolean = false
    private val interpreter : Interpreter = Interpreter()

    private fun run(source: String){
        val scanner = Scanner(source)
        val tokens = scanner.scanTokens()
        val parser = Parser(tokens)
        val expr = parser.parse()

        if(hadError || expr == null) return

        val result = interpreter.interpret(expr)
        //println(AstPrinter().print(expr))
    }

    fun runFile(file: String){
        val bytes: ByteArray = File(file).readBytes()
        run(bytes.toString())
        if (hadError) exitProcess(65)
        if (hadRuntimeError) exitProcess(70)
    }

    fun runPrompt(){
        while (true) {
            print("> ")
            val line: String = readlnOrNull() ?: break
            run(line)
            hadError = false
        }
    }

    fun error(line: Int, message: String) {
        report(line, "", message)
    }

    fun error(tok: Token, message: String) {
        if (tok.type == TokenType.EOF){
            report(tok.line, " at end ", message)
        } else {
            report(tok.line, " at '${tok.lexeme}' ", message)
        }
    }

    private fun report(line: Int, where: String, message: String){
        println("[line $line] Error $where : $message")
        hadError = true
    }

    fun runtimeError(e: RuntimeError) {
        System.err.println("${e.message} \n [line: ${e.token.line}]")
        hadRuntimeError = true
    }
}