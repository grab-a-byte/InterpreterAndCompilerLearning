abstract class Token {
  final String literal;
  Token(this.literal);
}

class Illegal extends Token {
  Illegal() : super("ILLEGAL");
}

class Plus extends Token {
  Plus() : super("+");
}

class Dash extends Token {
  Dash() : super("-");
}

class LeftParen extends Token {
  LeftParen() : super("(");
}

class RightParen extends Token {
  RightParen() : super(")");
}

class LeftBrace extends Token {
  LeftBrace() : super("{");
}

class RightBrace extends Token {
  RightBrace() : super("}");
}

class RightBracket extends Token {
  RightBracket() : super("]");
}

class LeftBracket extends Token {
  LeftBracket() : super("]");
}

class Comma extends Token {
  Comma() : super(",");
}

class SemiColon extends Token {
  SemiColon() : super(";");
}

class Colon extends Token {
  Colon() : super(":");
}

class EndOfFile extends Token {
  EndOfFile() : super("EOF");
}
