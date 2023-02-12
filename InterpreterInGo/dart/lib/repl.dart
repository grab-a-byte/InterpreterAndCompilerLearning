import 'dart:io';

import 'package:monkey_intepreter/lexer.dart';
import 'package:monkey_intepreter/token.dart';

class REPL {
  void start() {
    stdout.writeln(
        "Welcome to the Monkey programming language, dart ediiton, type commands to start");

    while (true) {
      var line = stdin.readLineSync();
      if (line == null) continue;
      var lexer = Lexer(line);

      var token = lexer.nextToken();
      while (token.runtimeType != EndOfFile) {
        stdout.writeln("${token.runtimeType} : ${token.literal}");
        token = lexer.nextToken();
      }
    }
  }
}
