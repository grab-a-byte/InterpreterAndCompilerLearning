import '../ast.dart';

class ExpressionStatement extends Statement {
  final Expression expression;

  ExpressionStatement(this.expression);

  @override
  String tokenLiteral() => expression.tokenLiteral();
}
