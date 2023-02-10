import 'package:monkey_intepreter/lexer.dart';
import 'package:monkey_intepreter/token.dart';
import 'package:test/test.dart';

void main() {
  test('When lexing token, generates correct tokens', () {
    var input = '''
    (){}-+[],:;
    ''';

    Lexer lexer = Lexer(input);
    List<Type> expected = [
      LeftParen,
      RightParen,
      LeftBrace,
      RightBrace,
      Dash,
      Plus,
      LeftBracket,
      RightBracket,
      Comma,
      Colon,
      SemiColon,
      EndOfFile
    ];

    for (var token in expected) {
      Token producedToken = lexer.nextToken();
      expect(producedToken.runtimeType, token);
    }
  });
}
