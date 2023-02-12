import 'package:monkey_intepreter/abstract_syntax_tree/ast.dart';

class IdentifierExpression extends Expression {
  String literal;

  IdentifierExpression(this.literal);

  @override
  String tokenLiteral() => literal;
}
