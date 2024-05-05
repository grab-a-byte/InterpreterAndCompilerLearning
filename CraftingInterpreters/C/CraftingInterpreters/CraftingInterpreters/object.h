#ifndef clox_object_h
#define clox_object_h

#include "common.h"
#include "value.h"
#include "chunk.h"

#define OBJECT_TYPE(obj) (AS_OBJ(obj)->type)

#define IS_STRING(obj) isObjType(obj, OBJ_STRING)
#define IS_FUNCTION(obj) isObjType(obj, OBJ_FUNCTION)
#define IS_CLOSURE(obj) isObjType(obj, OBJ_CLOSURE)

#define AS_STRING(obj) ((ObjString*)AS_OBJ(obj))
#define AS_CSTRING(obj) (((ObjString*)AS_OBJ(obj))->chars)
#define AS_FUNCTION(obj) ((ObjFunction*)AS_OBJ(obj))
#define AS_NATIVE(obj) (((ObjNative*)AS_OBJ(obj))->function)
#define AS_CLOSURE(obj) ((ObjClosure*)AS_OBJ(obj))

typedef enum {
	OBJ_STRING,
	OBJ_FUNCTION,
	OBJ_NATIVE,
	OBJ_CLOSURE,
	OBJ_UPVALUE,
} ObjType;

struct Obj {
	ObjType type;
	struct Obj* next;
};

struct ObjString {
	Obj obj;
	int length;
	char* chars;
	uint32_t hash;
};

typedef struct {
	Obj obj;
	ObjString* name;
	int arity;
	Chunk chunk;
	int upvalueCount;
} ObjFunction;

typedef struct ObjUpvalue {
	Obj obj;
	Value* location;
} ObjUpvalue;

typedef struct {
	Obj obj;
	ObjFunction* function;
	ObjUpvalue** upvalues;
	int upvalueCount;
} ObjClosure;

typedef Value (*NativeFn)(int argCount, Value* args);

typedef struct {
	Obj obj;
	NativeFn function;
} ObjNative;


ObjString* takeString(char* chars, int length);
ObjString* copyString(const char* chars, int length);
void printObject(Value value);
ObjFunction* newFunction();
ObjNative* newNative(NativeFn func);
ObjClosure* newClosure(ObjFunction* function);
ObjUpvalue* newUpvalue(Value* slot);

static inline bool isObjType(Value value, ObjType type) {
	return IS_OBJ(value) && AS_OBJ(value)->type == type;
}
#endif
