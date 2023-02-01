package evaluator

import "monkey/object"

var builtIns = map[string]*object.BuiltIn{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("too many args for builtin function 'len'")
			}

			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(obj.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(obj.Items))}
			default:
				return newError("Unsupported Type for builtin function 'len'")
			}
		},
	},
}
