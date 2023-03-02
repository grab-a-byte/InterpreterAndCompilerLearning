import 'package:monkey_intepreter/lexer.dart';
import 'package:monkey_intepreter/token.dart';
import 'package:test/test.dart';

class TokenLiteralPair {
  Type type;
  String? literal;

  TokenLiteralPair(this.type, this.literal);
}

void main() {
  test('When lexing token, generates correct tokens', () {
    var input = '''
    (){}-+[],:;=
    hello
    123
    123.456
    func
    let x = 5
    ==
    !
    !=
    /*
    ''';

    Lexer lexer = Lexer(input);
    List<TokenLiteralPair> expected = [
      TokenLiteralPair(LeftParen, null),
      TokenLiteralPair(RightParen, null),
      TokenLiteralPair(LeftBrace, null),
      TokenLiteralPair(RightBrace, null),
      TokenLiteralPair(Minus, null),
      TokenLiteralPair(Plus, null),
      TokenLiteralPair(LeftBracket, null),
      TokenLiteralPair(RightBracket, null),
      TokenLiteralPair(Comma, null),
      TokenLiteralPair(Colon, null),
      TokenLiteralPair(SemiColon, null),
      TokenLiteralPair(Assign, null),
      TokenLiteralPair(Identifier, "hello"),
      TokenLiteralPair(Integer, "123"),
      TokenLiteralPair(Float, "123.456"),
      TokenLiteralPair(Func, null),
      TokenLiteralPair(Let, null),
      TokenLiteralPair(Identifier, "x"),
      TokenLiteralPair(Assign, null),
      TokenLiteralPair(Integer, "5"),
      TokenLiteralPair(Equals, null),
      TokenLiteralPair(Bang, null),
      TokenLiteralPair(NotEquals, null),
      TokenLiteralPair(Slash, null),
      TokenLiteralPair(Star, null),
      TokenLiteralPair(EndOfFile, null)
    ];

    for (var pair in expected) {
      Token producedToken = lexer.nextToken();
      expect(producedToken.runtimeType, pair.type);
      if (pair.literal != null) {
        expect(producedToken.literal, pair.literal);
      }
    }
  });
}
