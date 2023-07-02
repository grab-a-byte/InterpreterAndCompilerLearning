package tool

import java.util.*
import kotlin.io.path.Path
import kotlin.system.exitProcess

fun main(args: Array<String>) {
    if (args.size != 1) {
        System.err.println("Usage: generate_ast <output_directory>");
        exitProcess(64)
    }

    val outputDir: String = args[0]
    defineAst(outputDir, "Expr", listOf(
        "Binary: Expr left, Token operator, Expr right",
        "Grouping: Expr expression",
        "Literal: Any? value",
        "Variable: Token name",
        "Unary: Token operator, Expr right"
    ))

    defineAst(outputDir, "Stmt", listOf(
        "Expression: Expr expression",
        "Var: Token name, Expr? initializer",
        "Print: Expr expression",
    ))
}

private fun defineAst(
    outputDir: String, baseName: String, types: List<String>
) {
    val builder = StringBuilder()
    builder.appendLine("package expressions")
    builder.appendLine()
    builder.appendLine("import lox.Token")
    builder.appendLine()
    builder.appendLine("abstract class $baseName {")
    builder.appendLine()

    //Base type methods
    builder.appendLine("    abstract fun <R> accept(visitor: Visitor<R>) : R")

    //Define visitor interface
    val visitor = defineVisitor(baseName, types)
    builder.appendLine(visitor)

    //Define each type
    types.forEach { type ->
        val className = type.split(":")[0].trim()
        val fields = type.split(":")[1].trim()
        val typeDef = defineType(baseName, className, fields)
        builder.appendLine(typeDef)
    }

    builder.appendLine("}")

    val file = Path(outputDir, "$baseName.kt").toAbsolutePath().toFile()
    println(file.path)
    file.createNewFile()
    file.writeText(builder.toString())
}

private fun defineType(baseName: String, className: String, fieldsList: String): String {
    val str = StringBuilder()
    str.appendLine("    class $className (")

    val fields = fieldsList.split(", ")
    fields.forEachIndexed { index, it ->
        val type = it.split(" ")[0]
        val name = it.split(" ")[1]
        if (index == fields.size - 1) str.append("        val $name: $type")
        else str.appendLine("        val $name: $type,")
    }
    str.appendLine(") : $baseName() {")
    str.appendLine("            override fun <R> accept(visitor: Visitor<R>): R {")
    str.appendLine("                return visitor.visit$className$baseName(this)")
    str.appendLine("            }")
    str.append("    }")
    return str.toString()
}

private fun defineVisitor(baseName: String, types: List<String>): String {
    val builder = StringBuilder()
    builder.appendLine("    interface Visitor<R> {")

    for (type in types) {
        val typeName = type.split(":")[0]
        builder.appendLine("        fun visit$typeName$baseName(${baseName.lowercase(Locale.getDefault())} : $typeName): R")
    }

    builder.appendLine("    }")

    return builder.toString()
}