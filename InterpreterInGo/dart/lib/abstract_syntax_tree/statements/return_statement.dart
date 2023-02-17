import 'package:monkey_intepreter/abstract_syntax_tree/ast.dart';

class ReturnStatement extends Statement {
  final Expression value;

  ReturnStatement(this.value);

  @override
  String tokenLiteral() => "return";
}
