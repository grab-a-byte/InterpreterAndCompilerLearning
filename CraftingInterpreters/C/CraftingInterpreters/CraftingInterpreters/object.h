#ifndef clox_object_h
#define clox_object_h

#include "common.h"
#include "value.h"

#define OBJECT_TYPE(obj) (AS_OBJ(obj)->type)
#define IS_STRING(obj) isObjType(obj, OBJ_STRING)
#define AS_STRING(obj) ((ObjString*)AS_OBJ(obj))
#define AS_CSTRING(obj) (((ObjString*)AS_OBJ(obj))->chars)

typedef enum {
	OBJ_STRING
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

ObjString* takeString(char* chars, int length);
ObjString* copyString(const char* chars, int length);

static inline bool isObjType(Value value, ObjType type) {
	return IS_OBJ(value) && AS_OBJ(value)->type == type;
}
#endif
