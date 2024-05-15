#include "memory.h"
#include "chunk.h"
#include "object.h"
#include "value.h"
#include "vm.h"
#include <stdlib.h>
#ifdef DEBUG_LOG_GC
#include <stdio.h>
#include "debug.h"
#endif

void *reallocate(void *pointer, size_t oldSize, size_t newSize) {

    if (newSize > oldSize) {
#ifdef DEBUG_STRESS_GC
        collectGarbage();
#endif
    }
  if (newSize == 0) {
    free(pointer);
    return NULL;
  }
  void *result = realloc(pointer, newSize);
  if (result == NULL)
    exit(1);
  return result;
}

void freeObject(Obj *object) {
#ifdef DEBUG_LOG_GC
    printf("%p free type %d\n", (void*)object, object->type);
#endif
  switch (object->type) {
  case OBJ_STRING: {
    ObjString *string = (ObjString *)object;
    FREE_ARRAY(char, string->chars, string->length + 1);
    FREE(ObjString, object);
    break;
  }
  case OBJ_FUNCTION: {
    ObjFunction *func = (ObjFunction *)object;
    freeChunk(&func->chunk);
    FREE(ObjFunction, object);
    break;
  }
  case OBJ_CLOSURE: {
    FREE(ObjClosure, object);
    ObjClosure* closure = (ObjClosure*)object;
    FREE_ARRAY(ObjUpvalue*, closure->upvalues, closure->upvalueCount);
    break;
  }

  case OBJ_NATIVE: {
    FREE(OBJ_NATIVE, object);
    break;
  }
  case OBJ_UPVALUE: {
      FREE(OBJ_UPVALUE, object);
      break;
  }
  }
}
void freeObjects() {
  Obj *object = vm.objects;
  while (object != NULL) {
    Obj *next = object->next;
    freeObject(object);
    object = next;
  }
}

void markObject(Obj* obj) {
    if (obj == NULL) return;
#ifdef DEBUG_LOG_GC
    printf("%p mark ", (void*)obj);
    printValue(OBJ_VAL(obj));
    printf("\n");
#endif
    obj->isMarked = true;
}

void markValue(Value value) {
    if (IS_OBJ(value)) markObject(AS_OBJ(value));
}

static void markRoots() {
    for (Value* slot = vm.stack; slot < vm.stackTop; slot++) {
        markValue(*slot);
    }

    markTable(&vm.globals);
}

void collectGarbage() {
#ifdef DEBUG_LOG_GC
    printf("-- gc begin \n");
#endif
    markRoots();

#ifdef DEBUG_LOG_GC
    printf("-- gc end\n");
#endif
}
