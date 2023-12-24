const std = @import("std");
const debug = @import("debug.zig");
const Chunk = @import("chunk.zig").Chunk;
const OpCode = @import("chunk.zig").OpCode;
const Value = @import("value.zig").Value;

pub const InterpretResult = enum {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR,
};

const STACK_MAX = 256;

pub const VM = struct {
    const Self = @This();

    chunk: *Chunk,
    ip: u8,
    stack: [STACK_MAX]Value,
    stackTop: u8, //Points to where the next value will be placed

    pub fn init(chunk: *Chunk) Self {
        var vm = VM{ .chunk = chunk, .ip = 0, .stack = std.mem.zeroes([STACK_MAX]Value), .stackTop = 0 };
        return vm;
    }

    pub fn interpret(self: *Self, isDebug: bool) InterpretResult {
        return self.run(isDebug);
    }

    fn run(self: *Self, isDebug: bool) InterpretResult {
        while (true) {
            const byte = self.chunk.code.items[self.ip];
            if (isDebug) {
                for (0..self.stackTop) |i| {
                    std.debug.print("[", .{});
                    debug.printValue(self.stack[i]);
                    std.debug.print("]", .{});
                }
                std.debug.print("\n", .{});
                _ = debug.dissassembleInstruction(self.chunk, @as(usize, self.ip), byte);
            }
            const instruction: OpCode = @enumFromInt(self.chunk.code.items[self.ip]);
            self.ip += 1;
            switch (instruction) {
                OpCode.OP_RETURN => {
                    debug.printValue(self.pop());
                    std.debug.print("\n", .{});
                    return InterpretResult.INTERPRET_OK;
                },
                OpCode.OP_CONSTANT => {
                    const constant = self.chunk.constants.items[self.chunk.code.items[self.ip]];
                    self.push(constant);
                    self.ip += 1;
                },
                OpCode.OP_NEGATE => {
                    self.push(-(self.pop()));
                },
                OpCode.OP_ADD => {
                    const b = self.pop();
                    const a = self.pop();
                    self.push(a + b);
                },
                OpCode.OP_SUBTRACT => {
                    const b = self.pop();
                    const a = self.pop();
                    self.push(a - b);
                },
                OpCode.OP_MULTIPLY => {
                    const b = self.pop();
                    const a = self.pop();
                    self.push(a * b);
                },
                OpCode.OP_DIVIDE => {
                    const b = self.pop();
                    const a = self.pop();
                    self.push(a / b);
                },
            }
        }
    }

    pub fn push(self: *Self, value: Value) void {
        self.stack[self.stackTop] = value;
        self.stackTop += 1;
    }

    pub fn pop(self: *Self) Value {
        self.stackTop -= 1;
        return self.stack[self.stackTop];
    }

    pub fn free(self: *Self) void {
        self.chunk.free();
    }

    fn resetStack(self: *Self) void {
        self.stackTop = &self.stack[0];
    }
};
