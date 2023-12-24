#include "common.h"
#include "chunk.h"
#include "vm.h"
#include "debug.h"
#include <stdio.h>

int main() {
    initVM();
    
    Chunk chunk;
    initChunk(&chunk);

    int constant = addConstant(&chunk, 1.2);
    writeChunk(&chunk, OP_CONSTANT, 123);
    writeChunk(&chunk, constant, 123);

    writeConstant(&chunk, 3.4, 123);
    writeChunk(&chunk, OP_ADD, 123);
    
    writeConstant(&chunk, 5.6, 123);
    writeChunk(&chunk, OP_DIVIDE, 123);

    writeChunk(&chunk, OP_NEGATE, 123);
    writeChunk(&chunk, OP_RETURN, 123);

    //disassembleChunk(&chunk, "1st Test Chunk");
    interpret(&chunk);
    freeVM();
    return 0;
}

/*
    //Testing Long Constant instruction works
    Chunk chunk;
    initChunk(&chunk);

    int constant = addConstant(&chunk, 1.2);
    writeChunk(&chunk, OP_CONSTANT, 123);
    writeChunk(&chunk, constant, 123);
    writeChunk(&chunk, OP_RETURN, 123);

    for (int i = 0; i < 300; i++) {
        writeConstant(&chunk, i, i);
    }

    disassembleChunk(&chunk, "1st Test Chunk");*/
