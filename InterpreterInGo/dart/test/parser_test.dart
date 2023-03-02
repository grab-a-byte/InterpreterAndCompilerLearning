import 'package:monkey_intepreter/abstract_syntax_tree/expressions/integer_expression.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/expression_statement.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/let_statement.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/return_statement.dart';
import 'package:monkey_intepreter/lexer.dart';
import 'package:monkey_intepreter/parser/parser.dart';
import 'package:test/test.dart';

void main() {
  test("Test Let Statement", () {
    var input = "let x = 5";
    var lexer = Lexer(input);
    var parser = Parser(lexer);

    var program = parser.parseProgram();
    expect(program.statements.length, 1);

    var statement = program.statements[0];

    expect(statement.runtimeType, LetStatement);
    var letStatement = statement as LetStatement;

    expect(letStatement.tokenLiteral(), "let");
    expect(letStatement.identifier.literal, "x");
    expect(letStatement.value.runtimeType, IntegerExpression);
  });

  test("Test Return Statement", () {
    var input = "return 5";
    var lexer = Lexer(input);
    var parser = Parser(lexer);

    var program = parser.parseProgram();
    expect(program.statements.length, 1);

    var statement = program.statements[0];

    expect(statement.runtimeType, ReturnStatement);
    var returnStatement = statement as ReturnStatement;

    expect(returnStatement.tokenLiteral(), "return");
    expect(returnStatement.value.runtimeType, IntegerExpression);
  });

  test("Test Expression Statement", () {
    var input = "5";
    var lexer = Lexer(input);
    var parser = Parser(lexer);

    var program = parser.parseProgram();
    expect(program.statements.length, 1);

    var statement = program.statements[0];

    expect(statement.runtimeType, ExpressionStatement);
    var expressionStatement = statement as ExpressionStatement;

    expect(expressionStatement.tokenLiteral(), "5");
    expect(expressionStatement.expression.runtimeType, IntegerExpression);
  });
}
