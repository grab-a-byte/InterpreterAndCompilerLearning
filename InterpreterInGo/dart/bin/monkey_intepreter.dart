import 'package:monkey_intepreter/lexer.dart';
import 'package:monkey_intepreter/token.dart';

void main(List<String> arguments) {
  Lexer l = Lexer("()");

  var t = l.nextToken();

  while (t.runtimeType != EndOfFile) {
    print(t.runtimeType);
    t = l.nextToken();
  }
}
