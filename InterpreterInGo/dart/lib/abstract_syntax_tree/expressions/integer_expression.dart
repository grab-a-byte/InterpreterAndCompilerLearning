import 'package:monkey_intepreter/abstract_syntax_tree/ast.dart';

class IntegerExpression extends Expression {
  final int value;

  IntegerExpression(this.value);

  @override
  String tokenLiteral() => "$value";
}
