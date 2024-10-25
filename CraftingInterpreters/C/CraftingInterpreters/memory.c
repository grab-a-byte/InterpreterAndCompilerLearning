#ifndef clox_memory_c
#define clox_memory_c

#include <stdlib.h>

#include "memory.h"
#include "vm.h"
#include "object.h"

void* reallocate(void* pointer, size_t oldSize, size_t newSize) {
	if (0 == newSize) {
		free(pointer);
		return NULL;
	}

	void* result = realloc(pointer, newSize);
	if (result == NULL) exit(1);
	return result;
}

static void freeObject(Obj* object) {
	switch (object->type) {
		case OBJ_STRING: {
			const ObjString* string = (ObjString*)object;
			FREE_ARRAY(char, string->chars, string->length +1);
			FREE(ObjString, object);
		}
	}
}

void freeObjects() {
	Obj* object = vm.objects;
	while(object != NULL) {
		Obj* next = object->next;
		freeObject(object);
		object = next;
	}
}

#endif // !clox_memory_c
