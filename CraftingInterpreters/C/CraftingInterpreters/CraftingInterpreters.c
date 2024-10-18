#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "vm.h"

static char *readFile(char *path) {
    FILE *file = fopen(path, "rb");
    if (file == NULL) {
        fprintf(stderr, "Could not open file %s\n", path);
       exit(74);
    }
    fseek(file, 0L, SEEK_END);
    size_t filesize = ftell(file);
    rewind(file);

    char *buffer = (char *) malloc(filesize + 1);
    if (buffer == NULL) {
        fprintf(stderr, "Not enough memory to read file %s\n", path);
        exit(74);
    }
    size_t bytesRead = fread(buffer, sizeof(char), filesize, file);

    if(bytesRead < filesize) {
        fprintf(stderr, "Could not read file %s", path);
        exit(74);
    }
    buffer[bytesRead] = "\0";

    fclose(file);
    return buffer;
}

static void repl() {
    char line[1024];
    for (;;) {
        printf(">");
        if (!fgets(line, sizeof(line), stdin)) {
            printf("\n");
            break;
        }
        printf("interpreting");
        interpret(line);
    }
}

static void runFile(char *path) {
    char *source = readFile(path);
    const InterpretResult result = interpret(source);
    free(source);
    if (result == INTERPRET_COMPILE_ERROR) exit(65);
    if (result == INTERPRET_RUNTIME_ERROR) exit(70);
}


int main(int argc, const char *argv[]) {
    initVM();
    if (argc == 1) {
        repl();
    } else if (argc == 2) {
        runFile(argv[1]);
    } else {
        fprintf(stderr, "Usages: clox [path]\n");
        exit(64);
    }
    return 0;
}
