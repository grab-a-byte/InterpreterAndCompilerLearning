#ifndef clox_chunk_h
#define clox_chunk_h

#include "common.h"
#include "value.h"

typedef enum {
    OP_RETURN,
    OP_CONSTANT,
} OpCode;

typedef struct {
    int count;
    int capacity;
    uint8_t *code;
    ValueArray constants;
    int* lines;
} Chunk;

void initChunk(Chunk *chunk);

void freeChunk(Chunk *chunk);

void writeChunk(Chunk *chunk, uint8_t byte, int line);

int addConstant(Chunk *chunk, Value constant);

#endif // !clox_chunk_h