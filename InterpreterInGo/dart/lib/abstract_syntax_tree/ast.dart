abstract class Node {
  String tokenLiteral();
}

abstract class Expression extends Node {}

abstract class Statement extends Node {}

class Program extends Node {
  List<Statement> statements = <Statement>[];

  @override
  String tokenLiteral() =>
      statements.isNotEmpty ? statements[0].tokenLiteral() : "";
}
