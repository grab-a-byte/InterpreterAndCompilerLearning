import 'package:monkey_intepreter/abstract_syntax_tree/ast.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/expressions/identifier_expression.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/let_statement.dart';
import 'package:monkey_intepreter/token.dart';

import 'lexer.dart';

class Parser {
  final Lexer _lexer;
  Token _currentToken;
  Token _peekToken;

  Parser(this._lexer)
      : _currentToken = _lexer.nextToken(),
        _peekToken = _lexer.nextToken();

  Program? parseProgram() {
    var program = Program();
    while (_currentToken.runtimeType != EndOfFile) {
      var statement = _parseStatement();
      //May fail, will be fixed when parseStatment is not nullable
      program.statements.add(statement!);
    }

    return program;
  }

  //TODO make no-nullable
  Statement? _parseStatement() {
    switch (_currentToken.runtimeType) {
      case Let:
        return _parseLetStatement();
    }

    return null;
  }

  LetStatement? _parseLetStatement() {
    _nextToken();

    if (!_expectPeek(Identifier)) {
      return null;
    }

    var ident = IdentifierExpression(_currentToken.literal);

    if (!_expectPeek(Equals)) {
      return null;
    }

    _nextToken();
    var expression = _parseExpression();

    return LetStatement(ident, expression!);
  }

  Expression? _parseExpression() {
    //TODO return
    return null;
  }

  bool _expectPeek(Type tokenType) {
    if (_peekToken.runtimeType != tokenType) {
      return false;
    }

    _nextToken();
    return true;
  }

  void _nextToken() {
    _currentToken = _peekToken;
    _peekToken = _lexer.nextToken();
  }
}
