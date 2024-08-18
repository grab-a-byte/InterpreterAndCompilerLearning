#ifndef clox_object_h
#define clox_object_h

#include "chunk.h"
#include "common.h"
#include "table.h"
#include "value.h"

#define OBJECT_TYPE(obj) (AS_OBJ(obj)->type)

#define IS_STRING(obj) isObjType(obj, OBJ_STRING)
#define IS_FUNCTION(obj) isObjType(obj, OBJ_FUNCTION)
#define IS_CLOSURE(obj) isObjType(obj, OBJ_CLOSURE)
#define IS_CLASS(obj) isObjType(obj, OBJ_CLASS)
#define IS_INSTANCE(obj) isObjType(obj, OBJ_INSTANCE)

#define AS_STRING(obj) ((ObjString *)AS_OBJ(obj))
#define AS_CSTRING(obj) (((ObjString *)AS_OBJ(obj))->chars)
#define AS_FUNCTION(obj) ((ObjFunction *)AS_OBJ(obj))
#define AS_NATIVE(obj) (((ObjNative *)AS_OBJ(obj))->function)
#define AS_CLOSURE(obj) ((ObjClosure *)AS_OBJ(obj))
#define AS_CLASS(obj) ((ObjClass *)AS_OBJ(obj))
#define AS_INSTANCE(obj) ((ObjInstance *)AS_OBJ(obj))

typedef enum
{
  OBJ_STRING,
  OBJ_FUNCTION,
  OBJ_NATIVE,
  OBJ_CLOSURE,
  OBJ_UPVALUE,
  OBJ_CLASS,
  OBJ_INSTANCE,
} ObjType;

struct Obj
{
  ObjType type;
  bool isMarked;
  struct Obj *next;
};

struct ObjString
{
  Obj obj;
  int length;
  char *chars;
  uint32_t hash;
};

typedef struct
{
  Obj obj;
  ObjString *name;
  int arity;
  Chunk chunk;
  int upvalueCount;
} ObjFunction;

typedef struct ObjUpvalue
{
  Obj obj;
  Value *location;
  struct ObjUpvalue *next;
  Value closed;
} ObjUpvalue;

typedef struct
{
  Obj obj;
  ObjFunction *function;
  ObjUpvalue **upvalues;
  int upvalueCount;
} ObjClosure;

typedef struct
{
  Obj obj;
  ObjString *name;
} ObjClass;

typedef Value (*NativeFn)(int argCount, Value *args);

typedef struct
{
  Obj obj;
  NativeFn function;
} ObjNative;

typedef struct
{
  Obj obj;
  ObjClass *klass;
  Table fields;
} ObjInstance;

ObjString *takeString(char *chars, int length);
ObjString *copyString(const char *chars, int length);
void printObject(Value value);
ObjFunction *newFunction();
ObjNative *newNative(NativeFn func);
ObjClosure *newClosure(ObjFunction *function);
ObjUpvalue *newUpvalue(Value *slot);
ObjClass *newClass(ObjString *name);
ObjInstance* newInstance(ObjClass* klass);

static inline bool isObjType(Value value, ObjType type)
{
  return IS_OBJ(value) && AS_OBJ(value)->type == type;
}
#endif
