import 'package:monkey_intepreter/constants.dart';
import 'package:monkey_intepreter/token.dart';

class Lexer {
  final RegExp _alphaRegex = RegExp(r'^[a-zA-Z]$');
  final String _input;

  int _readPosition = 0;

  Lexer(this._input);

  Token nextToken() {
    _skipWhitespace();
    if (_readPosition >= _input.length) return EndOfFile();

    Token Function()? tokFn = _singleCharTokens[_currentChar()];
    if (tokFn != null) {
      _readPosition += 1;
      return tokFn();
    }

    if (_isAlpha(_currentChar())) {
      var literal = _parseIdentifier();
      var keywordFunc = _keywords[literal];
      return keywordFunc == null ? Identifier(literal) : keywordFunc(literal);
    }

    if (_isNumber(_currentChar())) {
      var literal = _parseNumber();
      return Number(literal);
    }

    _readPosition += 1;

    return Illegal();
  }

  String? _currentChar() =>
      _readPosition >= _input.length ? null : _input[_readPosition];

  bool _isNumber(String? c) => c == null ? false : double.tryParse(c) != null;
  bool _isAlpha(String? c) => c == null ? false : _alphaRegex.hasMatch(c);

  void _skipWhitespace() {
    var currChar = _input[_readPosition];
    while (Constants.whitespace.contains(currChar) &&
        _readPosition < _input.length) {
      _readPosition++;
      if (_readPosition >= _input.length) break;
      currChar = _input[_readPosition];
    }
  }

  String _parseIdentifier() {
    var startPosition = _readPosition;
    while (_isAlpha(_currentChar())) {
      _readPosition += 1;
    }
    return _input.substring(startPosition, _readPosition);
  }

  String _parseNumber() {
    var firstPart = parseNumeric();
    if (_currentChar() != ".") return firstPart;
    _readPosition += 1;
    var secondPart = parseNumeric();
    return "$firstPart.$secondPart";
  }

  String parseNumeric() {
    var startPosition = _readPosition;
    while (_isNumber(_currentChar())) {
      _readPosition += 1;
    }
    return _input.substring(startPosition, _readPosition);
  }

  final Map<String, Token Function(String)> _keywords = {
    "func": Func.new,
    "let": Let.new
  };

  final Map<String, Token Function()> _singleCharTokens = {
    '+': Plus.new,
    '-': Dash.new,
    '{': LeftBrace.new,
    '}': RightBrace.new,
    '(': LeftParen.new,
    ')': RightParen.new,
    '[': LeftBracket.new,
    ']': RightBracket.new,
    ',': Comma.new,
    ':': Colon.new,
    ';': SemiColon.new,
    '=': Assign.new
  };
}
