import 'package:monkey_intepreter/abstract_syntax_tree/ast.dart';

import '../expressions/identifier_expression.dart';

class LetStatement extends Statement {
  IdentifierExpression identifier;
  Expression value;

  LetStatement(this.identifier, this.value);

  @override
  String tokenLiteral() => "let";
}
