#include "memory.h"
#include "chunk.h"
#include "compiler.h"
#include "object.h"
#include "table.h"
#include "value.h"
#include "vm.h"
#include <stdlib.h>
#include <time.h>
#ifdef DEBUG_LOG_GC
#include <stdio.h>
#endif

#define GC_HEAP_GROW_FACTOR 2

void *reallocate(void *pointer, size_t oldSize, size_t newSize)
{
  vm.bytesAllocated = oldSize + newSize;
  if (newSize > oldSize)
  {
#ifdef DEBUG_STRESS_GC
    collectGarbage();
#endif
    if (vm.bytesAllocated > vm.nextGC)
    {
      collectGarbage();
    }
  }
  if (newSize == 0)
  {
    free(pointer);
    return NULL;
  }
  void *result = realloc(pointer, newSize);
  if (result == NULL)
    exit(1);
  return result;
}

void freeObject(Obj *object)
{
#ifdef DEBUG_LOG_GC
  printf("%p free type %d\n", (void *)object, object->type);
#endif
  switch (object->type)
  {
  case OBJ_STRING:
  {
    ObjString *string = (ObjString *)object;
    FREE_ARRAY(char, string->chars, string->length + 1);
    FREE(ObjString, object);
    break;
  }
  case OBJ_FUNCTION:
  {
    ObjFunction *func = (ObjFunction *)object;
    freeChunk(&func->chunk);
    FREE(ObjFunction, object);
    break;
  }
  case OBJ_CLOSURE:
  {
    FREE(ObjClosure, object);
    ObjClosure *closure = (ObjClosure *)object;
    FREE_ARRAY(ObjUpvalue *, closure->upvalues, closure->upvalueCount);
    break;
  }

  case OBJ_NATIVE:
  {
    FREE(OBJ_NATIVE, object);
    break;
  }
  case OBJ_UPVALUE:
  {
    FREE(OBJ_UPVALUE, object);
    break;
  }
  case OBJ_CLASS:
  {
    FREE(OBJ_CLASS, object);
    break;
  }
  case OBJ_INSTANCE:
  {
    ObjInstance *instance = (ObjInstance *)object;
    freeTable(&instance->fields);
    FREE(ObjInstance, object);
  }
  }
}
void freeObjects()
{
  Obj *object = vm.objects;
  while (object != NULL)
  {
    Obj *next = object->next;
    freeObject(object);
    object = next;
  }

  free(vm.greyStack);
}

void markObject(Obj *obj)
{
  if (obj == NULL)
    return;

  if (obj->isMarked == true)
    return;

#ifdef DEBUG_LOG_GC
  printf("%p mark ", (void *)obj);
  printValue(OBJ_VAL(obj));
  printf("\n");
#endif
  obj->isMarked = true;
  if (vm.greyCapacity < vm.greyCount + 1)
  {
    vm.greyCapacity = GROW_CAPACITY(vm.greyCapacity);
    vm.greyStack =
        (Obj **)realloc(vm.greyStack, sizeof(Obj *) * vm.greyCapacity);

    if (vm.greyStack == NULL)
      exit(1);
  }
}
void markValue(Value value)
{

  if (IS_OBJ(value))
    markObject(AS_OBJ(value));
}

static void markRoots()
{
  for (Value *slot = vm.stack; slot < vm.stackTop; slot++)
  {
    markValue(*slot);
  }

  for (int i = 0; i < vm.frameCount; i++)
  {
    markObject((Obj *)vm.frames[i].closure);
  }

  for (ObjUpvalue *upvalue = vm.openUpvalues; upvalue != NULL;
       upvalue = upvalue->next)
  {
    markObject(&upvalue->obj);
  }

  markTable(&vm.globals);
  markCompilerRoots();
}

static void markArray(ValueArray *array)
{
  for (int i = 0; i < array->count; i++)
  {
    markValue(array->values[i]);
  }
}

static void blackenObject(Obj *object)
{
#ifdef DEBUG_LOG_GC
  printf("%p blacken", (void *)object);
  printValue((OBJ_VAL(object)));
  printf("\n");
#endif /* ifdef DEBUG_LOG_GC */
  switch (object->type)
  {
  case OBJ_NATIVE:
  case OBJ_STRING:
    break;
  case OBJ_FUNCTION:
  {
    ObjFunction *function = (ObjFunction *)object;
    markObject((Obj *)function);
    markArray(&function->chunk.constants);
    break;
  }
  case OBJ_CLOSURE:
  {
    ObjClosure *closure = (ObjClosure *)object;
    markObject((Obj *)closure->function);
    for (int i = 0; i < closure->upvalueCount; i++)
    {
      markObject((Obj *)closure->upvalues[i]);
    }
    break;
  }
  case OBJ_UPVALUE:
  {
    markValue(((ObjUpvalue *)object)->closed);
    break;
  }
  case OBJ_CLASS:
  {
    ObjClass *klass = (ObjClass *)object;
    markObject((Obj *)klass->name);
  }
  case OBJ_INSTANCE:
  {
    ObjInstance* instance = (ObjInstance*)object;
    markObject((Obj*)instance->klass);
    markTable(&instance->fields);
    break;
  }
  }
}

static void traceReferences()
{
  while (vm.greyCount > 0)
  {
    Obj *object = vm.greyStack[--vm.greyCount];
    blackenObject(object);
  }
}

static void sweep()
{
  Obj *previous = NULL;
  Obj *object = vm.objects;
  while (object != NULL)
  {
    if (object->isMarked)
    {
      object->isMarked = false;
      previous = object;
      object = object->next;
    }
    else
    {
      Obj *unreached = object;
      object = object->next;
      if (previous != NULL)
      {
        previous->next = object;
      }
      else
      {
        vm.objects = object;
      }

      freeObject(unreached);
    }
  }
}

void collectGarbage()
{
#ifdef DEBUG_LOG_GC
  printf("-- gc begin \n");
  size_t before = vm.bytesAllocated;
#endif
  markRoots();
  traceReferences();
  tableRemoveWhite(&vm.strings);
  sweep();

  vm.nextGC = vm.bytesAllocated * GC_HEAP_GROW_FACTOR;

#ifdef DEBUG_LOG_GC
  printf("-- gc end\n");
  printf("    collected %zu bytes (from %zu to %zu) next at %zu\n",
         before - vm.bytesAllocated, before, vm.bytesAllocated, vm.nextGC);
#endif
}
