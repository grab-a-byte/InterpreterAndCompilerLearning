package lox

import expressions.Expr
import expressions.Stmt
import java.util.Stack

class Resolver(private val interpreter: Interpreter) : Expr.Visitor<Unit>, Stmt.Visitor<Unit> {

    private enum class FunctionType {
        None, Function, Method
    }

    private enum class ClassType {
        None, Class
    }

    private val scopes : Stack<MutableMap<String, Boolean>> = Stack()
    private var currentFunction = FunctionType.None
    private var currentClass = ClassType.None


    override fun visitAssignExpr(expr: Expr.Assign) {
        resolve(expr.value)
        resolveLocal(expr, expr.name)
    }

    override fun visitBinaryExpr(expr: Expr.Binary) {
        resolve(expr.left)
        resolve(expr.right)
    }

    override fun visitCallExpr(expr: Expr.Call) {
        resolve(expr.callee)
        for(arg in expr.params) {
            resolve(arg)
        }
    }

    override fun visitGetExpr(expr: Expr.Get) {
      resolve(expr.obj)
    }

    override fun visitGroupingExpr(expr: Expr.Grouping) {
        resolve(expr.expression)
    }

    override fun visitLiteralExpr(expr: Expr.Literal) {}

    override fun visitLogicalExpr(expr: Expr.Logical) {
        resolve(expr.left)
        resolve(expr.right)
    }

    override fun visitSetExpr(expr: Expr.Set) {
       resolve(expr.value)
       resolve(expr.obj)
    }

    override fun visitThisExpr(expr: Expr.This) {
        if (currentClass != ClassType.Class) {
            Lox.error(expr.keyword, "cannot use this outside of a class definition")
            return
        }
        resolveLocal(expr, expr.keyword)
    }

    override fun visitVariableExpr(expr: Expr.Variable) {
        if (!scopes.isEmpty() && scopes.peek().get(expr.name.lexeme) == false) {
            Lox.error(expr.name, "Can't read local variable in its own initializer.")
        }
        resolveLocal(expr, expr.name)
    }

    override fun visitUnaryExpr(expr: Expr.Unary) {
        resolve(expr.right)
    }

    override fun visitBlockStmt(stmt: Stmt.Block) {
        beginScope()
        resolve(stmt.statements)
        endScope()
    }

    override fun visitClassStmt(stmt: Stmt.Class) {
        val enclosingClass = currentClass
        currentClass = ClassType.Class

        declare(stmt.name)
        define(stmt.name)

        beginScope()
        scopes.peek().set("this", true)

        for (method in stmt.methods) {
            val declaration = FunctionType.Method
            resolveFunction(method, declaration)
        }

        endScope()

        currentClass = enclosingClass
    }

    override fun visitExpressionStmt(stmt: Stmt.Expression) {
        resolve(stmt.expression)
    }

    override fun visitFunctionStmt(stmt: Stmt.Function) {
        declare(stmt.name)
        define(stmt.name)
        resolveFunction(stmt, FunctionType.Function)
    }

    override fun visitIfStmt(stmt: Stmt.If) {
        resolve(stmt.condition)
        resolve(stmt.branch)
        if (stmt.elseBranch != null) resolve(stmt.elseBranch)
    }

    override fun visitVarStmt(stmt: Stmt.Var) {
        declare(stmt.name)
        if (stmt.initializer != null) {
            resolve(stmt.initializer)
        }
        define(stmt.name)
    }

    override fun visitPrintStmt(stmt: Stmt.Print) {
        resolve(stmt.expression)
    }

    override fun visitReturnStmt(stmt: Stmt.Return) {
        if (currentFunction == FunctionType.None) {
            Lox.error(stmt.keyword, "Can't return from top level code")
        }
        if (stmt.value != null) resolve(stmt.value)
    }

    override fun visitWhileStmt(stmt: Stmt.While) {
        resolve(stmt.condition)
        resolve(stmt.body)
    }

    fun declare(name: Token){
        if (scopes.isEmpty()) return
        val scope = scopes.peek()
        if (scope.containsKey(name.lexeme)) {
            Lox.error(name, "Already a variable with this name in scope")
        }
        scope.set(name.lexeme, false)
    }

    fun define(name: Token){
        if (scopes.isEmpty()) return
        scopes.peek().set(name.lexeme, true)
    }

    fun beginScope() {
        scopes.push(mutableMapOf())
    }

    fun endScope() {
        scopes.pop()
    }

    fun resolveLocal(expr: Expr, name: Token) {
        for (i in scopes.size - 1 downTo 0) {
            if (scopes.get(i).containsKey(name.lexeme)) {
                interpreter.resolve(expr, scopes.size - 1 - i)
            }
        }
    }

    private fun resolveFunction(func: Stmt.Function, functionType: FunctionType) {
        currentFunction = functionType
        beginScope()
        for (param in func.args) {
            declare(param)
            define(param)
        }
        resolve(func.body)
        endScope()
    }

    fun resolve(statements: List<Stmt?>)  {
        for (stmt in statements) {
            resolve(stmt)
        }
    }

    fun resolve(stmt: Stmt?)  {
        if (stmt == null) return
        stmt.accept(this)
    }

    fun resolve(expr: Expr?) {
        if (expr == null) return
        expr.accept(this)
    }
}