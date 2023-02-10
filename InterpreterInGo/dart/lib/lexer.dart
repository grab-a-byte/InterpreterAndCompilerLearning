import 'package:monkey_intepreter/constants.dart';
import 'package:monkey_intepreter/token.dart';

class Lexer {
  final String _input;

  int _readPosition = 0;

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
    ';': SemiColon.new
  };

  Lexer(this._input);

  Token nextToken() {
    _skipWhitespace();
    if (_readPosition >= _input.length) return EndOfFile();

    Token Function()? tokFn = _singleCharTokens[_input[_readPosition]];
    if (tokFn == null) {
      return Illegal();
    }

    _readPosition += 1;

    return tokFn();
  }

  void _skipWhitespace() {
    var currChar = _input[_readPosition];
    while (Constants.whitespace.contains(currChar) &&
        _readPosition < _input.length) {
      _readPosition++;
      if (_readPosition >= _input.length) break;
      currChar = _input[_readPosition];
    }
  }
}
