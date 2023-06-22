package lox

import java.io.File
import kotlin.system.exitProcess

object Lox {

    private var hadError: Boolean = false

    private fun run(source: String){
        val scanner = Scanner(source)
        val tokens = scanner.scanTokens()
        for (token in tokens) {
            println(token)
        }
    }

    fun runFile(file: String){
        val bytes: ByteArray = File(file).readBytes()
        run(bytes.toString())
        if (hadError) exitProcess(65)
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

    private fun report(line: Int, where: String, message: String){
        println("[line $line] Error $where : $message")
        hadError = true
    }
}