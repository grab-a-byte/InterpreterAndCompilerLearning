import 'package:monkey_intepreter/abstract_syntax_tree/ast.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/expressions/identifier_expression.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/expressions/integer_expression.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/expression_statement.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/let_statement.dart';
import 'package:monkey_intepreter/abstract_syntax_tree/statements/return_statement.dart';
import 'package:monkey_intepreter/parser/precedences.dart';
import 'package:monkey_intepreter/token.dart';

import '../lexer.dart';

typedef PrefixFunction = Expression? Function();

typedef InfixFunction = Expression? Function(Expression);

class Parser {
  final Lexer _lexer;
  Token _currentToken;
  Token _peekToken;

  final Map<Type, PrefixFunction> _prefixFunctions = {};
  final Map<Type, InfixFunction> _infixFunctions = {};

  Parser(this._lexer)
      : _currentToken = _lexer.nextToken(),
        _peekToken = _lexer.nextToken() {
    // Map Prefix functions
    _prefixFunctions[Identifier] = _parseIdentiferExpression;
    _prefixFunctions[Integer] = _parseIntegerExpression;

    //Map Postfix functions
  }

  Program parseProgram() {
    var program = Program();

    while (_currentToken.runtimeType != EndOfFile) {
      var statement = _parseStatement();
      if (statement != null) {
        program.statements.add(statement);
      }
      _nextToken();
    }

    return program;
  }

  // Parsing Helpers
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

  Precedence currentPrecedence() =>
      precedenceMap[_currentToken.runtimeType] ?? Precedence.lowest;

  Precedence peekPrecedence() =>
      precedenceMap[_peekToken.runtimeType] ?? Precedence.lowest;

  //Statement Parsers
  Statement? _parseStatement() {
    switch (_currentToken.runtimeType) {
      case Let:
        return _parseLetStatement();
      case Return:
        return _parseReturnStatement();
      default:
        return _parseExpressionStatement();
    }
  }

  ExpressionStatement? _parseExpressionStatement() {
    Expression? exp = _parseExpression(Precedence.lowest);
    if (exp == null) {
      return null;
    }
    return ExpressionStatement(exp);
  }

  ReturnStatement? _parseReturnStatement() {
    _nextToken();

    Expression? expresison = _parseExpression(Precedence.lowest);
    if (expresison == null) {
      return null;
    }

    return ReturnStatement(expresison);
  }

  LetStatement? _parseLetStatement() {
    if (!_expectPeek(Identifier)) {
      return null;
    }

    var ident = IdentifierExpression(_currentToken.literal);

    if (!_expectPeek(Assign)) {
      return null;
    }

    _nextToken();
    var expression = _parseExpression(Precedence.lowest);

    if (expression == null) {
      return null;
    }
    return LetStatement(ident, expression);
  }

  //Expression Parsing
  Expression? _parseExpression(Precedence precedence) {
    var prefixFunc = _prefixFunctions[_currentToken.runtimeType];
    if (prefixFunc == null) return null;

    var left = prefixFunc();

    if (left == null) return null;

    if (_peekToken.runtimeType != SemiColon &&
        precedence.index < peekPrecedence().index) {
      var infixFunc = _infixFunctions[_peekToken.runtimeType];
      if (infixFunc == null) return left;
      _nextToken();
      left = infixFunc(left);
    }

    return left;
  }

  // Prefix Parsers
  Expression? _parseIdentiferExpression() {
    return IdentifierExpression(_currentToken.literal);
  }

  Expression? _parseIntegerExpression() {
    var value = int.tryParse(_currentToken.literal, radix: 10);
    if (value == null) return null;
    return IntegerExpression(value);
  }

  //Postfix Parsers
}
