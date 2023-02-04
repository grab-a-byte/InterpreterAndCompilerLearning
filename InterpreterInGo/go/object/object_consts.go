package object

const (
	INTEGER_OBJ           = "INTEGER"
	BOOLEAN_OBJ           = "BOOLEAN"
	NULL_OBJ              = "NULL"
	RETURN_VALUE_OBJECT   = "RETURN_VALUE"
	ERROR_VALUE_OBJECT    = "ERROR_VALUE"
	FUNCTION_VALUE_OBJECT = "FUNCTION"
	STRING_OBJECT         = "STRING_OBJECT"
	BUILT_IN_OBJECT       = "BUILT_IN_OBJECT"
	ARRAY_OBJECT          = "ARRAY_OBJECT"
	HASH_OBJECT           = "HASH_OBJECT"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)
