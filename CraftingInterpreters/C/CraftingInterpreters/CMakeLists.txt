cmake_minimum_required(VERSION 3.29)
project(CraftingInterpreters C)

set(CMAKE_C_STANDARD 11)

include_directories(.)

add_executable(CraftingInterpreters
        chunk.c
        chunk.h
        common.h
        compiler.h
        CraftingInterpreters.c
        memory.c
        memory.h
        debug.h
        debug.c
        value.h
        value.c
        vm.h
        vm.c
        compiler.c
        scanner.c
        scanner.h
        object.h
        object.c
        table.c
        table.h)
