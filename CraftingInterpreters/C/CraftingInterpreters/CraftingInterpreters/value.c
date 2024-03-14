#include "memory.h"
#include "value.h"
#include "object.h"
#include <stdio.h>
#include <string.h>

void initValueArray(ValueArray* arr) {
	arr->capacity = 0;
	arr->count = 0;
	arr->values = NULL;
}

void writeValueArray(ValueArray* arr, Value value) {
	if (arr->capacity < arr->count + 1) {
		int oldCapacity = arr->capacity;
		arr->capacity = GROW_CAPACITY(oldCapacity);
		arr->values = GROW_ARRAY(Value, arr->values, oldCapacity, arr->capacity);
	}
	arr->values[arr->count] = value;
	arr->count++;
}

void freeValueArray(ValueArray* arr) {
	FREE_ARRAY(Value, arr->values, arr->capacity);
	initValueArray(arr);
}

void printObject(Value value) {
	switch(OBJECT_TYPE(value)) {
		case OBJ_STRING: printf("%s", AS_CSTRING(value)); break;
	}
}

void printValue(Value value) {
	switch (value.type) {
		case VAL_BOOL: printf(AS_BOOL(value) ? "true" : "false"); break;
		case VAL_NIL: printf("nil"); break;
		case VAL_NUMBER: printf("%g", AS_NUMBER(value)); break;
		case VAL_OBJ: printObject(value);
	}
}

bool valuesEqual(Value a, Value b) {
	if (a.type != b.type) return false;
	switch (a.type) {
		case VAL_BOOL: return AS_BOOL(a) == AS_BOOL(b);
		case VAL_NIL: return true;
		case VAL_NUMBER: return AS_NUMBER(a) == AS_NUMBER(b);
		case VAL_OBJ: return AS_OBJ(a) == AS_OBJ(b); 
		default: return false; //unreacheeable
	}
}

