import 'package:monkey_intepreter/abstract_syntax_tree/expressions/integer_expression.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/let_statement.dart';
import 'package:monkey_intepreter/lexer.dart';
import 'package:monkey_intepreter/parser.dart';
import 'package:test/test.dart';

void main() {
  test("Test Let Expresison", () {
    var input = "let x = 5";
    var lexer = Lexer(input);
    var parser = Parser(lexer);

    var program = parser.parseProgram();
    expect(program != null, true);

    var statement = program!.statements[0];

    expect(statement.runtimeType, LetStatement);
    var letStatement = statement as LetStatement;

    expect(letStatement.tokenLiteral(), "let");
    expect(letStatement.identifier.literal, "x");
    expect(letStatement.value.runtimeType, IntegerExpression);
  });
}
