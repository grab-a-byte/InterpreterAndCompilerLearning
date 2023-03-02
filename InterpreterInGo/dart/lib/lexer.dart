import 'package:monkey_intepreter/constants.dart';
import 'package:monkey_intepreter/token.dart';

class Lexer {
  final RegExp _alphaRegex = RegExp(r'^[a-zA-Z]$');
  final String _input;

  int _readPosition = 0;

  Lexer(this._input);

  Token nextToken() {
    if (_readPosition >= _input.length) return EndOfFile();

    _skipWhitespace();

    if (_readPosition >= _input.length) return EndOfFile();

    Token? token;

    switch (_currentChar()) {
      case "+":
        token = Plus();
        break;
      case '-':
        token = Minus();
        break;
      case '*':
        token = Star();
        break;
      case '/':
        token = Slash();
        break;
      case '{':
        token = LeftBrace();
        break;
      case '}':
        token = RightBrace();
        break;
      case '(':
        token = LeftParen();
        break;
      case ')':
        token = RightParen();
        break;
      case '[':
        token = LeftBracket();
        break;
      case ']':
        token = RightBracket();
        break;
      case ',':
        token = Comma();
        break;
      case ':':
        token = Colon();
        break;
      case ';':
        token = SemiColon();
        break;
      case '=':
        if (_nextChar() == "=") {
          _readPosition += 1;
          token = Equals();
          break;
        }
        token = Assign();
        break;
      case '!':
        if (_nextChar() == "=") {
          _readPosition += 1;
          token = NotEquals();
          break;
        }
        token = Bang();
        break;
      default:
        if (_isAlpha(_currentChar())) {
          var literal = _readIdentifier();
          var keywordFunc = _keywords[literal];
          return keywordFunc == null
              ? Identifier(literal)
              : keywordFunc(literal);
        }

        if (_isNumber(_currentChar())) {
          token = _readNumber();
        }
    }

    _readPosition += 1;

    return token ?? Illegal();
  }

  String? _currentChar() =>
      _readPosition >= _input.length ? null : _input[_readPosition];

  String? _nextChar() =>
      _readPosition + 1 >= _input.length ? null : _input[_readPosition + 1];

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

  String _readIdentifier() {
    var startPosition = _readPosition;
    while (_isAlpha(_currentChar())) {
      _readPosition += 1;
    }
    return _input.substring(startPosition, _readPosition);
  }

  Token _readNumber() {
    var firstPart = _parseNumeric();
    if (_currentChar() != ".") return Integer(firstPart);
    _readPosition += 1;
    var secondPart = _parseNumeric();
    return Float("$firstPart.$secondPart");
  }

  String _parseNumeric() {
    var startPosition = _readPosition;
    while (_isNumber(_currentChar())) {
      _readPosition += 1;
    }
    return _input.substring(startPosition, _readPosition);
  }

  final Map<String, Token Function(String)> _keywords = {
    "func": Func.new,
    "let": Let.new,
    "if": If.new,
    "return": Return.new
  };
}
