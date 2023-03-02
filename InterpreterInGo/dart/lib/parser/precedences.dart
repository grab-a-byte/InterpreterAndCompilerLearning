import 'package:monkey_intepreter/token.dart';

enum Precedence { lowest, equals, lessGreater, sum, product, prefix, call }

final Map<Type, Precedence> precedenceMap = {
  Equals: Precedence.equals,
  NotEquals: Precedence.equals,
  Plus: Precedence.sum,
  Minus: Precedence.sum,
  Star: Precedence.product,
  Slash: Precedence.product,
};
